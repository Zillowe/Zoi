use crate::pkg::{hooks, install, local, pin, resolve, sync, transaction, types};
use crate::utils;
use anyhow::{Result, anyhow};
use colored::*;
use indicatif::{ProgressBar, ProgressStyle};
use rayon::prelude::*;
use std::collections::HashSet;
use std::sync::Mutex;

pub fn run(all: bool, package_names: &[String], yes: bool) {
    if all {
        if let Err(e) = run_update_all_logic(yes) {
            eprintln!("{}: {}", "Update failed".red().bold(), e);
        }
        return;
    }

    println!("{}", "--- Syncing Package Database ---".yellow().bold());
    if let Err(e) = sync::run(false, true, true) {
        eprintln!("{}: {}", "Sync failed".red().bold(), e);
        std::process::exit(1);
    }

    let mut failed_packages = Vec::new();

    for (i, package_name) in package_names.iter().enumerate() {
        if i > 0 {
            println!();
        }
        if let Err(e) = run_update_single_logic(package_name, yes) {
            eprintln!(
                "{}: Failed to update '{}': {}",
                "Error".red().bold(),
                package_name,
                e
            );
            failed_packages.push(package_name.clone());
        }
    }

    if !failed_packages.is_empty() {
        eprintln!(
            "\n{}: The following packages failed to update:",
            "Error".red().bold()
        );
        for pkg in &failed_packages {
            eprintln!("  - {}", pkg);
        }
        std::process::exit(1);
    } else if !package_names.is_empty() {
        println!(
            "\n{}",
            "Update process finished for all specified packages.".green()
        );
    }
}

fn run_update_single_logic(package_name: &str, yes: bool) -> Result<()> {
    println!("--- Updating package '{}' ---", package_name.blue().bold());

    let (new_pkg, new_version, _, _, registry_handle) =
        match resolve::resolve_package_and_version(package_name) {
            Ok(result) => result,
            Err(e) => {
                return Err(anyhow!(
                    "Could not resolve package '{}': {}",
                    package_name,
                    e
                ));
            }
        };

    if pin::is_pinned(package_name)? {
        println!(
            "Package '{}' is pinned. Skipping update.",
            package_name.yellow()
        );
        return Ok(());
    }

    let manifest = match local::is_package_installed(&new_pkg.name, types::Scope::User)?.or(
        local::is_package_installed(&new_pkg.name, types::Scope::System)?,
    ) {
        Some(m) => m,
        None => {
            return Err(anyhow!(
                "Package '{package_name}' is not installed. Use 'zoi install' instead."
            ));
        }
    };

    println!("Currently installed version: {}", manifest.version.yellow());
    println!("Available version: {}", new_version.green());

    if manifest.version == new_version {
        println!("\nPackage is already up to date.");
        return Ok(());
    }

    if !utils::ask_for_confirmation(
        &format!("Update from {} to {}?", manifest.version, new_version),
        yes,
    ) {
        return Ok(());
    }

    if let Some(hooks) = &new_pkg.hooks
        && let Err(e) = hooks::run_hooks(hooks, hooks::HookType::PreUpgrade)
    {
        return Err(anyhow!("Pre-upgrade hook failed: {}", e));
    }

    let mode = install::InstallMode::PreferPrebuilt;
    let tx = transaction::Transaction::new();

    let processed_deps = Mutex::new(HashSet::new());
    if let Err(e) = install::run_installation(
        package_name,
        mode,
        true,
        types::InstallReason::Direct,
        yes,
        false,
        &processed_deps,
        Some(manifest.scope),
        None,
        Some(tx.clone()),
    ) {
        let _ = tx.rollback();
        return Err(e);
    }

    // Defer cleanup until commit
    tx.defer_cleanup(
        manifest.scope,
        registry_handle.as_deref().unwrap_or("local"),
        &new_pkg.repo,
        &new_pkg.name,
    );

    if let Some(hooks) = &new_pkg.hooks
        && let Err(e) = hooks::run_hooks(hooks, hooks::HookType::PostUpgrade)
    {
        let _ = tx.rollback();
        return Err(anyhow!("Post-upgrade hook failed: {}", e));
    }

    if let Err(e) = tx.commit() {
        return Err(anyhow!("Failed to finalize update: {}", e));
    }

    println!("\n{}", "Update complete.".green());
    Ok(())
}

fn run_update_all_logic(yes: bool) -> Result<()> {
    println!("{}", "--- Syncing Package Database ---".yellow().bold());
    sync::run(false, true, true)?;

    let installed_packages = local::get_installed_packages()?;
    let pinned_packages = pin::get_pinned_packages()?;
    let pinned_sources: Vec<String> = pinned_packages.into_iter().map(|p| p.source).collect();

    let mut packages_to_upgrade = Vec::new();
    let mut upgrade_messages = Vec::new();

    println!("\n{}", "--- Checking for Upgrades ---".yellow().bold());
    let pb = ProgressBar::new(installed_packages.len() as u64);
    pb.set_style(
        ProgressStyle::default_bar()
            .template(
                "{spinner:.green} [{elapsed_precise}] [{bar:40.cyan/blue}] {pos}/{len} ({msg})",
            )?
            .progress_chars("#>-"),
    );
    pb.set_message("Checking packages...");

    for manifest in installed_packages {
        let source = format!("#{}@{}", manifest.registry_handle, manifest.repo);
        if pinned_sources.contains(&source) {
            upgrade_messages.push(format!("- {} is pinned, skipping.", manifest.name.cyan()));
            continue;
        }

        let (new_pkg, new_version, _, _, registry_handle) =
            match resolve::resolve_package_and_version(&source) {
                Ok(result) => result,
                Err(e) => {
                    upgrade_messages.push(format!(
                        "- Could not resolve package '{}': {}, skipping.",
                        manifest.name, e
                    ));
                    continue;
                }
            };

        if manifest.version != new_version {
            upgrade_messages.push(format!(
                "- {} can be upgraded from {} to {}",
                manifest.name.cyan(),
                manifest.version.yellow(),
                new_version.green()
            ));
            packages_to_upgrade.push((source.clone(), new_pkg, registry_handle, manifest.scope));
        } else {
            upgrade_messages.push(format!("- {} is up to date.", manifest.name.cyan()));
        }
        pb.inc(1);
    }
    pb.finish_and_clear();

    for msg in upgrade_messages {
        println!("{}", msg);
    }

    if packages_to_upgrade.is_empty() {
        println!("\n{}", "All packages are up to date.".green());
        return Ok(());
    }

    println!();
    if !utils::ask_for_confirmation("Do you want to upgrade these packages?", yes) {
        return Ok(());
    }

    let processed_deps = Mutex::new(HashSet::new());
    let tx = transaction::Transaction::new();

    let failures = Mutex::new(Vec::<String>::new());

    packages_to_upgrade
        .par_iter()
        .for_each(|(source, new_pkg, registry_handle, scope)| {
            println!(
                "\n--- Upgrading {} to {} ---",
                source.cyan(),
                new_pkg.version.as_deref().unwrap_or("N/A").green()
            );

            if let Some(hooks) = &new_pkg.hooks
                && let Err(e) = hooks::run_hooks(hooks, hooks::HookType::PreUpgrade)
            {
                eprintln!(
                    "{}: Pre-upgrade hook failed for '{}': {}",
                    "Error".red().bold(),
                    source,
                    e
                );
                failures.lock().unwrap().push(source.clone());
                return;
            }

            let mode = install::InstallMode::PreferPrebuilt;
            if let Err(e) = install::run_installation(
                source,
                mode,
                true,
                types::InstallReason::Direct,
                yes,
                false,
                &processed_deps,
                Some(*scope),
                None,
                Some(tx.clone()),
            ) {
                eprintln!("Failed to upgrade {}: {}", source, e);
                failures.lock().unwrap().push(source.clone());
                return;
            }

            // Defer cleanup to commit
            tx.defer_cleanup(
                *scope,
                registry_handle.as_deref().unwrap_or("local"),
                &new_pkg.repo,
                &new_pkg.name,
            );

            if let Some(hooks) = &new_pkg.hooks
                && let Err(e) = hooks::run_hooks(hooks, hooks::HookType::PostUpgrade)
            {
                eprintln!(
                    "{}: Post-upgrade hook failed for '{}': {}",
                    "Error".red().bold(),
                    source,
                    e
                );
                failures.lock().unwrap().push(source.clone());
            }
        });

    let failures = failures.into_inner().unwrap();
    if !failures.is_empty() {
        let _ = tx.rollback();
        return Err(anyhow!("One or more upgrades failed."));
    }

    tx.commit()?;

    println!("\n{}", "Upgrade complete.".green());
    Ok(())
}
