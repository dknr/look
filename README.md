# look

A flicker-free alternative to GNU `watch`.

`look` runs a command repeatedly, displaying the output in a terminal.
Unlike GNU `watch`, which writes output directly to the screen (causing flicker),
`look` captures the entire output into a buffer and writes it all at once.

## Installation

```bash
go build -o look .
sudo cp look /usr/local/bin/
```

## Usage

```bash
look [flags] command [args...]
```

### Flags

- `-n duration`: Update interval (default 2s). Supports suffixes like `100ms`, `5s`, `1m`.
- `-e`: Exit if the command exits with a non-zero status code.
- `-g`: Exit if the command output changes.
- `-t`: Don't display the header title (e.g., "Every 1s: ...").
- `-h`: Show help.

## Example

```bash
look -n 1s "ping -c 1 8.8.8.8"
```

## License

ISC License - see LICENSE file for details.
