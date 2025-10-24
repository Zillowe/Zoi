# Changelog: Lock File v1 Format and Project-Scope Commands

## New Features

### Lock File v1 Format

Added a new lock file format (`zoi.lock` version 1) that provides:

- **Reproducible Builds**: Exact package versions with SHA-512 integrity verification
- **Registry Tracking**: Records which registries are used and their URLs
- **Dependency Information**: Tracks dependencies, option dependencies, and optional dependencies per package
- **Integrity Verification**: Hash verification of entire package store directories

Format structure:
```json
{
  "version": "1",
  "packages": { ... },
  "registries": { ... },
  "#registry": { ... }
}
```

### New Commands

#### `zoi add`

Adds packages to the current project (alias for `install --local`):

```bash
zoi add <package>...
zoi add <package> --save    # Also update zoi.yaml
```

Features:
- Installs to project scope automatically
- Updates `zoi.lock` and `zoi.pkgs.json`
- With `--save`: adds to `zoi.yaml` pkgs field

#### `zoi remove`

Removes packages from the current project (alias for `uninstall --local`):

```bash
zoi remove <package>...
zoi remove <package> --save    # Also update zoi.yaml
```

Features:
- Uninstalls from project scope
- Updates `zoi.lock` and `zoi.pkgs.json`
- With `--save`: removes from `zoi.yaml` pkgs field

### Updated Commands

#### `zoi install`

Enhanced project-scope installation:

- Verifies registries from `zoi.lock` are synced
- Checks package integrity against lock file
- Warns about missing registries with helpful commands

#### `zoi uninstall`

Added `--save` flag:

```bash
zoi uninstall <package> --local --save
```

Removes package from `zoi.yaml` when using project scope.

### Project Configuration

#### zoi.yaml

New `registry` field to specify project registry:

```yaml
name: my-project
registry: https://github.com/user/registry.git

config:
  local: true

pkgs:
  - package1
  - package2
```

## Workflow

### Initial Project Setup

1. Create `zoi.yaml` with `config.local: true`
2. Add packages: `zoi add package1 package2 --save`
3. Commit `zoi.yaml`, `zoi.lock`, and `zoi.pkgs.json`

### Team Member Setup

1. Clone repository
2. Run `zoi sync` (if needed)
3. Run `zoi install` to install exact versions

### Adding/Removing Packages

```bash
# Add new dependency
zoi add new-package --save

# Remove old dependency
zoi remove old-package --save
```

## Technical Details

### Lock File V1 Module

New module: `src/project/lockfile_v1.rs`

Functions:
- `read_lock()`: Read the v1 lock file
- `write_lock()`: Write the v1 lock file
- `add_package()`: Add a package entry with integrity hash
- `remove_package()`: Remove a package entry
- `verify_integrity()`: Check package integrity against lock
- `verify_registries()`: Verify registries are synced

### Project Config Updates

- `src/project/config.rs`:
  - Added `registry` field to `ProjectConfig`
  - Added `remove_packages_from_config()` function

### CLI Updates

- New `Add` and `Remove` commands
- Updated `Uninstall` with `--save` flag
- Lock acquisition for `Add` and `Remove` commands

## Migration

The new lock file format is opt-in and coexists with the existing simple format. Projects with `config.local: true` in `zoi.yaml` will use the enhanced verification features.

## Files Changed

- `src/pkg/types.rs`: Added `ZoiLockV1` and `ZoiLockPackageInfo` types
- `src/project/config.rs`: Added `registry` field and `remove_packages_from_config()`
- `src/project/lockfile.rs`: Added v1 lock file read/write functions
- `src/project/lockfile_v1.rs`: New module for v1 lock file operations
- `src/project/mod.rs`: Exported `lockfile_v1` module
- `src/cli.rs`: Added `Add` and `Remove` commands, updated `Uninstall`
- `src/cmd/install.rs`: Added registry verification
- `src/cmd/uninstall.rs`: Added `save` parameter

## Documentation

- `docs/lock-file-v1.md`: Comprehensive documentation
- `examples/project-scope-usage.md`: Usage examples
