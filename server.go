package main

import (
	"net";
	"log";
    "os";
	"strings";
    "time";
    "regexp";
    "./util";
)

type TClient struct {
	local		string;
	forserver	string;
	msg			string;
}

type SFile struct {
    local   string;
    name    string;
}

func myReader(conn *net.TCPConn, client chan TClient){
    var rcvStr string;
    var localAddr string;
    for{
        rcvStr = "";
        read := true;
        for read {
            rcvd := make([]byte, 1);
            size, err := conn.Read(rcvd);
            switch err {
            case os.EOF:
               //log.Stdout("Warning: End of data reached: ", err);
               read = false;
            case nil:
                if(util.Streq(string(rcvd[0:1]),"\n")){
                    read = false;
                }else{
                    rcvStr = rcvStr + string(rcvd[0:size]);
                }
            default:
               log.Stdout("Error: Reading data: ", err);
               read = false;
            }
            if(util.Streq(string(rcvd[0:1]),"\n")){
               read = false;
            }
        }
        if regexp.MustCompile("^I'm").MatchString(rcvStr) {
            ladr, _ := net.ResolveTCPAddr(strings.Split(rcvStr," ",2)[1]);
            localAddr = ladr.String();
            log.Stdout(localAddr);
            client <-TClient{localAddr, conn.RemoteAddr().String(), "new"};
        }else{
            if len(rcvStr) > 0 {
                client <-TClient{localAddr, conn.RemoteAddr().String(), rcvStr};
                log.Stdout("Data sent by client: " + rcvStr);
            }
            time.Sleep(5*1000*1000);
        }
    }
}

func ProcessConn(conn *net.TCPConn, client chan TClient) {
	log.Stdout("connected\n");

	//get client's ip and port
	conn.Write(strings.Bytes("hi"));
    go myReader(conn, client);
}

func ListenConnections(listener *net.TCPListener, connections chan *net.TCPConn, clients chan TClient) {
	for {
		conn, err := listener.AcceptTCP();
		if err != nil {
			log.Stdout("error in Accept():", err)
		} else {
			conn.SetKeepAlive(true);
            conn.SetReadTimeout(5*1000*1000*1000);
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

	//1 channel for incoming connections, another for client communication
	connections := make(chan *net.TCPConn);
	clients := make(chan TClient);
    cMap := make(map[string] *net.TCPConn);
    fMap := make(map[string] string);

	go ListenConnections(listener, connections, clients);
	log.Stdout("Waiting for connections\n");
	for {
		select {
		case conn := <-connections:
            cMap[conn.RemoteAddr().String()] = conn;
		case client := <-clients:
            if regexp.MustCompile("^have ").MatchString(client.msg){
                fMap[string(client.msg[5:len(client.msg)])] = client.local;
            }
            if regexp.MustCompile("^list").MatchString(client.msg){
                for key, value := range fMap {
                    cMap[client.forserver].Write(strings.Bytes(key+"->"+value));
                }
                cMap[client.forserver].Write(strings.Bytes("\n"));
            }
		}
	}
}


