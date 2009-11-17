package main

import (
        "net";
        "log";
        "strings";
        "os";
    )

func ProcessConn(conn *net.TCPConn){
    log.Stdout("connected\n");
    conn.Write(strings.Bytes("hi\n"));
    b := make([]byte, 1024);
    conn.Read(b[0:15]);
    os.Stdout.WriteString(net.IP(b[0:15]).String());
    conn.Close();
}

func main(){
    addr, err := net.ResolveTCPAddr("127.0.0.1:4009");
    if err != nil {
        log.Exit("error:", err);
    }
    listener, err := net.ListenTCP("tcp", addr);
    if err != nil {
        log.Exit("error", err);
    }
    //c := make(chan int);
    for{
        log.Stdout("Waiting for connections\n");
        conn, err := listener.AcceptTCP();
        if err != nil {
            log.Stdout("error in Accept():", err);
        }else{
            go func () {
                ProcessConn(conn);
            }();
        }
    }
}
