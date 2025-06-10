<div align="center">
    <img width="120" height="120" hspace="10" alt="ZDS Logo" src="https://codeberg.org/Zusty/ZDS/media/branch/main/img/zds.png"/>
    <h1>GCT</h1>
    <strong>An intelligent, AI-powered Git assistant.</strong>
    <br/>
    <p>Go beyond simple commits. Let GCT explain, create, and conversationally refine your code changes with the power of AI.</p>
<br/>
<a href="https://codeberg.org/Zillowe/ZFPL">
<img alt="ZFPL-1.0" src="https://codeberg.org/Zillowe/ZFPL/raw/branch/main/badges/1-0/dark.svg"/>
</a>
</div>

<hr/>

## 🚀 Getting Started

The first thing you should do after installation is set up your AI provider. GCT includes an interactive wizard to make this easy. In your project directory, run:

```sh
gct init
```

This will guide you through selecting an AI provider, setting your model, and adding your API key and guideline files to create a local `gct.yaml` configuration.

## ⚙️ Configuration (`gct.yaml`)

The `gct init` command creates a `gct.yaml` file that holds the configuration for all AI-related commands.

**Security Note:** This file contains your API key. The `gct init` command will automatically add `gct.yaml` to your `.gitignore` file to prevent accidentally committing secrets.

**Example `gct.yaml`:**

```yaml
name: "My Project"
provider: "OpenAI"
model: "gpt-4o"
api: "sk-..."
endpoint: "" # Only used for "OpenAI Compatible" provider
commits:
  guides:
    - ./docs/COMMIT_STYLEGUIDE.md
changelogs:
  guides:
    - ./docs/CHANGELOG_STYLE.md
```

- `provider`: The AI service you want to use. Supported: `"OpenAI"`, `"Google AI Studio"`, `"Anthropic"`, `"OpenRouter"`, `"OpenAI Compatible"`.
- `model`: The specific model name from your chosen provider (e.g. `gpt-4o`, `claude-3-sonnet-20240229`).
- `api`: Your secret API key from the provider's dashboard.
- `endpoint`: (Optional) The base URL for an "OpenAI Compatible" provider.
- `commits.guides`: A list of local files containing formatting rules for the `ai commit` command.
- `changelogs.guides`: A list of local files containing formatting rules for the `ai log` command.

## ✨ Commands

GCT is a command-line tool. Here are the available commands, grouped by category:

### Core Commands

| Command | Description |
| :--- | :--- |
| `gct init` | Interactively creates the `gct.yaml` config file in the current directory. |
| `gct version` | Shows GCT version information. |
| `gct about` | Displays details and information about GCT. |
| `gct update` | Checks for and applies updates to GCT itself. |
| `gct help` | Shows the detailed help message. |

### Manual Git Commands

| Command | Description |
| :--- | :--- |
| `gct commit` | Creates a new git commit using a TUI form. |
| `gct commit edit`| Edits the previous commit's message using the same interactive TUI. |

### AI Git Commands

| Command | Description |
| :--- | :--- |
| `gct ai commit [context]` | Generates and conversationally refines a commit message from staged changes. |
| `gct ai diff [args]` | Asks AI to explain a set of code changes in a readable format. |
| `gct ai log [args]` | Generates a user-facing changelog entry from code changes. |

The `ai commit` command allows you to provide extra context to the AI and then iteratively refine the generated message:

```sh
# Provide extra context for the initial generation
gct ai commit "This change was co-authored by Jane Doe"

# After generation, you get new options:
# > Press [c] to chat/change, [e] to edit, [Enter] to commit, [q] to quit:
```

The `ai diff` and `ai log` commands can be used in several ways:

- `gct ai diff`: Explains unstaged changes.
- `gct ai diff --staged`: Explains changes staged for the next commit.
- `gct ai diff <commit-hash>`: Explains a specific commit's changes.
- `gct ai diff <branch-name>`: Explains the differences between your current branch and the specified branch.

## 💾 Installation

You can either build it from source or install it using our installer scripts.

### Scripts

To install the latest version of GCT, run the appropriate command for your system:

```sh
# Linux / macOS
curl -fsSL https://zusty.codeberg.page/GCT/@main/app/install.sh | bash

# Windows (in PowerShell)
irm https://zusty.codeberg.page/GCT/@main/app/install.ps1 | iex
```

### Build from Source

To build GCT from source you need to have [`go`](https://go.dev) installed.

Then, clone the repository and run the build script:

```sh
# For Linux/macOS
./build/build-all.sh

# For Windows
./build/build-all.ps1
```

## 📚 Documentation

To get started with GCT please refer to the [GCT Wiki](https://codeberg.org/Zusty/GCT/wiki).

## Footer

GCT is developed by Zusty < Zillowe Foundation, part of the Zillowe Development Suite (ZDS).

### License

GCT is licensed under the [ZFPL-1.0](https://codeberg.org/Zillowe/ZFPL) (Zillowe Foundation Public License, Version 1.0).
