package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type ClientInfo struct {
	Connection *net.TCPConn
	IP         string
	Port       string
	Command    string
}

const (
	InputError  = "Usage: <ip address> <port_number>"
	ExitMessage = "Exiting..."
	BufferSize  = 1024
	Timeout     = 5 * time.Second
)

func exitWithMessage(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func constructAddress(ip, port string) string {
	if ip == "localhost" {
		ip = "127.0.0.1"
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ""
	}

	return net.JoinHostPort(parsedIP.String(), port)
}

func validateAddress(ip, port string) {
	if constructAddress(ip, port) == "" {
		exitWithMessage(fmt.Sprintf("Invalid IP and port combination: %s:%s", ip, port))
	}
}

func parseArguments(clientInfo *ClientInfo) {
	if len(os.Args) < 3 {
		exitWithMessage(InputError)
	}

	clientInfo.IP = os.Args[1]
	clientInfo.Port = os.Args[2]
	validateAddress(clientInfo.IP, clientInfo.Port)
}

func promptForCommand(clientInfo *ClientInfo) {
	fmt.Print("Enter command to execute: ")
	reader := bufio.NewReader(os.Stdin)
	command, err := reader.ReadString('\n')
	if err != nil {
		exitWithMessage("Failed to read command: " + err.Error())
	}

	command = strings.TrimSpace(command)
	if command == "" {
		exitWithMessage("No command provided.")
	}

	clientInfo.Command = command
}

func connectToServer(clientInfo *ClientInfo) {
	serverAddr, err := net.ResolveTCPAddr("tcp", constructAddress(clientInfo.IP, clientInfo.Port))
	if err != nil {
		exitWithMessage("Failed to resolve server address: " + err.Error())
	}

	connection, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		exitWithMessage("Failed to connect to server: " + err.Error())
	}

	clientInfo.Connection = connection
	fmt.Printf("Connected to server at %s\n", clientInfo.Connection.RemoteAddr().String())
}

func sendCommand(clientInfo *ClientInfo) {
	_, err := clientInfo.Connection.Write([]byte(clientInfo.Command + "\n"))
	if err != nil {
		exitWithMessage("Error sending command: " + err.Error())
	}
}

func receiveResponse(clientInfo *ClientInfo) {
	buffer := make([]byte, BufferSize)
	clientInfo.Connection.SetReadDeadline(time.Now().Add(Timeout))
	n, err := clientInfo.Connection.Read(buffer)
	if err != nil {
		exitWithMessage("Error reading response: " + err.Error())
	}

	fmt.Printf("Server response: %s\n", strings.TrimSpace(string(buffer[:n])))
}

func main() {
	clientInfo := &ClientInfo{}

	parseArguments(clientInfo)

	promptForCommand(clientInfo)

	connectToServer(clientInfo)
	defer clientInfo.Connection.Close()

	sendCommand(clientInfo)
	receiveResponse(clientInfo)
}
