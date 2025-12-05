# LLMPack 📦

**LLMPack** is a blazing fast, zero-dependency CLI tool written in Go. It aggregates your codebase into a single, LLM-friendly context file (XML, Markdown, or ZIP), making it easy to feed entire projects to AI models like **ChatGPT (GPT-4o)**, **Claude 3.5**, or **Gemini**.

Designed for developers who are tired of manually copying and pasting files or struggling with `git archive`.

## 🚀 Features

* **Multi-Format Support:** Generate `XML` (best for prompting), `Markdown` (readable), or `ZIP` (for Code Interpreter).
* **Smart Context:** Generates a concise file list or a visual ASCII tree (`-f tree`) to help LLMs understand project structure.
* **Token Counting:** Built-in `TikToken` integration instantly estimates the token cost of your context.
* **Smart Filtering:**
    * Automatically respects `.gitignore` rules.
    * Detects and skips binary files to save tokens.
    * Security filters for sensitive folders (`.git`, `.env`, keys).
* **Clipboard Integration:** Copy the result directly to your clipboard with `-c`.
* **Flexible Inputs:** Accepts specific files, multiple directories, or wildcards as arguments.
* **High Performance:** Built with **Go 1.25+** using iterators and stream processing for minimal memory footprint.

## 📦 Installation

### Option 1: Go Install (Recommended)
If you have Go installed:

```bash
go install https://github.com/dehimik/llmpack/cmd/llmpack@latest
````

### Option 2: Build from Source

```bash
git clone https://github.com/dehimik/llmpack.git
cd llmpack
go build -o llmpack cmd/llmpack/main.go

# Optional: Move to path
sudo mv llmpack /usr/local/bin/
```

## 🛠 Usage

### Basic Usage

Pack the current directory into an XML file (default):

```bash
llmpack .
# Creates context.xml
```

### Copy to Clipboard

Pack specific files and folders, then copy to clipboard immediately:

```bash
llmpack main.go internal/ pkg/utils.go -c
```

### Output Formats

**XML (Default)** — Best for structured prompts (Claude/GPT):

```bash
llmpack . -f xml -o context.xml
```

**Markdown** — Readable format with code blocks:

```bash
llmpack . -f markdown -o context.md
```

**ZIP Archive** — For uploading to ChatGPT Code Interpreter:

```bash
llmpack . -f zip -o project.zip
```

**Visual Tree** — Generate an ASCII directory tree (no file content):

```bash
llmpack . -f tree
```

### Advanced Options

**Disable Tree Header:**
By default, LLMPack adds a file list at the top of the context. To disable it:

```bash
llmpack . --no-tree
```

**Token Counting:**
Enabled by default. To disable:

```bash
llmpack . --tokens=false
```

## ⚙️ Configuration Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output file path (or `-` for stdout) | `context.xml` |
| `--format` | `-f` | Output format (`xml`, `markdown`, `zip`, `tree`) | `xml` |
| `--clipboard` | `-c` | Copy output to system clipboard | `false` |
| `--ignore-git` | | Respect `.gitignore` rules | `true` |
| `--tokens` | | Calculate and display token count | `true` |
| `--no-tree` | | Disable file tree/list in the output header | `false` |

## 🏗 Architecture

LLMPack is built using a modular architecture in **Go**:

* **Walker:** Uses Go 1.25 iterators (`iter.Seq2`) for efficient file system traversal.
* **Streaming:** Uses `io.MultiWriter` to stream content to files and clipboard simultaneously without loading everything into RAM.
* **Tokenizer:** Uses `tiktoken-go` for accurate token estimation.

## 🤝 Contributing

Contributions are welcome\!

1.  Fork the repository.
2.  Create a feature branch.
3.  Commit your changes.
4.  Open a Pull Request.

## 📄 License

MIT License. See [LICENSE](https://www.google.com/search?q=LICENSE) for details.

```
```