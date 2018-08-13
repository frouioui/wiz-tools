package main whoisit

import (
	"fmt"
	"log"
	"net"

	"github.com/redcurrents/wiz-tools/whosonline/encode"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	encode.encrypt()
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	port := ":8754"
	ln, err := net.Listen("tcp4", port)
	if err != nil {
		log.Fatal("Impossible d'écouter sur le port" + port)
	}
	fmt.Printf("[*] Le serveur écoute sur le port %s\n", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			println(err)
		} else {
			fmt.Println("[+] New connexion")
			go handleConnection(conn)
		}
	}
}
