package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/redcurrents/wiz-tools/whosonline/encode"
)

const (
	//Port to connect
	port = 8754
	//Location of own public key
	locOwnPubKey = "/home/florian/Documents/wiz-tools/gpg/clefpub.asc"
	//Location of own private key
	locOwnPrivKey = "/home/florian/.gnupg/clefpriv.asc"
	//Location of the other person public key
	locOthPubKey = "/home/florian/Documents/wiz-tools/gpg/pubkey_client.asc"
)

func etablishNewCo(subnet string, i, port int) (net.Conn, error) {
	ip := fmt.Sprintf("%s.%d:%d", subnet, i, port)
	conn, err := net.DialTimeout("tcp", ip, 1*time.Second)
	return conn, err
}

func sendCmd(conn net.Conn, keys *encode.Keys, req encode.Requete) bool {
	reqjson, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return false
	}
	newWriter := bufio.NewWriter(conn)
	message := keys.Encrypt(reqjson)
	if message == nil {
		log.Println(err)
		return false
	}
	signature := keys.Sign(message)
	if signature == nil {
		log.Println(err)
		return false
	}
	n, err := newWriter.Write(signature)
	if err != nil || n != len(signature) {
		log.Println(err)
		return false
	}
	n, err = newWriter.Write(message)
	if err != nil || n != len(message) {
		log.Println(err)
		return false
	}
	err = newWriter.WriteByte(0x00)
	if err != nil {
		log.Println(err)
		return false
	}
	newWriter.Flush()
	return true
}

//checkActiveComputer will check for an ip if the computer is alive or not and who is it
func checkActiveComputer(i int, channel chan bool, keys *encode.Keys) {
	conn, err := etablishNewCo("192.168.2", i, port)
	if err != nil {
		channel <- false
	} else {
		defer conn.Close()
		req := encode.Requete{Cmd: "who", Opt: encode.Options{}}
		if sendCmd(conn, keys, req) == true {
			newReader := bufio.NewReader(conn)
			line, _, _ := newReader.ReadLine()
			fmt.Printf("%s #%s\n", conn.RemoteAddr(), line)
			channel <- true
		} else {
			channel <- false
		}
	}
}

func main() {
	channels := make(chan bool, 255)
	keys := encode.Init(locOwnPubKey, locOwnPrivKey, locOthPubKey)

	if string(keys.Othpubkey) == "" {
		log.Fatal("Impossible de get les clefs..")
	}
	for i := 0; i <= 254; i++ {
		go checkActiveComputer(i, channels, &keys)
	}
	for i := 0; i <= 254; i++ {
		<-channels
	}
	fmt.Println("done")
}
