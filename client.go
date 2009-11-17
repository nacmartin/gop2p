package main

import (
        "net";
        "log";
        "os";
        "flag";
//        "fmt";
//        "strings";
        )

var server = flag.String("s", "127.0.0.1:4009", "server address")
var local  = flag.String("c", "127.0.0.1:4005", "client address")

func main(){
    flag.Parse();
    addr, err := net.ResolveTCPAddr(*local);
    if err != nil {
        log.Exit("error resolving client:", err);
    }

    conn, err := net.Dial("tcp","",*server);
    if err != nil {
        log.Exit("error resolving server", err);
    }
    b := make([]byte, 1024);
    conn.Read(b);
    os.Stdout.Write(b);
    conn.Write(addr.IP);
}
