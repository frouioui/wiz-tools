package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"

	notify "github.com/mqu/go-notify"
	"github.com/redcurrents/wiz-tools/whosonline/encode"
)

const (
	//Location of own public key
	locOwnPubKey = "/home/florian/Documents/wiz-tools/gpg/pubkey_client.asc"
	//Location of own private key
	locOwnPrivKey = "/home/florian/Documents/wiz-tools/gpg/privkey_client.asc"
	//Location of the other person public key
	locOthPubKey = "/home/florian/Documents/wiz-tools/gpg/clefpub.asc"
)

//whoIsIt get the result of the command whoisit
func whoIsIt() []byte {
	out, err := exec.Command("whoisit").Output()
	if err != nil {
		return nil
	}
	return (out)
}

func getTheCmd(conn net.Conn, keys *encode.Keys) (encode.Requete, int) {
	buffer := bufio.NewReader(conn)
	lines, err := buffer.ReadBytes(0x00)
	if err != io.EOF && err != nil {
		log.Println(err)
		return encode.Requete{}, 1
	}
	if len(lines) <= 800 {
		log.Println("Minimum 800 bytes..")
		return encode.Requete{}, 1
	}
	signature := lines[:800]
	msg := lines[800 : len(lines)-1]
	if keys.Verify(msg, signature) == false {
		log.Printf("Mauvaise signature de %s\n", conn.RemoteAddr())
		return encode.Requete{}, 1
	}
	cmd := keys.Uncrypt(msg)
	var data encode.Requete
	if err := json.Unmarshal([]byte(cmd), &data); err != nil {
		log.Println("Impossible de récupérer la requete json..")
		return encode.Requete{}, 1
	}
	return data, 0
}

func handleConnection(conn net.Conn, keys *encode.Keys) {
	defer conn.Close()
	defer fmt.Printf("[-] Fin de la connexion\n")
	data, code := getTheCmd(conn, keys)
	if code == 1 {
		return
	}
	switch data.Cmd {
	case "who":
		_, err := conn.Write(whoIsIt())
		if err != nil {
			log.Printf("Impossible d'écrire le nom d'utilisateur..\n")
			return
		}
		fmt.Printf("[%s] : -> Message Send\n", data.Cmd)
	case "notif":
		//J'aime vraiment pas ça, à changer d'urgence pour utiliser la structure encode.Options
		options := data.Opt.(map[string]interface{})
		displayNotification(options["Title"].(string), options["Text"].(string), options["Image"].(string))
	default:
		fmt.Println("[ERREUR] Mauvaise Commande")
	}
}

func displayNotification(title, msg, img string) {
	notify.Init("ClientWiz")
	clientWiz := notify.NotificationNew(title,
		msg,
		img)

	if clientWiz == nil {
		fmt.Fprintf(os.Stderr, "Unable to create a new notification\n")
		return
	}

	notify.NotificationSetTimeout(clientWiz, 1000*10)

	if e := notify.NotificationShow(clientWiz); e != nil {
		fmt.Fprintf(os.Stderr, "%s\n", e.Message())
		return
	}
	notify.NotificationClose(clientWiz)
	notify.UnInit()
}

func main() {
	port := ":8754"
	keys := encode.Init(locOwnPubKey, locOwnPrivKey, locOthPubKey)
	if string(keys.Othpubkey) == "" {
		log.Fatal("impossible de get les clefs..")
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
			fmt.Printf("[+] New connexion (%s)\n", conn.RemoteAddr())
			go handleConnection(conn, &keys)
		}
	}
}
