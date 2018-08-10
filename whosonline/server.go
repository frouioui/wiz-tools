package main

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/jchavannes/go-pgp/pgp"
)

func encryptMessage(message string) string {
	pubEntity, err := pgp.GetEntity([]byte(pubkey), []byte{})
	if err != nil {
		println(fmt.Errorf("Error getting entity: %v", err))
	}
	fmt.Println("Created public key entity.")

	encrypted, err := pgp.Encrypt(pubEntity, []byte(TestMessage))
	if err != nil {
		println(err)
	}
	fmt.Println("Encrypted test message with public key entity.")

	println("msg chiffrer : ", string(encrypted))
}

func decryptMessage(message []byte) string {
	privEntity, err := pgp.GetEntity([]byte(pubkey), []byte(privkey))
	if err != nil {
		println(fmt.Errorf("Error getting entity: %v", err))
	}
	fmt.Println("Created private key entity.")

	decrypted, err := pgp.Decrypt(privEntity, encrypted)
	if err != nil {
		println(err)
	}
	fmt.Println("Decrypted message with private key entity.")

	decryptedMessage := string(decrypted)
	if decryptedMessage != TestMessage {
		println(errors.New("Decrypted message does not equal original."))
	}
	fmt.Println("Decrypted message equals original message.")
	fmt.Println("Entcrypt test: END\n")
	println(decryptedMessage)
}

func checkActiveComputer(i int, channel chan bool) {
	ip := fmt.Sprintf("%s.%d:%d", "192.168.2", i, 8754)
	conn, err := net.DialTimeout("tcp", ip, 2*time.Second)
	if err != nil {
		channel <- false
	} else {
		_, err := conn.Write([]byte("who"))
		checkErr(err)
		buffer := make([]byte, 16)
		_, err = conn.Read(buffer)
		checkErr(err)
		fmt.Printf("%s #%s\n", ip, buffer)
		channel <- true
		conn.Close()
	}
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	channels := make(chan bool, 255)

	for i := 0; i <= 254; i++ {
		go checkActiveComputer(i, channels)
	}
	for i := 0; i <= 254; i++ {
		<-channels
	}
	println("done")
}
