# LLMPack 📦

**LLMPack** is a blazing fast, zero-dependency CLI tool written in Go. It aggregates your codebase into a single, LLM-friendly context file (XML, Markdown, or ZIP), making it easy to feed entire projects to AI models like **ChatGPT (GPT-4o)**, **Claude 3.5**, or **Gemini**.

Designed for developers who are tired of manually copying and pasting files or struggling with `git archive`.

## 🚀 Key Features

* **Multi-Format Support:** Generate `XML` (best for prompting), `Markdown` (human-readable), or `ZIP` (for Code Interpreter).
* **Skeleton Mode:** A unique mode that parses AST (for Go) and strips function bodies, leaving only structures and interfaces. Reduces token usage by **up to 80%** when discussing architecture.
* **Cost Estimation:** Real-time token cost calculation for popular models (GPT-4o, Claude 3.5 Sonnet, Gemini 1.5).
* **Security Scanner:** Automatically detects and blocks sensitive data (API keys, `.env` files, private keys) to prevent accidental leakage.
* **Unix-way (Pipes):** Supports `STDIN`. You can pipe `git diff` or logs directly into LLMPack.
* **Smart Filtering:** Respects `.gitignore`, ignores binary files, and filters system directories (`.git`, `node_modules`).
* **Config Profiles:** Supports YAML configuration and profiles (e.g., different settings for `backend` vs `frontend`).

## 📦 Installation

### Option 1: Go Install (Recommended)

If you have Go (1.23+) installed:

```bash
go install github.com/dehimik/llmpack/cmd/llmpack@latest
````

### Option 2: Build from Source

```bash
git clone https://github.com/dehimik/llmpack.git
cd llmpack
go build -o llmpack cmd/llmpack/main.go

# Optional: Move to system path
sudo mv llmpack /usr/local/bin/
```

## 🛠 Usage

### Basic Usage

Pack the current directory into `context.xml` (default):

```bash
llmpack .
```

### Copy to Clipboard

Pack specific folders and copy the result directly to the clipboard:

```bash
llmpack internal/ cmd/ -c
```

### Skeleton Mode (Save Tokens)

Ideal for high-level architectural questions like "How do I refactor this module?". Leaves only signatures and types.

```bash
llmpack . -s
# Result: Compact context with "implementation hidden" bodies
```

### Cost Estimation

Check how much this context will cost for a specific model:

```bash
llmpack . --model claude-3-5-sonnet
# Output: Total Tokens: ~15400 ($0.04620 for claude-3-5-sonnet)
```

### Git Diff & Piping

Need an AI Code Review for your latest changes? Pipe the diff:

```bash
git diff main | llmpack --no-tree
```

### Output Formats

* **XML** (`-f xml`): Best structure for Claude/GPT prompts.
* **Markdown** (`-f md`): Readable format with code blocks.
* **Tree** (`-f tree`): Visual file tree only (no content).
* **Zip** (`-f zip`): Archive for file uploads.

## ⚙️ Configuration

You can create an `.llmpack.yaml` file in your project root or home directory:

```yaml
global:
  format: markdown
  tokens: true
  model_name: "gpt-4o"

profiles:
  backend:
    format: xml
    skeleton: true
  
ignore:
  - ".git"
  - "node_modules"
  - "images"
  - "*.lock"
```

Using a profile:

```bash
llmpack . -p backend
```

## 🚩 Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output file path (or `-` for stdout) | `context.xml` |
| `--format` | `-f` | Output format (`xml`, `markdown`, `zip`, `tree`) | `xml` |
| `--skeleton` | `-s` | **Skeleton Mode**: Strip function bodies | `false` |
| `--clipboard`| `-c` | Copy output to system clipboard | `false` |
| `--model` | `-m` | Model for cost estimation (`gpt-4o`, `claude-3-5`...) | `gpt-4o` |
| `--profile` | `-p` | Use settings from a specific config profile | - |
| `--config` | | Path to custom config file | `.llmpack.yaml` |
| `--tokens` | | Calculate token count | `true` |
| `--no-tree` | | Disable file tree header in output | `false` |
| `--no-security`| | Disable secrets detection (use with caution) | `false` |

## 🏗 Architecture

LLMPack is built with modularity and performance in mind:

* **Core:** Uses Go 1.23 iterators (`iter.Seq2`) for efficient file system traversal.
* **Streaming:** Utilizes `io.MultiWriter` to stream content to files and clipboard simultaneously without loading everything into RAM.
* **AST Parsing:** Uses `go/ast` for "Skeleton Mode" to ensure valid code structure after reduction.
* **Security:** Regex-based scanner to catch vulnerabilities before they enter the context.

## 📄 License

MIT License. See [LICENSE](https://www.google.com/search?q=LICENSE) for details.