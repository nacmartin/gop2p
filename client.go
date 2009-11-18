package main

import (
	"net";
	"log";
	"flag";
	"encoding/binary";
	//        "fmt";
	//        "strings";
)

var server = flag.String("s", "127.0.0.1:4009", "server address")
var local = flag.String("c", "127.0.0.1:4005", "client address")

func streq(s1 string, s2 string) bool {
	if len(s1) != len(s2) {
		log.Stdout("nomi");
		log.Stdout(len(s1));
		log.Stdout(len(s2));
		return false;
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true;
}

func main() {
	flag.Parse();
	addr, err := net.ResolveTCPAddr(*local);
	if err != nil {
		log.Exit("error resolving client:", err)
	}
	log.Stdout(addr.Port);
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
	if !streq("hi", string(b[0:size])) {
		log.Exit("Error initializing")
	}
	conn.Write(addr.IP);
	binary.LittleEndian.PutUint32(b, uint32(addr.Port));
	conn.Write(b[0:3]);
}
