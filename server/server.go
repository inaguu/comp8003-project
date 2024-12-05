package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

type ServerInfo struct {
	Socket     *net.TCPListener
	Connection net.Conn

	Ip, Port string
	Command  string
	Output   string
}

const (
	INPUT_ERROR = "Usage: <ip address> <port_number>"
)

func exit() {
	fmt.Println("Exiting...")
	os.Exit(0)
}

func address(ip string, port string) string {
	unparsedIp := ip
	if ip == "localhost" {
		unparsedIp = "127.0.0.1"
	}

	parsedIp := net.ParseIP(unparsedIp)

	if parsedIp == nil {
		return ""
	}

	return net.JoinHostPort(parsedIp.String(), port)
}

func checkArgs(serverInfo *ServerInfo) {
	if address(serverInfo.Ip, serverInfo.Port) == "" {
		fmt.Printf("%s and %s is not a valid ip and port combination", serverInfo.Ip, serverInfo.Port)
		exit()
	}
}

func parseArgs(serverInfo *ServerInfo) {
	if len(os.Args) < 3 {
		fmt.Println(INPUT_ERROR)
		exit()
	}

	serverInfo.Ip = os.Args[1]
	serverInfo.Port = os.Args[2]

	checkArgs(serverInfo)

	fmt.Printf("The TCP server is %s\n", address(serverInfo.Ip, serverInfo.Port))
}

func bindSocket(serverInfo *ServerInfo) {
	s, err := net.ResolveTCPAddr("tcp", address(serverInfo.Ip, serverInfo.Port))
	if err != nil {
		fmt.Println(err)
		exit()
	}

	socket, err := net.ListenTCP("tcp", s)
	if err != nil {
		fmt.Println(err)
		socket.Close()
		exit()
	}

	serverInfo.Socket = socket
}

func receiveClient(serverInfo *ServerInfo) {
	for {
		conn, err := serverInfo.Socket.Accept()
		if err != nil {
			fmt.Println(err)
			exit()
		}

		serverInfo.Connection = conn
		handleConnection(serverInfo)
	}
}

func handleConnection(serverInfo *ServerInfo) {
	fmt.Println("\nConnected to:", serverInfo.Connection.RemoteAddr())

	data, err := bufio.NewReader(serverInfo.Connection).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		output := []byte(fmt.Sprintf("Error reading from client: %v\n", err))
		serverInfo.Output = string(output)
		sendOutput((serverInfo))
		return
	}

	fmt.Print("User executed the command: ", data)
	serverInfo.Command = data
	runCommand(serverInfo)
}

func runCommand(serverInfo *ServerInfo) {
	commandArr := strings.TrimSpace(serverInfo.Command)
	parts := strings.Fields(commandArr)

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		output = []byte(fmt.Sprintf("Error executing command: %v\n", err))
		serverInfo.Output = string(output)
		sendOutput((serverInfo))
		return
	}

	serverInfo.Output = string(output)
	sendOutput((serverInfo))
}

func sendOutput(serverInfo *ServerInfo) {
	_, err := serverInfo.Connection.Write([]byte(serverInfo.Output))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connection closed with:", serverInfo.Connection.RemoteAddr())
	serverInfo.Connection.Close()
}

func main() {
	serverInfo := ServerInfo{}
	parseArgs(&serverInfo)
	bindSocket(&serverInfo)
	receiveClient(&serverInfo)
}
