# so - Minimalistic Stack Overflow CLI

A simple, non-verbose command-line tool for searching Stack Overflow directly from your terminal.

## Features

- 🔍 Quick Stack Overflow searches from the command line
- 📖 Extract code snippets directly
- ⚡ Fast and minimal output (accessibility-first)
- 🎯 Progressive verbosity with flags
- 🔧 No dependencies on Python or TUI libraries

## Installation

```bash
cd so
go build -o so
```

Move the binary to your PATH:

```bash
# Linux/Mac
sudo mv so /usr/local/bin/

# Or add to your local bin
mv so ~/bin/
```

## Usage

### Basic Search (snippet only)

```bash
so how to exit vim
so golang reverse string
so python list comprehension
```

Output: Just the code snippet from the top answer.

### Verbosity Levels

**Level 0 (default)** - Just the code snippet:
```bash
so how to exit vim
```

**Level 1** - Snippet with context text before/after:
```bash
so -v 1 golang http get
```

**Level 2** - Full details (question, answer score, URL):
```bash
so -v 2 python async await
```

### List Search Results

Show all matching questions with answer counts:
```bash
so -l golang channels
```

Output:
```
[1] How do channels work in Go? (15 answers)
[2] Buffered vs unbuffered channels (8 answers)
[3] Channel select statement (12 answers)
...
```

### Select Specific Answer

```bash
so -a 2 golang reverse string    # Show 2nd answer instead of 1st
so -a 3 -v 1 python decorators   # 3rd answer with context
```

### Control Number of Results

```bash
so -n 5 -l javascript promises   # Show only 5 results
so -n 20 react hooks            # Search through 20 results
```

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-v` | Verbosity level (0=snippet, 1=context, 2=full) | 0 |
| `-l` | List search results with answer counts | false |
| `-a N` | Show answer number N | 1 |
| `-n N` | Number of search results to fetch | 10 |

## Examples

```bash
# Quick snippet
so how to reverse a string in go

# With explanation context
so -v 1 golang error handling

# See all results first
so -l python type hints

# Get the second answer with full details
so -a 2 -v 2 rust ownership

# Search through more results
so -n 20 javascript async patterns
```

## Design Philosophy

This tool is intentionally minimal and accessibility-focused:

- **Non-verbose by default** - Just gives you the answer
- **Progressive disclosure** - Add `-v` flags only when you need context
- **CLI-friendly** - Works well in scripts, easy to pipe or redirect
- **Fast** - Minimal dependencies, quick startup
- **Focused** - Does one thing well: get Stack Overflow answers quickly

## API

Uses the official Stack Exchange API v2.3. No rate limiting for reasonable use.

## Requirements

- Go 1.18+

## License

MIT