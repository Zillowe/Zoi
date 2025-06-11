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

## ✨ Features

- **Conversational AI Commits:** Generate a commit message and then "chat" with the AI to refine it until it's perfect.
- **Multi-Provider Support:** Works with over 10 AI providers, including OpenAI, Anthropic, Google (AI Studio & Vertex AI), Mistral, Amazon Bedrock, and any OpenAI-compatible endpoint.
- **AI-Generated Changelogs:** Automatically create user-facing changelogs from any set of git changes (`gct ai log`).
- **AI-Powered Diff Analysis:** Get a high-level explanation of any commit, branch, or staged changes (`gct ai diff`).
- **Guided Setup:** An interactive wizard (`gct init`) makes setup for any provider simple and fast.
- **Custom Guidelines:** Enforce project-specific styles for commits and changelogs by providing your own guide files.

## 🚀 Getting Started

The easiest way to get started is with the **model preset wizard**. It provides a curated list of popular, high-performance models and configures the provider for you. In your project directory, run:

```sh
gct init model
```

For a fully manual setup where you enter the provider and model name yourself, run `gct init`.

## ⚙️ Configuration (`gct.yaml`)

The `init` command creates a `gct.yaml` file that holds the configuration for all AI commands.

**Security Note:** This file contains your API key. The `gct init` command will automatically add `gct.yaml` to your `.gitignore` file to prevent accidentally committing secrets.

### Supported Providers

GCT supports a wide range of providers:
`Google AI Studio`, `Google Vertex AI`, `OpenAI`, `OpenAI Compatible`, `Azure OpenAI`, `Anthropic`, `OpenRouter`, `DeepSeek`, `Mistral`, `Alibaba`, `Hugging Face`, `Amazon Bedrock`, and `xAI`.

### Configuration Fields

- `provider`: The AI service you want to use.
- `model`: The specific model/deployment name from your chosen provider.
- `api`: Your secret API key.
- `commits.guides` & `changelogs.guides`: Lists of local files containing formatting rules.
- `endpoint`: (Optional) The base URL, only for the `"OpenAI Compatible"` provider.
- `gcp_project_id`, `gcp_region`: (Optional) Required only for `"Google Vertex AI"`.
- `aws_region`, `aws_access_key_id`, `aws_secret_access_key`: (Optional) Required only for `"Amazon Bedrock"`.
- `azure_resource_name`: (Optional) Required only for `"Azure OpenAI"`.

## 🤖 Model Recommendations

Choosing a model can be tough. Here are some recommended starting points for GCT's use case:

| Recommendation       | Model ID                             | Best For...                                         |
| :------------------- | :----------------------------------- | :-------------------------------------------------- |
| **Best Overall**     | `gpt-4o` (OpenAI)                    | Top-tier reasoning, speed, and instruction following. |
| **Best Balance**     | `claude-3-sonnet-20240229` (Anthropic)| Excellent performance at a great price point.         |
| **Fastest & Best Value** | `gemini-1.5-flash-latest` (Google)   | High-speed, low-cost tasks like `ai log` & chat.    |
| **Best Open Model**  | `meta-llama/llama-3-70b-instruct` (via OpenRouter) | State-of-the-art open-source performance.          |
| **Code Specialist**  | `deepseek-coder` (DeepSeek)          | Specifically trained on code for superior analysis. |

## ✨ Commands

GCT is a command-line tool. Here are the available commands, grouped by category:

### Core Commands

| Command | Description |
| :--- | :--- |
| `gct init model`| Starts a wizard with recommended models for easy setup. |
| `gct init` | Interactively creates a `gct.yaml` config file with manual input. |
| `gct version` | Shows GCT version information. |
| `gct help` | Shows the detailed help message. |

### Manual Git Commands

| Command | Description |
| :--- | :--- |
| `gct commit` | Creates a new git commit using an interactive TUI form. |
| `gct commit edit`| Edits the previous commit's message using the same TUI. |

### AI Git Commands

| Command | Description |
| :--- | :--- |
| `gct ai commit [context]` | Generates and conversationally refines a commit message from staged changes. |
| `gct ai diff [args]` | Asks AI to explain a set of code changes in a readable format. |
| `gct ai log [args]` | Generates a user-facing changelog entry from code changes. |

## 💾 Installation

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

```sh
# For Linux/macOS
./build/build-all.sh

# For Windows
./build/build-all.ps1
```

## 📚 Documentation

For more detailed guides, please refer to the **[[GCT Wiki]]**.

## Footer

GCT is developed by Zusty < Zillowe Foundation, part of the Zillowe Development Suite (ZDS).

### License

GCT is licensed under the [ZFPL-1.0](https://codeberg.org/Zillowe/ZFPL) (Zillowe Foundation Public License, Version 1.0).
