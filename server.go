package main

import (
	"net";
	"log";
	"strings";
	"fmt";
	"encoding/binary";
)

func ProcessConn(conn *net.TCPConn) {
	defer conn.Close();
	log.Stdout("connected\n");
	conn.Write(strings.Bytes("hi"));
	ip16 := make([]byte, 16);
	conn.Read(ip16);
	ip := net.IP(ip16);
	var port uint32;
	portb := make([]byte, 4);
	conn.Read(portb);
	port = binary.LittleEndian.Uint32(portb);
	log.Stdout("new client on: " + ip.String() + ":" + fmt.Sprintf("%d", port) + "\n");
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
	//c := make(chan int);
	for {
		log.Stdout("Waiting for connections\n");
		conn, err := listener.AcceptTCP();
		if err != nil {
			log.Stdout("error in Accept():", err)
		} else {
			go func() { ProcessConn(conn) }()
		}
	}
}
