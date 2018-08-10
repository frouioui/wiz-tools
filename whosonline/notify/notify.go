package main

import (
	"fmt"
	"net"
	"time"
)

func checkActiveComputer() {
	ip := fmt.Sprintf("%s.%d:%d", "192.168.2", 195, 8754)
	conn, _ := net.DialTimeout("tcp", ip, 2*time.Second)
	_, err := conn.Write([]byte("msg"))
	checkErr(err)
	_, err = conn.Write([]byte("Salut a tousse"))
	checkErr(err)
	_, err = conn.Write([]byte("Diablox9 c le plus fort"))
	checkErr(err)
	_, err = conn.Write([]byte("/home/florian/Documents/wiz-tools/whosonline/res/Wiztopic-Logos.png"))
	checkErr(err)
	conn.Close()
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	checkActiveComputer()
}
