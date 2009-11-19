package main

import (
	"net";
	"log";
	"flag";
	"strings";
    "io";
    "os";
    "./util";
)

var server = flag.String("s", "127.0.0.1:4009", "server address")
var local = flag.String("c", "127.0.0.1:4005", "client address")
var dir = flag.String("d", "/data", "data dir")

func main() {
	flag.Parse();
	addr, err := net.ResolveTCPAddr(*local);
	if err != nil {
		log.Exit("error resolving client:", err)
	}
	conn, err := net.Dial("tcp", "", *server);
	if err != nil {
		log.Exit("error resolving server", err)
	}
	b := make([]byte, 1024);
	size, err := conn.Read(b);
	if err != nil {
		log.Exit("error reading ", err)
	}
	log.Stdout("server says: " + string(b[0:size]));
	if !util.Streq("hi", string(b[0:size])) {
		log.Exit("Error initializing", err)
	} else {
		log.Stdout("And that's fine")
	}
	conn.Write(strings.Bytes("I'm "+addr.String()+"\n"));
    files, err := io.ReadDir(*dir);
    if err != nil {
        log.Exit("error reading dir", err);
    }
    for _, v := range files{
        conn.Write(strings.Bytes("have "+v.Name+"\n"));
    }
	conn.Write(strings.Bytes("list\n"));

    for{
        rcvStr := "";
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
        }
        log.Stdout(rcvStr);
    }
}
