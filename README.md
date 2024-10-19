# Bobacom

Bobacom is a simple terminal-based serial communication monitor written in Go with a TUI using [Bubble Tea](https://github.com/charmbracelet/bubbletea/). It allows users to read and write to serial devices in real-time.

## Usage

```bash

go run main.go -b 115200 /dev/tty.usbmodem1101
```

### Command-line Options
- `-b` Baud rate, default (9600)
