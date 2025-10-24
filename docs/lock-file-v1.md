# Zoi Lock File v1 Format

The new `zoi.lock` format (version 1) provides reproducible builds with integrity verification for project-scoped package installations.

## Structure

```json
{
  "version": "1",
  "packages": {
    "#registry1@repo/path/to/package1": "1.0.0",
    "#registry2@repo/path/to/package2": "2.0.0"
  },
  "registries": {
    "registry1": "https://github.com/user/registry1.git",
    "registry2": "https://github.com/user/registry2.git"
  },
  "#registry1": {
    "@repo/path/to/package1": {
      "version": "1.0.0",
      "integrity": "sha512-of-the-entire-store-folder-that-created-locally",
      "dependencies": [
        "#registry3@repo/path/to/package3"
      ],
      "options_dependencies": [],
      "optionals_dependencies": []
    }
  },
  "#registry2": {
    "@repo/path/to/package2": {
      "version": "2.0.0",
      "integrity": "sha512-of-the-entire-store-folder-that-created-locally",
      "dependencies": [],
      "options_dependencies": [],
      "optionals_dependencies": []
    }
  }
}
```

## Fields

### Top Level

- **version**: The lock file format version (currently "1")
- **packages**: Map of full package IDs to their installed versions
- **registries**: Map of registry handles to their git URLs
- **#<registry>**: Per-registry package details (one section per registry)

### Package Details (per registry)

Each package entry under a registry section contains:

- **version**: The installed version
- **integrity**: SHA-512 hash of the entire package store directory
- **dependencies**: List of zoi dependencies in full ID format
- **options_dependencies**: List of chosen option dependencies
- **optionals_dependencies**: List of chosen optional dependencies

## Commands

### Add Packages

Add packages to the current project:

```bash
# Add a package to the project
zoi add package-name

# Add and save to zoi.yaml
zoi add package-name --save

# Alternative syntax
zoi install package-name --local
zoi install package-name --scope project
```

### Remove Packages

Remove packages from the current project:

```bash
# Remove a package from the project
zoi remove package-name

# Remove and update zoi.yaml
zoi remove package-name --save

# Alternative syntax
zoi uninstall package-name --local
```

### Install from Lock File

Install all packages defined in `zoi.yaml` and verify against `zoi.lock`:

```bash
zoi install
```

This will:
1. Check that all registries in the lock file are synced
2. Install packages with exact versions from the lock file
3. Verify integrity hashes match
4. Exit with error if any verification fails

## Project Configuration

### zoi.yaml

Add a `registry` field to specify the default registry for the project:

```yaml
name: my-project
registry: https://github.com/user/my-registry.git

config:
  local: true

pkgs:
  - package1
  - package2
```

### Registry Verification

When installing with project scope, Zoi verifies that:

1. All registries in `zoi.lock` are synced locally
2. Registry handles and URLs match exactly
3. Package integrity hashes match the lock file

If registries are missing, Zoi will display:

```
You don't have these registries:
  - registry1 (https://github.com/user/registry1.git)
  - registry2 (https://github.com/user/registry2.git)

You can add them with this command:
  zoi sync add <registry-url>
```

## Workflow

### Initial Setup

1. Create `zoi.yaml` with `local: true` and `registry` field
2. Add packages: `zoi add package1 package2 --save`
3. This creates/updates `zoi.lock` and `zoi.pkgs.json`

### Team Collaboration

1. Clone the repository
2. Run `zoi sync` to ensure registries are available
3. Run `zoi install` to install exact versions from lock file
4. Zoi verifies integrity to ensure reproducible environment

### Adding/Removing Packages

```bash
# Add new package
zoi add new-package --save

# Remove package
zoi remove old-package --save

# Both commands update zoi.lock, zoi.pkgs.json, and zoi.yaml
```

## Migration

The new lock file format coexists with the existing simple format. Projects using `config.local: true` in `zoi.yaml` will automatically use the new verification features.
