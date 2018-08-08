package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
)

func whoIsIt() string {
	out, err := exec.Command("whoisit").Output()
	checkErr(err)
	user := fmt.Sprintf("%s", out)
	return (user)
}

func handleConnection(conn net.Conn) {
	message := whoIsIt()
	_, err := conn.Write([]byte(message))
	checkErr(err)
	fmt.Println("-> Message Send")
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
