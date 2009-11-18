package main

import (
	"net";
	"log";
	"strings";
	"fmt";
	"encoding/binary";
)

type TClient struct {
	local		net.TCPAddr;
	forserver	net.TCPAddr;
	msg			string;
}

func ProcessConn(conn *net.TCPConn, client chan TClient) {
	defer conn.Close();
	log.Stdout("connected\n");

	//get client's ip and port
	conn.Write(strings.Bytes("hi"));
	ip16 := make([]byte, 16);
	conn.Read(ip16);
	ip := net.IP(ip16);
	var port uint32;
	portb := make([]byte, 4);
	conn.Read(portb);
	port = binary.LittleEndian.Uint32(portb);

	clientLocalAddr := net.TCPAddr{ip, int(port)};
	log.Stdout(clientLocalAddr.Port);
	log.Stdout("new client on: " + ip.String() + ":" + fmt.Sprintf("%d", port) + "\n");
}

func ListenConnections(listener *net.TCPListener, connections chan *net.TCPConn, clients chan TClient) {
	for {
		conn, err := listener.AcceptTCP();
		if err != nil {
			log.Stdout("error in Accept():", err)
		} else {
			go ProcessConn(conn, clients);
			connections <- conn;
		}
	}
}

func main() {
	addr, err := net.ResolveTCPAddr("127.0.0.1:4009");
	if err != nil {
		log.Exit("error:", err)
	}
	listener, err := net.ListenTCP("tcp", addr);
	if err != nil {
		log.Exit("error", err)
	}

	connections := make(chan *net.TCPConn);
	clients := make(chan TClient);

	go ListenConnections(listener, connections, clients);
	log.Stdout("Waiting for connections\n");
	for {
		select {
		case conn := <-connections:
			log.Stdout(conn)
		case client := <-clients:
			log.Stdout(client)
		}
	}
}
