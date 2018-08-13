package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/redcurrents/wiz-tools/whosonline/encode"
)

//Location of own public key
const locOwnPubKey = "/home/florian/Documents/wiz-tools/gpg/clefpub.asc"

//Location of own private key
const locOwnPrivKey = "/home/florian/Documents/wiz-tools/gpg/clefpub.asc"

//Location of the other person public key
const locOthPubKey = "/home/florian/Documents/wiz-tools/gpg/clefpub.asc"

func handleConnection(conn net.Conn) {
	defer conn.Close()
	privkey, err := ioutil.ReadFile(locOwnPrivKey)
	if err != nil {
		log.Fatal(err)
	}
	encode.Encrypt("", "")
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	port := ":8754"
	pubkey, err := ioutil.ReadFile(locOwnPubKey)
	if err != nil {
		log.Fatal(err)
	}
	othpubkey, err := ioutil.ReadFile(locOthPubKey)
	if err != nil {
		log.Fatal(err)
	}
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
