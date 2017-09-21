package network

import (
	"net"
	"fmt"
	"log"
	"time"
	"bufio"
	"github.com/mjolnir92/kdfs/kademliaid"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/routingtable"
	"github.com/mjolnir92/kdfs/kvstore"
	"github.com/vmihailenco/msgpack"
)

type T struct {
	timeout time.Duration
	id *kademliaid.T
	routingtable *routingtable.T
	kvstore *kvstore.T
	conn *net.UDPConn
}

func New(timeoutms int64, id *kademliaid.T) T {
	// TODO: do more of the setup here
	return T{timeout: time.Duration(timeoutms) * time.Millisecond, id: id}
}

const (
	PING = 0
	PING_RESPONSE = 1
	FIND_NODE = 2
	FIND_NODE_RESPONSE = 3
	FIND_VALUE = 4
	FIND_VALUE_RESPONSE = 5
	STORE = 6
)

type RPCPing struct {
	RPCType int
	SenderID kademliaid.T
}

type RPCPingResponse struct {
	RPCType int
	SenderID kademliaid.T
}

type RPCFindNode struct {
	RPCType int
	SenderID kademliaid.T
	FindID kademliaid.T
}

type RPCFindNodeResponse struct {
	RPCType int
	SenderID kademliaid.T
	Contacts []contact.T
}

type RPCFindValue struct {
	RPCType int
	SenderID kademliaid.T
	FindID kademliaid.T
}

type RPCFindValueResponse struct {
	RPCType int
	SenderID kademliaid.T
	Value []byte
}

type RPCStore struct {
	RPCType int
	SenderID kademliaid.T
	StoreValue []byte
}

func (nw *T) Listen(ip string, port int) {
	b := make([]byte, 2048)
	addrStr := fmt.Sprintf("%s:%d", ip, port)
	laddr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		log.Fatalf("Error listening on %v: %v\n", laddr, err)
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalf("Error listening on %v: %v\n", laddr, err)
	}
	nw.conn = conn
	for {
		_, raddr, err := conn.ReadFromUDP(b)
		if err != nil {
			log.Printf("Error reading UDP: %v", err)
			continue
		}
		go nw.resolveRPC(b, raddr)
	}
	// unreachable
}

func (nw *T) send(c *contact.T, msg []byte) (*net.UDPConn, error) {
	// TODO: contact should probably store the resolved address already
	// right now it's a string (who wrote this sample code?!)
	raddr, err := net.ResolveUDPAddr("udp", c.Address)
	if err != nil {
		return nil, err
		log.Printf("Could not resolve address %v: %v", c.Address, err)
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Printf("Error dialing %v: %v\n", raddr, err)
		// TODO: do we need to close conn here too? might depend on the error
		return nil, err
	}
	_, err = conn.Write(msg)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

func (nw *T) respond(msg interface{}, raddr *net.UDPAddr) error {
	b, err := msgpack.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling PingResponse: %v\n", err)
		return err
	}
	_, err = nw.conn.WriteTo(b, raddr)
	if err != nil {
		log.Printf("Error writing PingResponse: %v\n", err)
		return err
	}
	return nil
}

func (nw *T) receive(conn *net.UDPConn) ([]byte, error) {
	// TODO: make this buffer size configurable somewhere
	p := make([]byte, 2048)
	conn.SetReadDeadline(time.Now().Add(nw.timeout))
	_, err := bufio.NewReader(conn).Read(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (nw *T) rpc(c *contact.T, msg interface{}) (map[string]interface{}, error) {
	b, err := msgpack.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling FindNode RPC: %v\n", err)
		return nil, err
	}
	conn, err := nw.send(c, b)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	rb, err := nw.receive(conn)
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

func (nw *T) Ping(c *contact.T) error {
	msg := RPCPing{RPCType: PING, SenderID: *nw.id}
	_, err := nw.rpc(c, msg)
	if err != nil {
		return err
	}
	// routing table is updated as a side effect of receiving the response
	return nil
}

func (nw *T) FindNode(c *contact.T, findID *kademliaid.T) ([]contact.T, error) {
	msg := RPCFindNode{RPCType: FIND_NODE, SenderID: *nw.id, FindID: *findID}
	response, err := nw.rpc(c, msg)
	if err != nil {
		return nil, err
	}
	contacts := response["Contacts"].([]contact.T)
	return contacts, nil
}

func (nw *T) FindValue(c *contact.T, findID *kademliaid.T) ([]byte, []contact.T, bool, error) {
	msg := RPCFindValue{RPCType: FIND_VALUE, SenderID: *nw.id, FindID: *findID}
	// response, err := nw.rpc(c, msg)
	response, err := nw.rpc(c, msg) // just to mute errors until we use response
	if err != nil {
		return nil, nil, false, err
	}
	// TODO: remove assumption that it is int8 after marshalling
	rpcType := int(response["RPCType"].(int8))
	if rpcType == FIND_VALUE_RESPONSE {
		val := response["Value"].([]byte)
		return val, nil, true, nil
	} else {
		contacts := response["Contacts"].([]contact.T)
		return nil, contacts, false, nil
	}
}

func (nw *T) Store(c *contact.T, data []byte) error {
	msg := RPCStore{RPCType: STORE, SenderID: *nw.id, StoreValue: data}
	// not using rpc() since this rpc doesn't need a response
	b, err := msgpack.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling FindNode RPC: %v\n", err)
		return err
	}
	conn, err := nw.send(c, b)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func (nw *T) resolveRPC(message []byte, raddr *net.UDPAddr) {
	var args map[string]interface{}
	err := msgpack.Unmarshal(message, &args)
	if err != nil {
		log.Printf("Unable to unpack message from %v: %v\n", raddr, err)
		return
	}
	// This depends on what type msgpack decided to encode RPCType with and might break in the future.
	// TODO: read the RPC type in a safer way
	rpcType := int(args["RPCType"].(int8))
	switch rpcType {
	case PING:
		nw.pingResponse(raddr)
	case FIND_NODE:
		nw.findNodeResponse(&args, raddr)
	case FIND_VALUE:
		nw.findValueResponse(&args, raddr)
	case STORE:
		nw.storeResponse(&args)
	default:
		log.Printf("Unknown RPC: %v\n", args["RPCType"])
		// garbage message, don't update routing table
		return
	}
	// TODO: update our routing table
	// this should be a function RPC -> raddr -> contact
	// senderID := args["SenderID"].(kademliaid.T)
	// senderContact := contact.New(senderID, raddr)
	// nw.routingtable.Insert(senderContact)
}

func (nw *T) storeResponse(args *map[string]interface{}) {
	// TODO: use the kvstore Value type
	val := (*args)["Value"].([]byte)
	nw.kvstore.Store(val)
	// no confirmation is sent
}

func (nw *T) pingResponse(raddr *net.UDPAddr) {
	msg := RPCPingResponse{RPCType: PING_RESPONSE, SenderID: *nw.id}
	err := nw.respond(msg, raddr)
	if err != nil {
		log.Println("Failed to respond to ping: %v\n", err)
	}
}

func (nw *T) findValueResponse(args *map[string]interface{}, raddr *net.UDPAddr) {
	target := (*args)["FindID"].(kademliaid.T)
	val, ok := nw.kvstore.Get(target)
	if ok {
		// TODO: drop the pinned state
		msg := RPCFindValueResponse{RPCType: FIND_VALUE_RESPONSE, SenderID: *nw.id, Value: val}
		nw.respond(msg, raddr)
		return
	} else {
		// if we can't find it, just treat it as a FindNode
		nw.findNodeResponse(args, raddr)
	}
}

func (nw *T) findNodeResponse(args *map[string]interface{}, raddr *net.UDPAddr) {
	target := (*args)["FindID"].(kademliaid.T)
	contacts := nw.routingtable.FindKClosestContacts(&target)
	msg := RPCFindNodeResponse{RPCType: FIND_NODE_RESPONSE, SenderID: *nw.id, Contacts: contacts}
}
