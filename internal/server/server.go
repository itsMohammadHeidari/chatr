package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// ServerInterface defines methods for the server.
type ServerInterface interface {
	Start() error
	Stop() error
	Broadcast(sender string, msg string)
}

// Server tracks all clients and their usernames.
type Server struct {
	host    string
	port    int
	mu      sync.Mutex
	clients map[net.Conn]string // connection -> username
}

// NewServer creates a new server instance.
func NewServer(host string, port int) ServerInterface {
	return &Server{
		host:    host,
		port:    port,
		clients: make(map[net.Conn]string),
	}
}

// Start listens for new connections and handles them concurrently.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return fmt.Errorf("failed to listen on %s:%d: %v", s.host, s.port, err)
	}
	defer listener.Close()

	log.Printf("Server started on %s:%d", s.host, s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// Stop closes all connections and stops the server.
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.clients {
		_ = conn.Close()
		delete(s.clients, conn)
	}
	log.Println("Server stopped.")
	return nil
}

// Broadcast sends a message to all connected clients.
func (s *Server) Broadcast(sender string, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fullMsg := fmt.Sprintf("[%s] %s\n", sender, msg)
	var failedConns []net.Conn

	for c := range s.clients {
		if _, err := c.Write([]byte(fullMsg)); err != nil {
			log.Printf("error writing to %s: %v", c.RemoteAddr(), err)
			failedConns = append(failedConns, c)
		}
	}
	// Remove stale connections
	for _, fc := range failedConns {
		_ = fc.Close()
		delete(s.clients, fc)
	}
}

// broadcastUserList sends the updated user list to all clients.
func (s *Server) broadcastUserList() {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []string
	for _, user := range s.clients {
		users = append(users, user)
	}
	userListMsg := "USERS:" + strings.Join(users, ",") + "\n"

	var failedConns []net.Conn
	for conn := range s.clients {
		if _, err := conn.Write([]byte(userListMsg)); err != nil {
			log.Printf("error writing user list to %s: %v", conn.RemoteAddr(), err)
			failedConns = append(failedConns, conn)
		}
	}
	// Remove stale connections
	for _, fc := range failedConns {
		_ = fc.Close()
		delete(s.clients, fc)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	// Cleanup logic when the function returns
	defer func() {
		s.mu.Lock()
		username := s.clients[conn]
		delete(s.clients, conn)
		s.mu.Unlock()

		s.Broadcast("SERVER", fmt.Sprintf("%s has left the chat.", username))
		s.broadcastUserList()
		log.Printf("Client disconnected: %s", conn.RemoteAddr())
		_ = conn.Close()
	}()

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("error reading username: %v", err)
		return
	}
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "USERNAME:") {
		log.Printf("invalid username handshake from %s", conn.RemoteAddr())
		conn.Write([]byte("Invalid handshake. Expected USERNAME:<yourname>\n"))
		return
	}

	username := strings.TrimSpace(strings.TrimPrefix(line, "USERNAME:"))
	if username == "" {
		username = "Guest"
	}

	// Add client to the map
	s.mu.Lock()
	s.clients[conn] = username
	s.mu.Unlock()

	s.broadcastUserList()
	s.Broadcast("SERVER", fmt.Sprintf("%s has joined the chat.", username))
	log.Printf("Client %s joined as '%s'", conn.RemoteAddr(), username)

	// Read messages in a loop, applying read deadlines
	for {
		conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
		msg, err := reader.ReadString('\n')
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("Client %s timed out", conn.RemoteAddr())
			} else {
				log.Printf("error reading from %s: %v", conn.RemoteAddr(), err)
			}
			return
		}
		conn.SetReadDeadline(time.Time{}) // reset read deadline after success

		s.Broadcast(username, strings.TrimSpace(msg))
	}
}
