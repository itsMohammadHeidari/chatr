package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
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
		conn.Close()
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
	for c := range s.clients {
		_, err := c.Write([]byte(fullMsg))
		if err != nil {
			log.Printf("error writing to %s: %v", c.RemoteAddr(), err)
		}
	}
}

// Add new method to broadcast user list
func (s *Server) broadcastUserList() {
	s.mu.Lock()
	var users []string
	for _, user := range s.clients {
		users = append(users, user)
	}
	s.mu.Unlock()

	userListMsg := "USERS:" + strings.Join(users, ",") + "\n"
	s.mu.Lock()
	for conn := range s.clients {
		conn.Write([]byte(userListMsg))
	}
	s.mu.Unlock()
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		s.mu.Lock()
		username := s.clients[conn]
		delete(s.clients, conn)
		s.mu.Unlock()
		// Broadcast user left and updated user list
		s.Broadcast("SERVER", fmt.Sprintf("%s has left the chat.", username))
		s.broadcastUserList() // Broadcast updated user list
		log.Printf("Client disconnected: %s", conn.RemoteAddr())
		conn.Close()
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

	username := strings.TrimPrefix(line, "USERNAME:")
	username = strings.TrimSpace(username)
	if username == "" {
		username = "Guest"
	}

	// Store the client and broadcast updated user list
	s.mu.Lock()
	s.clients[conn] = username
	s.mu.Unlock()

	// Broadcast updated user list to all clients
	s.broadcastUserList()

	// Broadcast that the user has joined
	s.Broadcast("SERVER", fmt.Sprintf("%s has joined the chat.", username))
	log.Printf("Client %s joined as '%s'", conn.RemoteAddr(), username)

	// Read messages in a loop and broadcast them
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		s.Broadcast(username, strings.TrimSpace(msg))
	}
}
