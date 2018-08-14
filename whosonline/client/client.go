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

//Location of own public key
const locOwnPubKey = "/home/florian/Documents/wiz-tools/gpg/pubkey_client.asc"

//Location of own private key
const locOwnPrivKey = "/home/florian/Documents/wiz-tools/gpg/privkey_client.asc"

//Location of the other person public key
const locOthPubKey = "/home/florian/Documents/wiz-tools/gpg/clefpub.asc"

//whoIsIt get the result of the command whoisit
func whoIsIt() []byte {
	out, err := exec.Command("whoisit").Output()
	checkErr(err)
	return (out)
}

func getTheCmd(conn net.Conn, keys *encode.Keys) encode.Requete {
	buffer := bufio.NewReader(conn)
	lines, err := buffer.ReadBytes(0x00)
	if err != io.EOF && err != nil {
		panic(err)
	}
	if len(lines) <= 800 {
		panic("line trop court")
	}
	signature := lines[:800]
	msg := lines[800 : len(lines)-1]
	if keys.Verify(msg, signature) == false {
		panic("[ERROR] : Mauvaise signature..")
	}
	cmd := keys.Uncrypt(string(msg))
	var dat encode.Requete
	if err := json.Unmarshal([]byte(cmd), &dat); err != nil {
		panic(err)
	}
	return dat
}

func handleConnection(conn net.Conn, keys *encode.Keys) {
	defer conn.Close()
	defer println("[-] Fin de la connexion")
	data := getTheCmd(conn, keys)
	print(data.Cmd)
	switch data.Cmd {
	case "who":
		_, err := conn.Write(whoIsIt())
		checkErr(err)
		fmt.Println("[who] : -> Message Send")
	default:
		fmt.Println("[ERREUR] Mauvaise Commande")
	}
	/*requestb := make([]byte, 3)
	_, err := conn.Read(requestb)
	checkErr(err)
	switch request := fmt.Sprintf("%s", requestb); request {
	case "who":
		_, err = conn.Write(whoIsIt())
		checkErr(err)
		fmt.Println("[who] : -> Message Send")
	case "msg":
		title, body, img := make([]byte, 32), make([]byte, 256), make([]byte, 128)
		_, err := conn.Read(title)
		checkErr(err)
		fmt.Printf("Title : %s\n", title)
		_, err = conn.Read(body)
		checkErr(err)
		fmt.Printf("Body : %s\n", body)
		_, err = conn.Read(img)
		checkErr(err)
		fmt.Printf("Image : %s\n", img)
		displayNotification(string(title), string(body), string(img))
	default:
		fmt.Printf("[%s] : Unknow command..", request)
	}*/
}

func checkErr(e error) {
	if e != nil {
		panic(e)
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
			fmt.Println("[+] New connexion")
			go handleConnection(conn, &keys)
		}
	}
}
