package main

import (
	"fmt"
	"net"
	"os"
)

type ServerInfo struct {
	Socket        *net.TCPListener
	ClientAddress net.Conn

	Ip, Port string
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
	if len(os.Args) < 2 {
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

	connection, err := net.ListenTCP("tcp", s)
	if err != nil {
		fmt.Println(err)
		exit()
	}

	serverInfo.Socket = connection
}

func receiveClient(serverInfo *ServerInfo) {
	for {
		conn, err := serverInfo.Socket.Accept()
		if err != nil {
			fmt.Println(err)
			exit()
		}

		serverInfo.ClientAddress = conn
		handleConnection(serverInfo)
	}
}

func handleConnection(serverInfo *ServerInfo) {
	defer serverInfo.ClientAddress.Close()
	fmt.Println("Connected to:", serverInfo.ClientAddress.RemoteAddr())

	buffer := make([]byte, 1024)
	_, err := serverInfo.ClientAddress.Read(buffer)
	if err != nil {
		fmt.Println(err)
		exit()
	}

	fmt.Println(string(buffer))

	output := []byte("hello")

	_, err = serverInfo.ClientAddress.Write(output)
	if err != nil {
		fmt.Println(err)
		exit()
	}

	serverInfo.ClientAddress.Close()
}

func main() {
	serverInfo := ServerInfo{}
	parseArgs(&serverInfo)
	bindSocket(&serverInfo)
	receiveClient(&serverInfo)
}
