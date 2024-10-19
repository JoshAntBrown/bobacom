# Bobacom

Bobacom is a simple terminal-based serial communication monitor written in Go with a TUI using [Bubble Tea](https://github.com/charmbracelet/bubbletea/). It allows users to read and write to serial devices in real-time.

## Usage

```bash

go run main.go -b 115200 /dev/tty.usbmodem1101
```

### Command-line Options
- `-b` Baud rate, default (9600)

## Example Preview
Using bobacom to communicate with an Arduino over serial.
![Screenshot 2024-10-19 at 17 21 55](https://github.com/user-attachments/assets/31f8fab1-52e7-4b43-9089-3b5d593929c9)
