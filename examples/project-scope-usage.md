# Project-Scope Package Management Example

This example demonstrates the new project-scope package management features with lock file support.

## Setup

Create a new project with a `zoi.yaml` file:

```yaml
name: my-awesome-project
registry: https://github.com/Zillowe/Zoidberg.git

config:
  local: true

commands:
  - cmd: dev
    run: npm run dev

environments:
  - name: development
    cmd: dev
    run:
      - npm install
```

## Adding Packages

Add packages to your project:

```bash
# Add a package (updates zoi.lock and zoi.pkgs.json)
zoi add ripgrep

# Add and save to zoi.yaml
zoi add fd --save

# Add multiple packages
zoi add bat exa --save
```

After running these commands:
- `zoi.lock` will contain package versions, integrity hashes, and dependencies
- `zoi.pkgs.json` will have detailed package installation records
- `zoi.yaml` (with --save) will list the packages in the `pkgs` field

## Removing Packages

Remove packages from your project:

```bash
# Remove a package
zoi remove ripgrep

# Remove and update zoi.yaml
zoi remove fd --save
```

## Installing from Lock File

When someone clones your project:

```bash
# Clone the repository
git clone https://github.com/user/project.git
cd project

# Sync registries if needed
zoi sync

# Install all packages with exact versions from lock file
zoi install
```

This will:
1. Verify all registries are available
2. Install packages with exact versions
3. Verify integrity hashes
4. Fail if any verification doesn't match

## Lock File Format

After running `zoi add package1 package2`, your `zoi.lock` might look like:

```json
{
  "version": "1",
  "packages": {
    "#zoidberg@main/ripgrep": "14.1.0",
    "#zoidberg@main/fd": "10.1.0"
  },
  "registries": {
    "zoidberg": "https://github.com/Zillowe/Zoidberg.git"
  },
  "#zoidberg": {
    "@main/ripgrep": {
      "version": "14.1.0",
      "integrity": "sha512-abc123...",
      "dependencies": [],
      "options_dependencies": [],
      "optionals_dependencies": []
    },
    "@main/fd": {
      "version": "10.1.0",
      "integrity": "sha512-def456...",
      "dependencies": [],
      "options_dependencies": [],
      "optionals_dependencies": []
    }
  }
}
```

## Updated zoi.yaml

With `--save` flag, your `zoi.yaml` will be updated:

```yaml
name: my-awesome-project
registry: https://github.com/Zillowe/Zoidberg.git

config:
  local: true

pkgs:
  - ripgrep
  - fd
  - bat
  - exa

commands:
  - cmd: dev
    run: npm run dev
```

## Alternative Commands

The following commands are equivalent:

```bash
# Adding packages
zoi add package-name
zoi install package-name --local
zoi install package-name --scope project

# Removing packages
zoi remove package-name
zoi uninstall package-name --local
zoi uninstall package-name --scope project
```

## Benefits

1. **Reproducibility**: Exact versions and integrity hashes ensure consistent environments
2. **Team Collaboration**: Everyone gets the same package versions
3. **Integrity Verification**: Detect corrupted or tampered packages
4. **Registry Tracking**: Know which registries are needed for the project
5. **Dependency Tracking**: See all dependencies for each package
