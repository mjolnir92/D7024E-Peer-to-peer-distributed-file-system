package d7024e

import (
	"net"
	"fmt"
	"log"
)

type Network struct {
}

func Listen(ip string, port int) {
	b := make([]byte, 2048)
	addrStr := fmt.Sprintf("%s:%d", ip, port)
	addr := net.ResolveUDPAddr("udp", addrStr)
	srv, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Error listening on %v: %v\n", addr, err)
	}
	for {
		_, raddr, err := srv.ReadFromUDP(b)
		if err != nil {
			log.Printf("Error reading UDP: %v", err)
		}
		// TODO: unmarshal
		// TODO: in new thread
		// TODO: marshal response
		// TODO: send it
		go resolveRPC(b, raddr)
}

func (network *Network) SendPingMessage(contact *Contact) {
}

func (network *Network) SendFindContactMessage(contact *Contact) {
}

func (network *Network) SendFindDataMessage(hash string) {
}

func (network *Network) SendStoreMessage(data []byte) {
}

// TODO: Send{Ping, FindContact, FindData, Store}Response
