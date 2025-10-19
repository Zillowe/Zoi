use crate::pkg::{install, local, resolve, transaction, types, types::InstallReason};
use crate::utils;
use colored::*;
use std::collections::HashSet;
use std::sync::Mutex;

pub fn run(sources: &[String], force: bool, yes: bool) {
    for source in sources {
        println!(
            "--- Building package '{}' from source ---",
            source.blue().bold()
        );

        if let Ok(Some(manifest)) = local::is_package_installed(source, types::Scope::User)
            && !force
        {
            println!(
                "Package '{}' is already installed. Use --force to rebuild.",
                manifest.name.yellow()
            );
            continue;
        }

        match resolve::resolve_source(source) {
            Ok(resolved_source) => {
                if let Err(e) = utils::confirm_untrusted_source(&resolved_source.source_type, yes) {
                    eprintln!("\n{}", e.to_string().red());
                    return;
                }

                if let Some(repo_name) = &resolved_source.repo_name {
                    utils::print_repo_warning(repo_name);
                }

                let processed_deps = Mutex::new(HashSet::new());
                let tx = transaction::Transaction::new();
                if let Err(e) = install::run_installation(
                    resolved_source.path.to_str().unwrap(),
                    install::InstallMode::ForceBuild,
                    true,
                    InstallReason::Direct,
                    yes,
                    false,
                    &processed_deps,
                    None,
                    None,
                    Some(tx.clone()),
                ) {
                    let _ = tx.rollback();
                    eprintln!("\n{}: {}", "Build failed".red().bold(), e);
                    std::process::exit(1);
                }
                if let Err(e) = tx.commit() {
                    eprintln!("\n{}: {}", "Build finalize failed".red().bold(), e);
                    std::process::exit(1);
                }
            }
            Err(e) => {
                eprintln!("\n{}: {}", "Error".red().bold(), e);
                std::process::exit(1);
            }
        }
        println!("\n{}", "Build complete.".green());
    }
}
