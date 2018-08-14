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

//Location of own public key
const locOwnPubKey = "/home/florian/Documents/wiz-tools/gpg/clefpub.asc"

//Location of own private key
const locOwnPrivKey = "/home/florian/.gnupg/clefpriv.asc"

//Location of the other person public key
const locOthPubKey = "/home/florian/Documents/wiz-tools/gpg/pubkey_client.asc"

//checkActiveComputer will check for an ip if the computer is alive or not and who is it
func checkActiveComputer(i int, channel chan bool, keys *encode.Keys) {
	ip := fmt.Sprintf("%s.%d:%d", "192.168.2", i, 8754)
	conn, err := net.DialTimeout("tcp", ip, 2*time.Second)
	if err != nil {
		channel <- false
	} else {
		defer conn.Close()
		newWriter := bufio.NewWriter(conn)
		test := encode.Requete{Cmd: "who", Opt: encode.Options{}}
		b, err := json.Marshal(test)
		checkErr(err)
		message := keys.Encrypt(string(b))
		signature := keys.Sign(message)
		println(len(signature), len(message))
		newWriter.Write(signature)
		newWriter.Write(message)
		newWriter.WriteByte(0x00)
		newWriter.Flush()
		newReader := bufio.NewReader(conn)
		line, _, _ := newReader.ReadLine()
		fmt.Printf("%s #%s\n", ip, line)
		//var avi map[string]string{} = {"avion": "avion", "toz": "poubelle"}
		/*test := "{\"cmd\": \"who\"}"
		message := keys.Encrypt(test)
		signature := keys.Sign(message)
		_, err := conn.Write(signature)
		checkErr(err)
		_, err = conn.Write(message)
		checkErr(err)
		time.Sleep(time.Second * 10)
		buf, err := ioutil.ReadAll(conn)
		checkErr(err)
		println(buf)*/
		//println(buf.String())
		//fmt.Printf("%s #%s\n", ip, buffer)
		channel <- true
	}
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	channels := make(chan bool, 255)
	keys := encode.Init(locOwnPubKey, locOwnPrivKey, locOthPubKey)
	if string(keys.Othpubkey) == "" {
		log.Fatal("impossible de get les clefs..")
	}
	for i := 0; i <= 254; i++ {
		go checkActiveComputer(i, channels, &keys)
	}
	for i := 0; i <= 254; i++ {
		<-channels
	}
	println("done")
}
