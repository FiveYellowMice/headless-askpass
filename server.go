package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

type Request struct {
	conn net.Conn
}

type Server struct {
	queue chan *Request
}

func (self *Server) handleConnection(conn net.Conn) {
	// Queue the request
	request := &Request{
		conn: conn,
	}

	select {
	case self.queue <- request:
		// Request queued successfully
	default:
		// Queue is full
		fmt.Fprintf(os.Stderr, "Queue full, rejecting request\n")
		conn.Close()
	}
}

func (self *Server) processRequest() {
	request := <- self.queue
	defer request.conn.Close()

	// Read the prompt from the client
	reader := bufio.NewReader(request.conn)
	prompt, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "Error reading from connection: %v\n", err)
		return
	}

	// Remove trailing newline
	prompt = strings.TrimSuffix(prompt, "\n")

	// Use provided prompt or default to "Password: "
	if prompt == "" {
		prompt = "Password: "
	}

	// Print the prompt to the terminal
	fmt.Print(prompt)

	// Get the file descriptor for stdin
	fd := int(os.Stdin.Fd())

	// Flush any pending input typed before the prompt was shown
	if _, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.TCFLSH, unix.TCIFLUSH); errno != 0 {
		fmt.Println()
		fmt.Fprintf(os.Stderr, "Failed to flush terminal input: %v\n", errno)
		return
	}

	// Read password using term.ReadPassword
	response, err := term.ReadPassword(fd)
	if err != nil {
		fmt.Println()
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		return
	}

	// Send response back to client
	_, err = fmt.Fprintf(request.conn, "%s\n", string(response))
	if err != nil {
		fmt.Println()
		fmt.Fprintf(os.Stderr, "Error sending response: %v\n", err)
		return
	}

	// Give success feedback
	fmt.Print("[sent]\n")
}

func RunServer() {
	for _, file := range []*os.File{os.Stdin, os.Stdout} {
		// Check if stdin/stdout is connected to a terminal
		if !term.IsTerminal(int(file.Fd())) {
			fmt.Fprintf(os.Stderr, "Server must be run in a terminal\n")
			os.Exit(1)
		}

		// Stat the terminal file
		stat, err := file.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to stat terminal: %v\n", err)
		}

		// Check that group and others don't have any permissions
		permissions := stat.Mode().Perm()
		if permissions & 0077 != 0 {
			fmt.Fprintf(os.Stderr, "Terminal must have only read/write permissions by the owner\n")
		}

	}

	path := socketPath()

	// Remove existing socket file if it exists
	os.Remove(path)

	// Listen on socket
	listener, err := net.Listen("unix", path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to listen on socket: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		listener.Close()
		os.Remove(path)
	}()

	// Restrict socket to owner only
	if err := os.Chmod(path, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set socket permissions: %v\n", err)
		os.Exit(1)
	}

	// Initialize server object
	server := Server{
		queue: make(chan *Request, 10), // Buffer for up to 10 pending requests
	}

	fmt.Fprintf(os.Stderr, "Server listening on %s\n", path)

	// Start the request processor
	go func() {
		for {
			server.processRequest()
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
			continue
		}

		go server.handleConnection(conn)
	}
}
