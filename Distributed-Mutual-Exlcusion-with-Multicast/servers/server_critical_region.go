package main

import (
	"encoding/gob"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func make_random_string(length int, charset string) string {

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generate_random_string(length int) string {
	return make_random_string(length, charset)
}

func handle_connection(c net.Conn) {

	// log.Println("received a request")

	var buffer string
	dec := gob.NewDecoder(c)
	err := dec.Decode(&buffer)

	if err != nil {
		log.Fatal(err)
		log.Fatal("Fail to Decode")
	}

	info := strings.Split(buffer, "|")
	log.Println("\n\n")
	log.Println(info[0], " entered in critical region with timestamp ", info[1])

	time.Sleep(time.Duration(seededRand.Intn(500)) * time.Millisecond)
	enc := gob.NewEncoder(c)
	err = enc.Encode(generate_random_string(1e5))

	if err != nil {
		log.Fatal("Fail to Encode")
	}

	log.Println(info[0], " left in critical region with timestamp ", info[1])
}

func run_server() {

	l, err := net.Listen("tcp", "localhost:8081")

	if err != nil {
		log.Fatal("Fail to connect")
	}

	log.Println("Server of Critical Region is running")

	defer l.Close()

	for {

		c, err := l.Accept()
		if err != nil {
			log.Fatal("Fail to connect")
		}

		handle_connection(c)
	}
}

func main() {

	waiter := make(chan int)
	go run_server()
	<-waiter
}
