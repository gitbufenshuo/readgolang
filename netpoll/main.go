package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	Server("0.0.0.0", "9999")
	time.Sleep(time.Hour)
}

func Server(ip, port string) {
	fmt.Println("Server")
	hostPort := ip + ":" + port
	l, e := net.Listen("tcp4", hostPort)
	if e != nil {
		log.Println(e.Error())
		os.Exit(0)
	}
	fmt.Println("ready for accept")
	go func() {
		time.Sleep(time.Second * 3)
		fmt.Println("ready for accept    333")
		l.Accept()
	}()
	{
		coo, err := l.Accept()
		if err != nil {
			log.Println(err.Error())
			os.Exit(0)
		}
		coo.Read(nil)
	}
}
