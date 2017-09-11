package network

import (
	"net"
	"fmt"
	"log"
	"time"
	"kademliaid"
	"github.com/vmihailenco/msgpack"
)

type Network struct {
	timeout time.Duration
}

func New(timeoutms int64) Network {
	return Network{timeout: timeoutms * time.Millisecond}
}

const (
	PING = 0
	FIND_NODE = 1
	FIND_VALUE = 2
	STORE = 3
)

type RPC struct {
	rpc_type int
	sender_id kademliaid.KademliaID
}

type RPCPing struct {
	RPC
}

type RPCFindNode struct {
	RPC
	find_id kademliaid.KademliaID
}

type RPCFindValue struct {
	RPC
	find_id kademliaid.KademliaID
}

type RPCStore struct {
	RPC
	store_value []byte
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
			continue
		}
		go resolveRPC(b, raddr)
	}
	// unreachable
}

func (network *Network) send(contact *Contact, msg []byte) (net.UDPConn, error) {
	// TODO: contact should probably store the resolved address already
	// right now it's a string (who wrote this sample code?!)
	raddr := net.ResolveUDPAddr("udp", contact.Address)
	// TODO: would it be safe to use the port we are listening on?
	// net.Conn docs seem ok with it:
	// Multiple goroutines may invoke methods on a Conn simultaneously.
	conn, err = net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Printf("Error dialing %v: %v\n", raddr, err)
		// TODO: do we need to close conn here too? might depend on the error
		return nil, err
	}
	conn.WriteToUDP(msg, raddr)
	return conn, nil
}

func (network *Network) receive(conn net.UDPConn) ([]byte, error) {
	// TODO: make this buffer size configurable somewhere
	p := make([]byte, 2048)
	conn.SetReadDeadline(network.timeout)
	_, err = bufio.NewReader(conn).Read(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (network *Network) rpc(contact *Contact, msg *RPC) (map[string]interface{}, error) {
	b, err := msgpack.Marshal(&msg)
	if err != nil {
		log.Printf("Error marshalling FindNode RPC: %v\n", err)
		return nil, err
	}
	conn, err := send(contact, b)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	rb, err := receive(conn)
	if err != nil {
		return nil, err
	}
	var response map[string]interface{}
	err = msgpack.Unmarshal(rb, &response)
	if err != nil {
		return nil, err
	}
	// TODO: update routing table
	return response, nil
}

func (network *Network) Ping(contact *Contact) error {
	msg := &RPCPing{rpc_type: PING, sender_id: network.id}
	response, err := rpc(contact, msg)
	if err != nil {
		return err
	}
	// routing table is updated as a side effect of receiving the response
	return nil
}

func (network *Network) FindNode(contact *Contact, findID *kademliaid.KademliaID) ([]Contact, error) {
	msg := &RPCFindNode{rpc_type: FIND_NODE, sender_id: network.id, find_id: *findID}
	response, err := rpc(contact, msg)
	if err != nil {
		return err
	}
	// TODO: do something with response and return
}

func (network *Network) FindValue(contact *Contact, findID *kademliaid.KademliaID) {
	msg := &RPCFindValue{rpc_type: FIND_VALUE, sender_id: network.id, find_id: *findID}
	response, err := rpc(contact, msg)
	if err != nil {
		return err
	}
	// TODO: do something with response and return
}

func (network *Network) Store(data []byte) error {
	msg := &RPCStore{rpc_type: STORE, sender_id: network.id, store_value: data}
	// not using rpc() since this rpc doesn't need a response
	b, err := msgpack.Marshal(&msg)
	if err != nil {
		log.Printf("Error marshalling FindNode RPC: %v\n", err)
		return err
	}
	conn, err := send(contact, b)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func resolveRPC(message []byte, raddr *net.UDPAddr) {
	// TODO: extract the RPC enum from the message
	args := map[string]interface{}
	err := msgpack.Unmarshal(b, &args)
	switch args["rpc_type"] {
	case PING:
		pingResponse(raddr)
	case FIND_NODE:
		findNodeResponse(&args, raddr)
	case FIND_VALUE:
		findValueResponse(&args, raddr)
	case STORE:
		storeResponse(&args)
	default:
		log.Printf("Unknown RPC: %v\n", args["rpc_type"])
		// garbage message, don't update routing table
		return
	}
	// TODO: update our routing table
}

func storeResponse(args *map[string]interface{}) {
	// no confirmation is sent
	// TODO: actually store the value from args
}

func (network *Network) pingResponse(raddr *net.UDPAddr) {
}

func (network *Network) storeResponse() {
	// TODO: try to find it in the kv store
	// return
	// if we can't find it, just treat it as a FindNode
	network.FindNode(contact, findID)
}
