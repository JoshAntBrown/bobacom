package main

import (
  "log"
  "flag"
  "fmt"
  "strings"
  "time"
  "go.bug.st/serial"
  "github.com/charmbracelet/bubbles/textinput"
  "github.com/charmbracelet/bubbles/viewport"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
  inputStyle = func() lipgloss.Style {
    b := lipgloss.RoundedBorder()
    return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
  }()
)

type serialMsg string

type model struct {
  ready bool
  content string
  port serial.Port
  portName string
  baudRate int
  viewport viewport.Model
  message textinput.Model
}

func New(port serial.Port, portName string, baudRate int) *model {
  message := textinput.New()
  message.Placeholder = "Type something..."
  message.Focus()

  return &model{
    message: message,
    port: port,
    portName: portName,
    baudRate: baudRate,
  }
}

func (m model) Init() tea.Cmd {
  return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var (
    cmds []tea.Cmd
    cmd tea.Cmd
  )

  switch msg := msg.(type) {
  case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
    verticalMarginHeight := headerHeight + footerHeight

    if !m.ready {
      m.viewport = viewport.New(msg.Width, msg.Height - verticalMarginHeight)
      m.viewport.YPosition = headerHeight
      m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
      m.viewport.SetContent(m.content)
      m.viewport.YPosition = headerHeight + 1
      m.ready = true
    } else {
      m.viewport.Width = msg.Width
      m.viewport.Height = msg.Height - verticalMarginHeight
    }

    if useHighPerformanceRenderer {
      cmds = append(cmds, viewport.Sync(m.viewport))
    }

  case tea.KeyMsg:
    switch msg.String() {
    case "ctrl+c":
      return m, tea.Quit

    case "enter":
      m.sendToSerial(m.message.Value() + "\n")
      m.message.SetValue("")
      return m, nil
    }

  case serialMsg:
    m.addContent(string(msg))
    return m, nil
  }

  m.viewport, cmd = m.viewport.Update(msg)
  cmds = append(cmds, cmd)

  m.message, cmd = m.message.Update(msg)
  cmds = append(cmds, cmd)
  
  return m, tea.Batch(cmds...)
}

func (m *model) sendToSerial(msg string) {
  _, err := m.port.Write([]byte(msg))
  if err != nil {
    log.Fatalf("Failed to write to serial port: %v", err)
  } else {
    m.addContent(msg)
  }
}

func (m *model) addContent(newContent string) {
  m.content += newContent
  m.viewport.SetContent(m.content)
  m.viewport.GotoBottom()
}

func (m model) View() string {
  if !m.ready {
    return "\n  Initializing..."
  }

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
  title := titleStyle.Render("Bobacom – " + m.portName + " @ " + fmt.Sprintf("%d", m.baudRate) + " baud")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
  contentWidth := m.viewport.Width - inputStyle.GetHorizontalFrameSize()
  return inputStyle.Width(contentWidth).Render(m.message.View())
}

func readSerial(port serial.Port, p *tea.Program, readInterval time.Duration) {
  buf := make([]byte, 256)

  for {
    n, err := port.Read(buf)
    if err != nil {
      log.Fatalf("Failed to read from serial port: %v", err)
      continue
    }

    if n > 0 {
      p.Send(serialMsg(string(buf[:n])))
    }

    time.Sleep(readInterval)
  }
}

func main() {
  baudRate := flag.Int("b", 9600, "Set baud rate for serial communication")
  readInterval := 100 * time.Millisecond
  flag.Parse()

  args := flag.Args()
  portName := args[0]

  mode := &serial.Mode{
    BaudRate: *baudRate,
  }
  port, err := serial.Open(portName, mode)
  if err != nil {
    log.Fatal(err)
  }
  defer port.Close()

  m := New(port, portName, *baudRate)

  p := tea.NewProgram(
    m,
    tea.WithAltScreen(),
    tea.WithMouseCellMotion(),
  )

  go func() {
    for {
      readSerial(port, p, readInterval)
    }
  }()

  if _, err := p.Run(); err != nil {
    log.Fatal(err)
  }
}
