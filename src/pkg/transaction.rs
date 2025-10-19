use crate::pkg::{local, recorder, types};
use anyhow::{Result, anyhow};
use colored::*;
use semver::Version;
use std::collections::{HashMap, HashSet};
use std::fs;
use std::path::{Path, PathBuf};
use std::sync::{Arc, Mutex};

#[cfg(windows)]
use junction;

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
struct PackageKey {
    scope: types::Scope,
    registry_handle: String,
    repo: String,
    name: String,
}

impl PackageKey {
    fn from_manifest(m: &types::InstallManifest) -> Self {
        Self {
            scope: m.scope,
            registry_handle: m.registry_handle.clone(),
            repo: m.repo.clone(),
            name: m.name.clone(),
        }
    }

    fn from_parts(scope: types::Scope, registry_handle: &str, repo: &str, name: &str) -> Self {
        Self {
            scope,
            registry_handle: registry_handle.to_string(),
            repo: repo.to_string(),
            name: name.to_string(),
        }
    }
}

pub type SharedTransaction = Arc<Transaction>;

#[derive(Default)]
pub struct Transaction {
    baseline: HashSet<PackageKey>,
    pre_state: Mutex<HashMap<PackageKey, Option<types::InstallManifest>>>,
    // Keep order of successes to facilitate reverse processing if needed
    successes: Mutex<Vec<PackageKey>>,
    // deferred cleanup for update operations (package_dir pruning)
    cleanup: Mutex<Vec<PackageKey>>,
}

impl Transaction {
    pub fn new() -> SharedTransaction {
        let mut baseline = HashSet::new();
        if let Ok(installed) = local::get_installed_packages() {
            for m in installed {
                baseline.insert(PackageKey::from_manifest(&m));
            }
        }
        Arc::new(Transaction {
            baseline,
            pre_state: Mutex::new(HashMap::new()),
            successes: Mutex::new(Vec::new()),
            cleanup: Mutex::new(Vec::new()),
        })
    }

    pub fn register_pre_state(
        &self,
        scope: types::Scope,
        registry_handle: &str,
        repo: &str,
        name: &str,
    ) -> Result<()> {
        let key = PackageKey::from_parts(scope, registry_handle, repo, name);
        let mut pre = self.pre_state.lock().unwrap();
        pre.entry(key)
            .or_insert(local::is_package_installed(name, scope)?);
        Ok(())
    }

    pub fn register_success(
        &self,
        scope: types::Scope,
        registry_handle: &str,
        repo: &str,
        name: &str,
    ) {
        let key = PackageKey::from_parts(scope, registry_handle, repo, name);
        self.successes.lock().unwrap().push(key);
    }

    pub fn defer_cleanup(
        &self,
        scope: types::Scope,
        registry_handle: &str,
        repo: &str,
        name: &str,
    ) {
        let key = PackageKey::from_parts(scope, registry_handle, repo, name);
        self.cleanup.lock().unwrap().push(key);
    }

    pub fn commit(&self) -> Result<()> {
        // Perform any deferred cleanups (old versions pruning)
        for key in self.cleanup.lock().unwrap().iter() {
            self.cleanup_old_versions(&key.name, key.scope, &key.repo, &key.registry_handle)?;
        }
        Ok(())
    }

    pub fn rollback(&self) -> Result<()> {
        println!(
            "{}",
            "An error occurred. Rolling back changes...".yellow().bold()
        );

        // Determine installed packages after operation
        let mut current: HashMap<PackageKey, types::InstallManifest> = HashMap::new();
        if let Ok(installed) = local::get_installed_packages() {
            for m in installed {
                current.insert(PackageKey::from_manifest(&m), m);
            }
        }

        // Remove any package that did not exist in baseline
        let to_remove_keys: Vec<PackageKey> = current
            .iter()
            .filter(|(k, _)| !self.baseline.contains(*k))
            .map(|(k, _)| k.clone())
            .collect();

        for key in to_remove_keys {
            if let Some(manifest) = current.get(&key)
                && let Err(e) = self.uninstall_by_manifest(manifest)
            {
                eprintln!(
                    "Warning: failed to remove newly installed '{}': {}",
                    manifest.name, e
                );
            }
        }

        // For packages that existed before, try reverting to previous version
        for (key, pre) in self.pre_state.lock().unwrap().iter() {
            if let Some(prev_manifest) = pre {
                // If the version changed, revert to prev
                if let Some(current_manifest) = current.get(key)
                    && current_manifest.version != prev_manifest.version
                    && let Err(e) = self.revert_to_manifest(prev_manifest)
                {
                    eprintln!(
                        "Warning: failed to revert '{}' to {}: {}",
                        prev_manifest.name, prev_manifest.version, e
                    );
                }
            } else {
                // pre was None, but current might still exist (e.g., installed succeeded partially then baseline didn't have it). Already handled in to_remove above.
            }
        }

        println!("{}", "Rollback complete.".green());
        Ok(())
    }

    fn cleanup_old_versions(
        &self,
        package_name: &str,
        scope: types::Scope,
        repo: &str,
        registry_handle: &str,
    ) -> Result<()> {
        let config = crate::pkg::config::read_config()?;
        let rollback_enabled = config.rollback_enabled;

        let package_dir = local::get_package_dir(scope, registry_handle, repo, package_name)?;
        let mut versions = Vec::new();
        if let Ok(entries) = fs::read_dir(&package_dir) {
            for entry in entries.flatten() {
                let path = entry.path();
                if path.is_dir()
                    && let Some(version_str) = path.file_name().and_then(|s| s.to_str())
                    && version_str != "latest"
                    && let Ok(version) = Version::parse(version_str)
                {
                    versions.push(version);
                }
            }
        }
        if versions.is_empty() {
            return Ok(());
        }
        versions.sort();

        let versions_to_keep = if rollback_enabled { 2 } else { 1 };
        if versions.len() > versions_to_keep {
            let num_to_delete = versions.len() - versions_to_keep;
            let versions_to_delete = &versions[..num_to_delete];
            for version in versions_to_delete {
                let dir = package_dir.join(version.to_string());
                if dir.exists() {
                    fs::remove_dir_all(dir)?;
                }
            }
        }
        Ok(())
    }

    fn get_bin_root(&self, scope: types::Scope) -> Result<PathBuf> {
        match scope {
            types::Scope::User => {
                let home_dir =
                    home::home_dir().ok_or_else(|| anyhow!("Could not find home directory."))?;
                Ok(home_dir.join(".zoi/pkgs/bin"))
            }
            types::Scope::System => {
                if cfg!(target_os = "windows") {
                    Ok(PathBuf::from("C:\\ProgramData\\zoi\\pkgs\\bin"))
                } else {
                    Ok(PathBuf::from("/usr/local/bin"))
                }
            }
            types::Scope::Project => {
                let current_dir = std::env::current_dir()?;
                Ok(current_dir.join(".zoi").join("pkgs").join("bin"))
            }
        }
    }

    fn create_symlink(&self, target: &Path, link: &Path) -> Result<()> {
        if link.exists() {
            fs::remove_file(link)?;
        }
        #[cfg(unix)]
        {
            std::os::unix::fs::symlink(target, link)?;
        }
        #[cfg(windows)]
        {
            fs::copy(target, link)?;
        }
        Ok(())
    }

    fn revert_to_manifest(&self, prev_manifest: &types::InstallManifest) -> Result<()> {
        let package_dir = local::get_package_dir(
            prev_manifest.scope,
            &prev_manifest.registry_handle,
            &prev_manifest.repo,
            &prev_manifest.name,
        )?;

        let latest_symlink_path = package_dir.join("latest");
        // Try determine current version dir via manifest to be robust across symlink/junctions
        let current_version_dir = {
            let manifest_in_latest = latest_symlink_path.join("manifest.yaml");
            if manifest_in_latest.exists() {
                if let Ok(content) = fs::read_to_string(&manifest_in_latest) {
                    if let Ok(cur_manifest) =
                        serde_yaml::from_str::<types::InstallManifest>(&content)
                    {
                        package_dir.join(cur_manifest.version)
                    } else {
                        fs::read_link(&latest_symlink_path)
                            .unwrap_or_else(|_| latest_symlink_path.clone())
                    }
                } else {
                    fs::read_link(&latest_symlink_path)
                        .unwrap_or_else(|_| latest_symlink_path.clone())
                }
            } else if latest_symlink_path.exists() || latest_symlink_path.is_symlink() {
                fs::read_link(&latest_symlink_path).unwrap_or_else(|_| latest_symlink_path.clone())
            } else {
                latest_symlink_path.clone()
            }
        };

        let previous_version_dir = package_dir.join(&prev_manifest.version);

        if latest_symlink_path.exists() || latest_symlink_path.is_symlink() {
            if latest_symlink_path.is_dir() {
                fs::remove_dir_all(&latest_symlink_path)?;
            } else {
                let _ = fs::remove_file(&latest_symlink_path);
            }
        }

        #[cfg(unix)]
        std::os::unix::fs::symlink(&previous_version_dir, &latest_symlink_path)?;
        #[cfg(windows)]
        {
            junction::create(&previous_version_dir, &latest_symlink_path)?;
        }

        if let Some(bins) = &prev_manifest.bins {
            let bin_root = self.get_bin_root(prev_manifest.scope)?;
            for bin in bins {
                let symlink_path = bin_root.join(bin);
                if symlink_path.exists() {
                    let _ = fs::remove_file(&symlink_path);
                }
                let bin_path_in_store = previous_version_dir.join("bin").join(bin);
                if bin_path_in_store.exists() {
                    self.create_symlink(&bin_path_in_store, &symlink_path)?;
                }
            }
        }

        // Remove current version dir if it differs and exists
        if current_version_dir != previous_version_dir && current_version_dir.exists() {
            // Before removing, try to delete installed files listed by the current manifest
            let current_manifest_path = current_version_dir.join("manifest.yaml");
            if let Ok(content) = fs::read_to_string(&current_manifest_path)
                && let Ok(cur_manifest) = serde_yaml::from_str::<types::InstallManifest>(&content)
            {
                self.remove_installed_files(&cur_manifest)?;
            }
            fs::remove_dir_all(current_version_dir)?;
        }

        Ok(())
    }

    fn remove_installed_files(&self, manifest: &types::InstallManifest) -> Result<()> {
        // Remove files recorded outside the store (usrroot/usrhome)
        for path_str in &manifest.installed_files {
            let path = PathBuf::from(path_str);
            if path.exists() {
                if path.is_dir() {
                    let _ = fs::remove_dir_all(&path);
                } else {
                    let _ = fs::remove_file(&path);
                }
            }
        }
        Ok(())
    }

    fn uninstall_by_manifest(&self, manifest: &types::InstallManifest) -> Result<()> {
        // Remove symlink(s)
        let bin_root = self.get_bin_root(manifest.scope)?;
        if let Some(bins) = &manifest.bins {
            for bin in bins {
                let symlink_path = bin_root.join(bin);
                if symlink_path.exists() {
                    let _ = fs::remove_file(&symlink_path);
                }
            }
        } else {
            let symlink_path = bin_root.join(&manifest.name);
            if symlink_path.exists() {
                let _ = fs::remove_file(&symlink_path);
            }
        }

        // Remove package directory
        let package_dir = local::get_package_dir(
            manifest.scope,
            &manifest.registry_handle,
            &manifest.repo,
            &manifest.name,
        )?;
        if package_dir.exists() {
            let _ = fs::remove_dir_all(&package_dir);
        }

        // Remove external installed files
        self.remove_installed_files(manifest)?;

        // Remove lockfile record
        if let Err(e) = recorder::remove_package_from_record(&manifest.name, manifest.scope) {
            eprintln!(
                "{} Failed to remove package from lockfile: {}",
                "Warning:".yellow(),
                e
            );
        }

        Ok(())
    }
}
