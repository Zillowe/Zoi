<div align="center">
    <img width="120" height="120" hspace="10" alt="ZDS Logo" src="https://codeberg.org/Zusty/ZDS/media/branch/main/img/zds.png"/>
    <h1>GCT</h1>
    An intelligent, AI-powered Git assistant
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

This will guide you through selecting an AI provider, setting your model, and adding your API key to create a local `gct.yaml` configuration file.

## ⚙️ Configuration

The `gct init` command creates a `gct.yaml` file in your project's root directory. This file holds the configuration for all AI-related commands.

**Example `gct.yaml`:**

```yaml
name: "GCT"
guides:
  - ./etc/Commits.md
provider: "OpenRouter" # Supported: "OpenAI", "Google AI Studio", "Anthropic", "OpenRouter"
model: "google/gemma-2-9b-it" # The model name from your chosen provider
api: "sk-or-v1-abc...123" # Your API Key
```

- `provider`: The AI service you want to use.
- `model`: The specific model name from that provider (e.g. `gpt-4o`, `claude-3-haiku-20240307`).
- `api`: Your secret API key from the provider's dashboard.
- `guides`: A list of local Markdown or text files that provide style guidelines for the AI to follow.

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
| `gct commit edit` | Edits the previous commit's message using the same interactive TUI. |

### AI Git Commands

| Command | Description |
| :--- | :--- |
| `gct ai commit` | Uses AI to automatically generate a commit message from your staged changes. |
| `gct ai diff [args]`| Asks AI to explain a set of code changes in a readable format. |

The `gct ai diff` command can be used in several ways:

- `gct ai diff`: Explains unstaged changes in your working directory.
- `gct ai diff --staged`: Explains changes that are staged for the next commit.
- `gct ai diff <commit-hash>`: Explains the changes introduced by a specific commit.
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
./build/build.sh

# For Windows
./build/build.ps1
```

## 📚 Documentation

To get started with GCT please refer to the [GCT Wiki](https://codeberg.org/Zusty/GCT/wiki).

## Footer

GCT is developed by Zusty < Zillowe Foundation, part of the Zillowe Development Suite (ZDS).

### License

GCT is licensed under the [ZFPL-1.0](https://codeberg.org/Zillowe/ZFPL) (Zillowe Foundation Public License, Version 1.0).
