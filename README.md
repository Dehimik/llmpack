# LLMPack 📦

**LLMPack** is a blazing fast, zero-dependency CLI tool written in Go. It aggregates your codebase into a single, LLM-friendly context file (XML, Markdown, or ZIP), making it easy to feed entire projects to AI models like **ChatGPT (GPT-4o)**, **Claude 3.5**, or **Gemini**.

Designed for developers who are tired of manually copying and pasting files or struggling with `git archive`.

## 🚀 Key Features

* **Multi-Format Support:** Generate `XML` (best for prompting), `Markdown` (human-readable), `ZIP` (for Code Interpreter), or **`llms.txt`** (standardized project index).
* **Universal Skeleton Mode:** Strips function/method bodies while preserving signatures. Supports **Go** (via AST), **Python** (via indentation), and **20+ C-family languages** (JS, TS, Java, C++, Rust, etc.).
* **Semantic Indexing:** List all symbols, extract specific implementations, or search for symbols across the project.
* **Auto-Metadata Extraction:** Automatically detects project name and summary from your `README.md` to enrich `llms.txt` and `Markdown` outputs.
* **Cost Estimation:** Real-time token cost calculation for popular models (GPT-4o, Claude 3.5 Sonnet, Gemini 1.5).
* **Security Scanner:** Automatically detects and blocks sensitive data (API keys, `.env` files, private keys) to prevent accidental leakage.
* **MCP Compatible:** Works as a toolset for AI agents (Claude Code, Gemini CLI, etc.).
* **Smart Filtering:** Respects `.gitignore`, ignores binary files, and filters system directories.

## 📦 Installation

### Option 1: Using Makefile (Recommended)

If you have Go (1.23+) and `make` installed:

```bash
git clone https://github.com/dehimik/llmpack.git
cd llmpack
make install
```
This builds the binary and moves it to `/usr/local/bin/`.

### Option 2: Go Install

```bash
go install github.com/dehimik/llmpack/cmd/llmpack@latest
```

## 🛠 Usage

### Semantic Features (New!)

*   **List Symbols:** Get a compact index of all functions, classes, and methods.
    ```bash
    llmpack . --symbols
    ```
*   **Extract Implementation:** Show full code for a specific symbol and skeletonize everything else.
    ```bash
    llmpack . --implementation MyFunction
    ```
*   **Smart Search & Focus:** Find files containing a symbol and return them with that symbol's body expanded.
    ```bash
    llmpack . --find MyMethod --focus
    ```

### llms.txt Support (Answer.AI Spec)

Generate a standardized project index for LLMs:
```bash
llmpack pack -f llms-txt .
```
This will automatically create an `llms.txt` file with your project's H1 title and summary extracted from `README.md`.

### Shell Autocompletion

Enable intelligent autocompletion for commands and flags (supports bash, zsh, fish, powershell).

**For Zsh:**
```bash
echo 'source <(llmpack completion zsh)' >> ~/.zshrc
```

**For Bash:**
```bash
source <(llmpack completion bash)
```

### AI Agents & MCP Support 🤖

LLMPack supports the **Model Context Protocol (MCP)**. Add it to your AI agent to give it "superpowers" over your codebase.

**Tools exposed via MCP:**
- `list_symbols`: Browse project architecture.
- `get_code`: Pull specific implementation of a symbol.
- `search`: Find where a symbol is defined.
- `scan_security`: Proactively check for secrets.
- `estimate_cost`: Calculate token budget.
- `get_tree`: Visualize project structure.

**Claude Desktop Configuration:**
Add this to your `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "llmpack": {
      "command": "llmpack",
      "args": ["mcp"]
    }
  }
}
```

## 🚩 Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output file path (or `-` for stdout) | `context.xml` |
| `--format` | `-f` | Output format (`xml`, `markdown`, `zip`, `tree`, `llms-txt`) | `xml` |
| `--skeleton` | `-s` | **Skeleton Mode**: Strip function bodies | `false` |
| `--symbols` | | List all symbols in AI-friendly format | `false` |
| `--implementation` | | Extract full implementation of a symbol | - |
| `--find` | | Find files containing a specific symbol | - |
| `--focus` | | In find mode, return only implementation + skeleton | `false` |
| `--clipboard`| `-c` | Copy output to system clipboard | `false` |
| `--model` | `-m` | Model for cost estimation (`gpt-4o`, `claude-3-5`...) | `gpt-4o` |
| `--tokens` | | Calculate token count | `true` |
| `--no-tree` | | Disable file tree header in output | `false` |
| `--no-security`| | Disable secrets detection | `false` |

## ⚙️ Configuration

Create an `.llmpack.yaml` in your project root:

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
  - "*.lock"
```

## 🏗 Architecture

LLMPack is built for speed and AI-compatibility:
* **Core:** Go 1.23 iterators for high-performance FS traversal.
* **Streaming:** `io.MultiWriter` for efficient data flow.
* **Heuristics:** Multi-language symbol extraction without heavy AST parsers.

## 📄 License

MIT License.
