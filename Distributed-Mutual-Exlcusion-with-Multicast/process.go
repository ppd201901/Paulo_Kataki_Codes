package main

import (
	"encoding/gob"
	"log"
	"net"
	"strconv"
	"sync"
)

type process struct {
	mutex              sync.Mutex
	address            string
	id                 int
	state              int
	timestamp          int
	requestTimestamp   int
	processesAddresses []string
	s                  replyCounter
	q                  messageQueue
	channels           []chan message
	receivedAllReplies chan bool
	channelIndex       map[string]int
}

func (p *process) getState() string {
	state := ""
	if p.state == WANTED {
		state = "WANTED"
	}

	if p.state == RELEASED {
		state = "RELEASED"
	}

	if p.state == HELD {
		state = "HELD"
	}
	return state
}

func (p *process) numberOfProcesses() int {
	return len(p.processesAddresses)
}

func (p *process) incrementReply() {
	p.updateTimestamp(p.timestamp)
	p.s.addReply()
	if p.s.size() == p.numberOfProcesses()-1 {
		p.receivedAllReplies <- true
	}
}

func (p *process) clearReplyCounter() {
	p.updateTimestamp(p.timestamp)
	p.s.clear()
}

func (p *process) enqueueMessage(msg message) {
	p.updateTimestamp(p.timestamp)
	p.q.push(msg)
}

func (p *process) getMessageQueueTopRequest() string {
	p.updateTimestamp(p.timestamp)
	return p.q.pop().Address
}

func (p *process) isMessageQueueEmpty() bool {
	p.updateTimestamp(p.timestamp)
	return p.q.empty() == true
}

func (p *process) updateTimestamp(timestamp int) {
	p.mutex.Lock()
	p.timestamp = max(p.timestamp, timestamp) + 1
	p.mutex.Unlock()
}

func (p *process) changeState(state int) {
	p.updateTimestamp(p.timestamp)
	p.mutex.Lock()
	p.state = state
	p.mutex.Unlock()
}

func (p *process) updateRequestTimestamp() {
	p.requestTimestamp = p.timestamp
}

func (p *process) getIndexFromAddress(address string) int {
	return p.channelIndex[address]
}

func (p *process) startProcess(address string, id int) {

	p.s = NewReplyCounter()
	p.q = NewMessageQueue()
	p.channelIndex = make(map[string]int)
	p.receivedAllReplies = make(chan bool)
	p.address = address
	p.id = id

	p.changeState(RELEASED)

	if err := p.startListenPort(); err != nil {
		log.Fatal("Error on startListenPort")
	}

	if err := p.getOtherProcessesAddresses(); err != nil {
		log.Fatal("Error on getOtherProcessesAddresses")
	}

	if err := p.openAllProcessesTCPConnections(); err != nil {
		log.Fatal("Error on openAllProcessesTCPConnections")
	}

	log.Println("Process ", p.id, " is ready")
}

func (p *process) sendPermissionToAllProcesses() {
	p.doMulticast(PERMISSION)
	p.waitAllProcessesReplies()
	p.clearReplyCounter()
}

func (p *process) openTCPConnection(address string, TCPWaiter chan bool) error {
	//open TCP connection on address
	connection, err := net.Dial("tcp", address)
	if err != nil {
		log.Println("Error in opening TCP port: ", address)
	}
	TCPWaiter <- true

	go func(address string, connection net.Conn) {
		var msg message

		//create encoder serializer
		encoder := gob.NewEncoder(connection)
		defer connection.Close()

		for {
			// channel waiting some message
			msg = <-p.channels[p.getIndexFromAddress(address)]
			log.Println(p.timestamp, " Process ", p.id, " is sending a ", msg.getType(), "with timestamp ", msg.Timestamp, " to ", address)

			//sending message
			if err := msg.encodeAndSendMessage(encoder); err != nil {
				log.Println(p.timestamp, " Error on Process ", p.id)
				log.Println(err)
			}
		}
	}(address, connection)

	return nil
}

func (p *process) openAllProcessesTCPConnections() error {
	//creating slice of channels
	p.channels = make([]chan message, p.numberOfProcesses())

	//channel to wait all the TCP
	TCPWaiter := make(chan bool, p.numberOfProcesses()-1)
	for i, address := range p.processesAddresses {

		p.channelIndex[address] = i

		//creating each channel
		p.channels[i] = make(chan message)

		if address != p.address {
			if err := p.openTCPConnection(address, TCPWaiter); err != nil {
				return err
			}
		}
	}

	//waiting all TCP connections to be ready
	for i := 0; i < p.numberOfProcesses()-1; i++ {
		<-TCPWaiter
	}
	return nil

}

func (p *process) sendMessage(typeMessage int, address string) {
	//creating new message
	msg := message{
		Timestamp:        p.timestamp,
		RequestTimestamp: p.requestTimestamp,
		TypeMessage:      typeMessage,
		Address:          p.address,
		Id:               p.id,
	}

	// sending message to address's channel
	p.channels[p.getIndexFromAddress(address)] <- msg
}

func (p process) doMulticast(typeMessage int) {
	p.updateTimestamp(p.timestamp)

	//send message to all processes
	for _, address := range p.processesAddresses {
		if address != p.address {
			go p.sendMessage(typeMessage, address)
		}
	}
}

//Method related to servers
func (p *process) getOtherProcessesAddresses() error {

	conn, err := net.Dial("tcp", REGISTER_SERVER_ADDRESS)
	if err != nil {
		return err
	}

	defer conn.Close()

	enc := gob.NewEncoder(conn)
	if err := enc.Encode(p.address); err != nil {
		return err
	}

	dec := gob.NewDecoder(conn)
	if err = dec.Decode(&p.processesAddresses); err != nil {
		return err
	}

	return nil
}

//Method related to server
func (p process) getRandomString() {
	conn, err := net.Dial("tcp", CRITICAL_REGION_SERVER_ADDRESS)

	defer conn.Close()

	if err != nil {
		log.Println("Dial on ", p.address)
		log.Fatal(err)
	}

	msg := p.address + "|" + strconv.Itoa(p.timestamp)

	enc := gob.NewEncoder(conn)
	err = enc.Encode(msg)

	dec := gob.NewDecoder(conn)
	err = dec.Decode(&msg)

	log.Println(p.timestamp, " PROCESS ", p.id, " received a message")

	if err != nil {
		log.Fatal("ERROR")
	}
	//time.Sleep(100 * time.Millisecond)
}
