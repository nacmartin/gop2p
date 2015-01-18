package main

import (
	"io";
	"net";
	"fmt";
	"strings";
	"time";
	"regexp";
)

type TClient struct {
	local		string;
	forserver	string;
	msg			string;
}

type SFile struct {
	local	string;
	name	string;
}

func myReader(conn *net.TCPConn, client chan TClient) {
	var rcvStr string;
	var localAddr string;
	for {
		rcvStr = "";
		read := true;
		for read {
			rcvd := make([]byte, 1);
			size, err := conn.Read(rcvd);
			switch err {
			case io.EOF:
				//log.Stdout("Warning: End of data reached: ", err);
				read = false
			case nil:
				if (string(rcvd[0:1]) == "\n") {
					read = false
				} else {
					rcvStr = rcvStr + string(rcvd[0:size])
				}
			default:
				fmt.Println("Error: Reading data: ", err);
				read = false;
			}
			if (string(rcvd[0:1])== "\n") {
				read = false
			}
		}
		if regexp.MustCompile("^I'm").MatchString(rcvStr) {
			ladr, _ := net.ResolveTCPAddr("tcp", strings.Split(rcvStr, " ")[1]);
			localAddr = ladr.String();
			fmt.Println(localAddr);
			client <- TClient{localAddr, conn.RemoteAddr().String(), "new"};
		} else {
			if len(rcvStr) > 0 {
				client <- TClient{localAddr, conn.RemoteAddr().String(), rcvStr};
				fmt.Println("Data sent by client: " + rcvStr);
			}
			time.Sleep(5 * 1000 * 1000);
		}
	}
}

func ProcessConn(conn *net.TCPConn, client chan TClient) {
	fmt.Println("connected\n");

	//get client's ip and port
	conn.Write([]byte("hi"));
	go myReader(conn, client);
}

func ListenConnections(listener *net.TCPListener, connections chan *net.TCPConn, clients chan TClient) {
	for {
		conn, err := listener.AcceptTCP();
		if err != nil {
			fmt.Println("error in Accept():", err)
		} else {
			conn.SetKeepAlive(true);
			conn.SetReadDeadline(time.Now().Add(10 * time.Second));
			go ProcessConn(conn, clients);
			connections <- conn;
		}
	}
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp","127.0.0.1:4009");
	
	if err != nil {
		fmt.Println("error:", err)
	}

	listener, err := net.ListenTCP("tcp", addr);
	if err != nil {
		fmt.Println("error", err)
	}

	//1 channel for incoming connections, another for client communication
	connections := make(chan *net.TCPConn);
	clients := make(chan TClient);
	cMap := make(map[string]*net.TCPConn);
	fMap := make(map[string]string);

	go ListenConnections(listener, connections, clients);
	fmt.Println("Waiting for connections\n");
	for {
		select {
		case conn := <-connections:
			cMap[conn.RemoteAddr().String()] = conn
		case client := <-clients:
			if regexp.MustCompile("^have ").MatchString(client.msg) {
				fMap[string(client.msg[5:len(client.msg)])] = client.local
			}
			if regexp.MustCompile("^list").MatchString(client.msg) {
				for key, value := range fMap {
					cMap[client.forserver].Write([]byte(key + "->" + value))
				}
				cMap[client.forserver].Write([]byte("\n"));
			}
		}
	}
}
