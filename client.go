package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func RunClient(prompt string) {
	path := socketPath()

	// Connect to the server
	conn, err := net.Dial("unix", path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "headless-askpass: Error connecting to server: %v\n", err)
		fmt.Fprintf(os.Stderr, "Make sure to start the server first.\n")
		os.Exit(1)
	}
	defer conn.Close()

	// Send the prompt to the server
	_, err = fmt.Fprintf(conn, "%s\n", prompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "headless-askpass: Error sending prompt to server: %v\n", err)
		os.Exit(1)
	}

	// Give user instructions
	fmt.Fprintf(os.Stderr, "Enter password in headless-askpass server\n")

	// Read the response from the server
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "headless-askpass: Error reading response from server: %v\n", err)
		os.Exit(1)
	}

	// Remove trailing newline and print the response
	fmt.Print(strings.TrimSuffix(response, "\n"))
}
