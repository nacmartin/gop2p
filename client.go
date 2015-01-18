package main

import (
	//"io/ioutil";
	"io";
	"net";
	"fmt";
	"flag";
)

var server = flag.String("s", "127.0.0.1:4009", "server address")
var local = flag.String("c", "127.0.0.1:4005", "client address")
var dir = flag.String("d", "/data", "data dir")

func main() {
	flag.Parse();
	addr, err := net.ResolveTCPAddr("tcp",*local);
	if err != nil {
		fmt.Println("error resolving client:", err)
	}
	conn, err := net.Dial("tcp", *server);
	if err != nil {
		fmt.Println("error resolving server", err)
	}
	b := make([]byte, 1024);
	size, err := conn.Read(b);
	if err != nil {
		fmt.Println("error reading ", err)
	}
	fmt.Println("server says: " + string(b[0:size]));
	if !("hi" == string(b[0:size])) {
		fmt.Println("Error initializing", err)
	} else {
		fmt.Println("And that's fine")
	}
	conn.Write([]byte("I'm " + addr.String() + "\n"));
	//files, _ := ioutil.ReadDir(*dir);
	/*
	for _, v := range files {
        if !v.Fileinfo().IsDir(){
		    conn.Write([]byte("have " + v.Name + "\n"))
        }
	}
	*/
	conn.Write([]byte("list\n"));

	for {
		rcvStr := "";
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
		}
		fmt.Println(rcvStr);
	}
}
