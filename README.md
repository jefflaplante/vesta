# vesta

A command-line tool for sending formatted messages to [Vestaboard](https://www.vestaboard.com/) Note and Flagship devices.

Features auto-wrapping, centering, color/symbol escape codes, and a dry-run preview mode.

## Setup

### Prerequisites

- Go 1.23+
- A Vestaboard with a [Read/Write API key](https://www.vestaboard.com/documentation) (cloud) or a Local API key (local)

### Install

```sh
# Clone and build
git clone https://github.com/jeff/vesta.git
cd vesta
make build

# Or install directly
go install .
```

### Configure

Set your API token securely (hidden input, avoids shell history):

```sh
vesta config set token
```

Or provide it directly:

```sh
vesta config set token YOUR_API_KEY
```

Or use an environment variable:

```sh
export VESTABOARD_API_KEY=YOUR_API_KEY
```

Set your device type (defaults to `note`):

```sh
vesta config set device flagship
```

Verify your configuration:

```sh
vesta config show
```

### Local API Setup (Optional)

The local API allows direct communication with your Vestaboard over your local network, bypassing the cloud. This provides lower latency and works without internet connectivity.

1. Get your Local API key from your Vestaboard's settings
2. Configure the local URL and token:

```sh
vesta config set local-url 192.168.1.100:7000  # Your board's IP
vesta config set local-token                    # Prompts for token
```

3. Use local mode:

```sh
# Per-command with --local flag
vesta send --local "Hello"
vesta read --local

# Or set as default
vesta config set api-mode local
vesta send "Hello"  # Now uses local API by default
```

Environment variables are also supported:

```sh
export VESTABOARD_LOCAL_URL=192.168.1.100:7000
export VESTABOARD_LOCAL_API_KEY=YOUR_LOCAL_KEY
export VESTABOARD_API_MODE=local
```

## Usage

### Send a message

```sh
vesta send "Hello World"
```

Messages auto-wrap to fit the board. Use `\n` for explicit line breaks:

```sh
vesta send "Line one\nLine two"
```

### Center a message

```sh
vesta send -c "Centered"
```

### Preview without sending

```sh
vesta send --dry-run "Testing 123"
```

```
Character array:
Row 0: [20 5 19 20 9 14 7 0 27 28 29 0 0 0 0]
Row 1: [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
Row 2: [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]

Preview:
┌───────────────┐
│TESTING 123    │
│               │
│               │
└───────────────┘
```

### Read from stdin (scripting)

```sh
echo "Hello" | vesta send -
cat message.txt | vesta send -c -
```

### Read the board

```sh
vesta read
```

### Escape codes

Use escape codes in your message for colors and symbols:

| Code | Description |
|------|-------------|
| `{red}` | Red tile |
| `{orange}` | Orange tile |
| `{yellow}` | Yellow tile |
| `{green}` | Green tile |
| `{blue}` | Blue tile |
| `{violet}` | Violet tile |
| `{white}` | White tile |
| `{black}` | Black tile |
| `{filled}` | Filled tile |
| `{deg}` | Degree symbol (Flagship only) |
| `{<3}` or `<3` | Heart (Note only) |
| `{0}`-`{71}` | Raw character code |

```sh
vesta send "{red}{orange}{yellow}{green}{blue}{violet} Rainbow"
vesta send "I <3 Go"
```

### Global flags

| Flag | Description |
|------|-------------|
| `--device note\|flagship` | Override device type for this command |
| `-l`, `--local` | Use local API instead of cloud |
| `-v`, `--verbose` | Show detailed error information |

## Device types

| Device | Rows | Columns |
|--------|------|---------|
| Note | 3 | 15 |
| Flagship | 6 | 22 |

### Version information

```sh
vesta version
```

## Development

```sh
make test       # Run tests
make test-v     # Run tests with verbose output
make lint       # Format and vet
make clean      # Remove binary
make all        # Clean, build, and test
```

## License

MIT
