package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type ClientInfo struct {
	Connection *net.TCPConn

	Ip, Port string
	Command  string
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

func checkArgs(clientInfo *ClientInfo) {
	if address(clientInfo.Ip, clientInfo.Port) == "" {
		fmt.Printf("%s and %s is not a valid ip and port combination", clientInfo.Ip, clientInfo.Port)
		exit()
	}
}

func parseArgs(clientInfo *ClientInfo) {
	if len(os.Args) < 3 {
		fmt.Println(INPUT_ERROR)
		exit()
	}

	clientInfo.Ip = os.Args[1]
	clientInfo.Port = os.Args[2]

	checkArgs(clientInfo)

}

func getCommand(clientInfo *ClientInfo) {
	fmt.Print("Enter command to execute: ")
	reader := bufio.NewReader(os.Stdin)
	command, _ := reader.ReadString('\n')

	checkCommand := strings.TrimSuffix(command, "\n")

	if checkCommand == "" {
		fmt.Println("No command provided.")
		exit()
	}

	clientInfo.Command = command
}

func bindAndConnect(clientInfo *ClientInfo) {
	s, _ := net.ResolveTCPAddr("tcp", address(clientInfo.Ip, clientInfo.Port))
	c, err := net.DialTCP("tcp", nil, s)
	if err != nil {
		fmt.Println(err)
		c.Close()
		exit()
	}

	clientInfo.Connection = c

	fmt.Printf("The TCP server is %s\n", clientInfo.Connection.RemoteAddr().String())
}

func sendCommand(clientInfo *ClientInfo) {
	_, err := clientInfo.Connection.Write([]byte(clientInfo.Command))
	if err != nil {
		fmt.Println("Error sending command:", err)
		clientInfo.Connection.Close()
		exit()
	}
}

func getResponse(clientInfo *ClientInfo) {
	buffer := make([]byte, 1024)
	_, err := clientInfo.Connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading output:", err)
		clientInfo.Connection.Close()
		exit()
	}

	fmt.Println(string(buffer))
	clientInfo.Connection.Close()
}

func main() {
	clientInfo := ClientInfo{}
	parseArgs(&clientInfo)
	getCommand(&clientInfo)
	bindAndConnect(&clientInfo)
	sendCommand(&clientInfo)
	getResponse(&clientInfo)
}
