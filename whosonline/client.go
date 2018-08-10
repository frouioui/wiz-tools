package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	notify "github.com/mqu/go-notify"
)

//whoIsIt get the result of the command whoisit
func whoIsIt() []byte {
	out, err := exec.Command("whoisit").Output()
	checkErr(err)
	return (out)
}

//
func handleConnection(conn net.Conn) {
	defer conn.Close()
	requestb := make([]byte, 3)
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
	}
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
