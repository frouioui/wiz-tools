package main

import (
	"bytes"
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

func getTheCmd(conn net.Conn, keys *encode.Keys) map[string]string {
	signature := make([]byte, 800)
	n, err := conn.Read(signature)
	if n != 800 {
		panic("bad signature")
	}
	checkErr(err)
	var buf bytes.Buffer
	io.Copy(&buf, conn)
	if keys.Verify(buf.Bytes(), signature) == false {
		panic("[ERROR] : Mauvaise signature")
	}
	cmd := keys.Uncrypt(buf.String())
	var dat map[string]string
	if err := json.Unmarshal([]byte(cmd), &dat); err != nil {
		panic(err)
	}
	//println(dat["cmd"])
	/*for key, value := range dat {
		fmt.Printf("Key : %s | Value : %s\n", key, value)
	}*/
	return dat
}

func handleConnection(conn net.Conn, keys *encode.Keys) {
	defer conn.Close()
	data := getTheCmd(conn, keys)
	switch data["cmd"] {
	case "who":
		_, err := conn.Write(whoIsIt())
		checkErr(err)
		fmt.Println("[who] : -> Message Send")
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
