use crate::pkg::{hash, local, types};
use anyhow::{Result, anyhow};
use std::collections::HashMap;
use std::fs;

fn get_lockfile_path() -> Result<std::path::PathBuf> {
    Ok(std::env::current_dir()?.join("zoi.lock"))
}

pub fn read_lock() -> Result<types::ZoiLockV1> {
    let path = get_lockfile_path()?;
    if !path.exists() {
        return Ok(types::ZoiLockV1 {
            version: "1".to_string(),
            packages: HashMap::new(),
            registries: HashMap::new(),
            registry_packages: HashMap::new(),
        });
    }
    let content = fs::read_to_string(path)?;
    let lockfile: types::ZoiLockV1 = serde_json::from_str(&content)?;
    Ok(lockfile)
}

pub fn write_lock(lockfile: &types::ZoiLockV1) -> Result<()> {
    let path = get_lockfile_path()?;
    let content = serde_json::to_string_pretty(lockfile)?;
    fs::write(path, content)?;
    Ok(())
}

pub fn add_package(
    pkg: &types::Package,
    registry_handle: &str,
    registry_url: &str,
    dependencies: &[String],
    options_dependencies: &[String],
    optionals_dependencies: &[String],
) -> Result<()> {
    let mut lockfile = read_lock()?;

    let version = pkg
        .version
        .as_ref()
        .ok_or_else(|| anyhow!("Package {} has no version", pkg.name))?
        .clone();

    let package_store_dir =
        local::get_package_dir(types::Scope::Project, registry_handle, &pkg.repo, &pkg.name)?;
    let latest_dir = package_store_dir.join("latest");

    let integrity = if latest_dir.exists() {
        hash::calculate_dir_hash(&latest_dir)?
    } else {
        String::new()
    };

    let full_package_id = format!("#{}@{}/{}", registry_handle, pkg.repo, pkg.name);

    lockfile
        .packages
        .insert(full_package_id.clone(), version.clone());

    lockfile
        .registries
        .insert(registry_handle.to_string(), registry_url.to_string());

    let registry_key = format!("#{}", registry_handle);
    let registry_packages = lockfile
        .registry_packages
        .entry(registry_key)
        .or_insert_with(HashMap::new);

    let package_path = format!("@{}/{}", pkg.repo, pkg.name);
    registry_packages.insert(
        package_path,
        types::ZoiLockPackageInfo {
            version,
            integrity,
            dependencies: dependencies.to_vec(),
            options_dependencies: options_dependencies.to_vec(),
            optionals_dependencies: optionals_dependencies.to_vec(),
        },
    );

    write_lock(&lockfile)
}

pub fn remove_package(package_name: &str, registry_handle: &str, repo: &str) -> Result<()> {
    let mut lockfile = read_lock()?;

    let full_package_id = format!("#{}@{}/{}", registry_handle, repo, package_name);
    lockfile.packages.remove(&full_package_id);

    let registry_key = format!("#{}", registry_handle);
    if let Some(registry_packages) = lockfile.registry_packages.get_mut(&registry_key) {
        let package_path = format!("@{}/{}", repo, package_name);
        registry_packages.remove(&package_path);

        if registry_packages.is_empty() {
            lockfile.registry_packages.remove(&registry_key);
            lockfile.registries.remove(registry_handle);
        }
    }

    write_lock(&lockfile)
}

pub fn verify_integrity(package_name: &str, registry_handle: &str, repo: &str) -> Result<bool> {
    let lockfile = read_lock()?;

    let registry_key = format!("#{}", registry_handle);
    let package_path = format!("@{}/{}", repo, package_name);

    if let Some(registry_packages) = lockfile.registry_packages.get(&registry_key)
        && let Some(pkg_info) = registry_packages.get(&package_path)
    {
        let package_store_dir =
            local::get_package_dir(types::Scope::Project, registry_handle, repo, package_name)?;
        let latest_dir = package_store_dir.join("latest");

        if !latest_dir.exists() {
            return Ok(false);
        }

        let current_hash = hash::calculate_dir_hash(&latest_dir)?;
        return Ok(current_hash == pkg_info.integrity);
    }

    Ok(false)
}

pub fn verify_registries(
    _project_config: &crate::project::config::ProjectConfig,
) -> Result<Vec<String>> {
    let lockfile = read_lock()?;
    let zoi_config = crate::pkg::config::read_config()?;

    let mut missing_registries = Vec::new();

    for (registry_handle, registry_url) in &lockfile.registries {
        let mut found = false;

        if let Some(default_reg) = &zoi_config.default_registry
            && &default_reg.handle == registry_handle
            && &default_reg.url == registry_url
        {
            found = true;
        }

        for added_reg in &zoi_config.added_registries {
            if &added_reg.handle == registry_handle && &added_reg.url == registry_url {
                found = true;
                break;
            }
        }

        if !found {
            missing_registries.push(format!("{} ({})", registry_handle, registry_url));
        }
    }

    Ok(missing_registries)
}

pub fn get_registry_url_from_config(
    registry_handle: &str,
    zoi_config: &types::Config,
) -> Option<String> {
    if let Some(default_reg) = &zoi_config.default_registry
        && default_reg.handle == registry_handle
    {
        return Some(default_reg.url.clone());
    }

    for added_reg in &zoi_config.added_registries {
        if added_reg.handle == registry_handle {
            return Some(added_reg.url.clone());
        }
    }

    None
}
