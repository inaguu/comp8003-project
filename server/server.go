package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ServerInfo struct {
	Socket     *net.TCPListener
	Connection net.Conn

	IP, Port string
	Command  string
	Output   string
}

const (
	InputError  = "Usage: <ip address> <port_number>"
	ExitMessage = "Exiting..."
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

func parseArguments(serverInfo *ServerInfo) {
	if len(os.Args) < 3 {
		exitWithMessage(InputError)
	}

	serverInfo.IP = os.Args[1]
	serverInfo.Port = os.Args[2]
	validateAddress(serverInfo.IP, serverInfo.Port)

	fmt.Printf("The TCP server is %s\n", constructAddress(serverInfo.IP, serverInfo.Port))
}

func bindSocket(serverInfo *ServerInfo) {
	serverAddr, err := net.ResolveTCPAddr("tcp", constructAddress(serverInfo.IP, serverInfo.Port))
	if err != nil {
		exitWithMessage("Failed to resolve server address: " + err.Error())
	}

	socket, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		exitWithMessage("Failed to bind socket: " + err.Error())
	}

	serverInfo.Socket = socket
	fmt.Println("Server is listening for connections...")
}

func receiveClients(serverInfo *ServerInfo) {
	for {
		conn, err := serverInfo.Socket.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		serverInfo.Connection = conn
		handleConnection(serverInfo)
	}
}

func handleConnection(serverInfo *ServerInfo) {
	fmt.Println("\nConnected to:", serverInfo.Connection.RemoteAddr())
	defer serverInfo.Connection.Close()

	serverInfo.Connection.SetReadDeadline(time.Now().Add(Timeout))
	data, err := bufio.NewReader(serverInfo.Connection).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from client:", err)
		serverInfo.Output = fmt.Sprintf("Error reading from client: %v\n", err)
		sendOutput(serverInfo)
		return
	}

	serverInfo.Command = strings.TrimSpace(data)
	fmt.Println("User executed the command:", serverInfo.Command)
	runCommand(serverInfo)
}

func runCommand(serverInfo *ServerInfo) {
	parts := strings.Fields(serverInfo.Command)
	if len(parts) == 0 {
		serverInfo.Output = "No command provided."
		sendOutput(serverInfo)
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		serverInfo.Output = fmt.Sprintf("Error executing command: %v\n", err)
		sendOutput(serverInfo)
		return
	}

	serverInfo.Output = string(output)
	sendOutput(serverInfo)
}

func sendOutput(serverInfo *ServerInfo) {
	_, err := serverInfo.Connection.Write([]byte(serverInfo.Output + "\n"))
	if err != nil {
		fmt.Println("Error sending output to client:", err)
	}

	fmt.Println("Connection closed with:", serverInfo.Connection.RemoteAddr())
}

func main() {
	serverInfo := &ServerInfo{}

	parseArguments(serverInfo)

	bindSocket(serverInfo)

	receiveClients(serverInfo)
}
