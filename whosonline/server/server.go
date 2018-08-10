package main

import (
	"fmt"
	"net"
	"time"
)

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
