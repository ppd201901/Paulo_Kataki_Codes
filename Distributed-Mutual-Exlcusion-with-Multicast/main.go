package main

import (
	"log"
	"strconv"
)

func main() {

	const NUMBER_PROCESSES = 3
	const INITIAL_PORT = 8090

	var p [NUMBER_PROCESSES]process

	waiter := make(chan bool)

	for i := 0; i < NUMBER_PROCESSES; i++ {
		port := INITIAL_PORT + i
		go func(i int) {
			s := strconv.Itoa(port)
			log.Println(port, "localhost:"+s)
			p[i].startProcess("localhost:"+s, port)
			waiter <- true
		}(i)
	}

	for i := 0; i < NUMBER_PROCESSES; i++ {
		<-waiter
	}

	for i := 0; i < NUMBER_PROCESSES; i++ {
		go p[i].runProcess()
	}

	<-waiter
}
