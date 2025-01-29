package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ClientInterface defines methods for the client.
type ClientInterface interface {
	Start() error
	Stop() error
	Send(msg string)
}

// Client represents a chat client with TView UI.
type Client struct {
	host           string
	port           int
	username       string
	conn           net.Conn
	mu             sync.RWMutex
	wg             sync.WaitGroup
	app            *tview.Application
	usersTable     *tview.Table
	chatBox        *tview.TextView
	inputField     *tview.InputField
	connectedUsers map[string]bool
	colorMap       map[string]string
	colors         []string
}

// NewClient creates a new client instance.
func NewClient(host string, port int, username string) ClientInterface {
	return &Client{
		host:           host,
		port:           port,
		username:       username,
		connectedUsers: make(map[string]bool),
		colorMap:       make(map[string]string),
		colors: []string{
			"green", "yellow", "blue", "purple", "cyan", "red",
			"orange", "lime", "pink", "skyblue", "violet", "gold",
			"silver", "magenta", "teal", "olive", "maroon", "navy",
		},
	}
}

// Start connects to the server, sets up the TUI, and runs it.
func (c *Client) Start() error {
	var err error
	c.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.host, c.port))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	// Send handshake: "USERNAME:<username>"
	_, err = c.conn.Write([]byte("USERNAME:" + c.username + "\n"))
	if err != nil {
		return fmt.Errorf("failed to send username: %v", err)
	}

	// Build the TUI
	c.app = tview.NewApplication().EnableMouse(true)

	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	contentFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	contentFlex.SetBorder(true)
	contentFlex.SetTitle(" Messenger ")

	// Create the users table
	c.usersTable = tview.NewTable().
		SetSelectable(false, false)

	c.usersTable.SetBorder(true).SetTitle(" Users ")

	// Create the chat box
	c.chatBox = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			c.app.Draw()
		})

	c.chatBox.SetBorder(true).SetTitle(fmt.Sprintf("[::b]Chat as %s", c.username))

	// Create the input field
	c.inputField = tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				text := c.inputField.GetText()
				c.inputField.SetText("")
				c.Send(text)
			}
		})

	c.inputField.SetBorder(true).SetTitle(" Type a message (Enter to send) ")

	contentFlex.AddItem(c.usersTable, 20, 1, false)
	contentFlex.AddItem(c.chatBox, 0, 3, true)
	rootFlex.AddItem(contentFlex, 0, 1, true)
	rootFlex.AddItem(c.inputField, 3, 1, false)

	c.app.SetRoot(rootFlex, true).
		SetFocus(c.inputField).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyUp:
				row, _ := c.chatBox.GetScrollOffset()
				c.chatBox.ScrollTo(row-1, 0)
				return nil
			case tcell.KeyDown:
				row, _ := c.chatBox.GetScrollOffset()
				c.chatBox.ScrollTo(row+1, 0)
				return nil
			case tcell.KeyPgUp:
				_, _, _, height := c.chatBox.GetInnerRect()
				row, _ := c.chatBox.GetScrollOffset()
				c.chatBox.ScrollTo(row-height, 0)
				return nil
			case tcell.KeyPgDn:
				_, _, _, height := c.chatBox.GetInnerRect()
				row, _ := c.chatBox.GetScrollOffset()
				c.chatBox.ScrollTo(row+height, 0)
				return nil
			}
			return event
		})
	// Start goroutine to read server messages
	c.wg.Add(1)
	go c.readMessages()

	// Run the TUI
	err = c.app.Run()
	if err != nil {
		return fmt.Errorf("error running TUI: %v", err)
	}

	// Wait for the read goroutine
	c.wg.Wait()

	return nil
}

// Stop closes the connection.
func (c *Client) Stop() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Send writes a message to the server.
func (c *Client) Send(msg string) {
	// Guard clause: no connection or empty message
	if c.conn == nil || strings.TrimSpace(msg) == "" {
		return
	}
	if c.conn != nil && strings.TrimSpace(msg) != "" {
		_, _ = c.conn.Write([]byte(msg + "\n"))
	}
}

// readMessages reads lines from the server and updates the TUI.
func (c *Client) readMessages() {
	defer c.wg.Done()

	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				c.printToChat("[red::b]Server closed the connection.")
			} else {
				c.printToChat(fmt.Sprintf("[red::b]Disconnected: %v", err))
			}
			c.app.QueueUpdateDraw(func() {
				c.app.Stop()
			})
			return
		}

		c.handleServerLine(strings.TrimSpace(line))
	}
}

func (c *Client) handleServerLine(line string) {
	if strings.HasPrefix(line, "USERS:") {
		// Handle user list
		userList := strings.TrimPrefix(line, "USERS:")
		users := strings.Split(userList, ",")

		c.mu.Lock()
		c.connectedUsers = make(map[string]bool)
		for _, user := range users {
			user = strings.TrimSpace(user)
			if user != "" {
				c.connectedUsers[user] = true
				if _, ok := c.colorMap[user]; !ok {
					idx := len(c.colorMap) % len(c.colors)
					c.colorMap[user] = c.colors[idx]
				}
			}
		}
		c.mu.Unlock()

		c.refreshUsers()
		return
	}

	// Otherwise, treat as broadcast: "[sender] message"
	parts := strings.SplitN(line, "]", 2)
	if len(parts) == 2 {
		sender := strings.TrimPrefix(parts[0], "[")
		text := parts[1]

		// Assign color to sender if needed
		c.mu.Lock()
		color, exists := c.colorMap[sender]
		if !exists {
			idx := len(c.colorMap) % len(c.colors)
			c.colorMap[sender] = c.colors[idx]
			color = c.colors[idx]
		}
		c.mu.Unlock()

		// Print message in chatBox
		c.app.QueueUpdateDraw(func() {
			timestamp := time.Now().Format("15:04:05")
			fmt.Fprintf(c.chatBox, "[%s][%s] [::b]%s[::-] %s\n",
				color, timestamp, sender, text)
		})
	}
}

// refreshUsers updates the users table in the TUI.
func (c *Client) refreshUsers() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.usersTable.Clear()
	row := 0
	c.usersTable.SetCell(row, 0,
		tview.NewTableCell("[::b]Connected Users").
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false),
	)
	row++

	for user := range c.connectedUsers {
		c.usersTable.SetCell(row, 0,
			tview.NewTableCell(user).
				SetTextColor(tcell.ColorWhite),
		)
		row++
	}
}

// printToChat is a helper to safely append text to the chat box.
func (c *Client) printToChat(msg string) {
	c.app.QueueUpdateDraw(func() {
		c.chatBox.ScrollToEnd()
		fmt.Fprintf(c.chatBox, "%s\n", msg)
	})
}
