# so - Minimalistic Stack Overflow CLI

A simple, non-verbose command-line tool for searching Stack Overflow directly from your terminal.

## Features

- 🔍 Quick Stack Overflow searches from the command line
- 📖 View questions and answers in your terminal
- 🌐 Open results in your browser
- ⚡ Fast and minimal output
- 🎯 No dependencies on Python or TUI libraries

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

### Basic Search

```bash
so how to reverse a string in go
so golang slice vs array
so python list comprehension
```

### Interactive Mode

After running a search, you'll see numbered results:

```
Found 10 results:

[1] How to reverse a string in Go?
    (5 answers)

[2] Reverse a string in Golang
    (3 answers)

...

Enter number to view (0 to exit, URL to open browser):
```

- Enter a number (e.g., `1`) to view the question and answers
- Enter `0` or press Enter to exit
- When viewing a question, you can open it in your browser

### Example Session

```bash
$ so golang http get request

Found 10 results:

[1] How to make an HTTP GET request in Go?
    (8 answers)

[2] Making HTTP requests in Golang
    (4 answers)

Enter number to view (0 to exit, URL to open browser): 1

================================================================================

How to make an HTTP GET request in Go?

--------------------------------------------------------------------------------

I'm trying to make a simple HTTP GET request in Go. What's the standard way
to do this using the net/http package?

================================================================================

Top 1 Answer(s):

[Answer 1]
--------------------------------------------------------------------------------

You can use http.Get() for simple GET requests:

resp, err := http.Get("https://example.com")
if err != nil {
    // handle error
}
defer resp.Body.Close()
body, err := ioutil.ReadAll(resp.Body)

[URL: https://stackoverflow.com/questions/...]

Open in browser? (y/N):
```

## Design Philosophy

This tool is intentionally minimal:

- **No verbose output** - Clean, readable results
- **CLI-friendly** - Works well in scripts and terminals
- **No TUI complexity** - Simple text output, easy to pipe or redirect
- **Fast** - Minimal dependencies, quick startup
- **Focused** - Does one thing well: search Stack Overflow

## Differences from Python Version

The original Python version (`sample.py`) included:

- Full TUI interface with urwid
- Code execution and error detection
- Complex scrolling widgets
- Real-time output capture

This Go version focuses on:

- Simple search functionality
- Clean terminal output
- Minimal dependencies
- Fast execution
- Easy to use in scripts

## Requirements

- Go 1.25+ (for `golang.org/x/net/html`)

## Building

```bash
go build -o so
```

## License

MIT