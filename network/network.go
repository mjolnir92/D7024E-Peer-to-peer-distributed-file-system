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

func New(timeoutms int64, id *kademliaid.T, rt *routingtable.T, kvs *kvstore.T) T {
	// TODO: do more of the setup here
	return T{timeout: time.Duration(timeoutms) * time.Millisecond, id: id, routingtable: rt, kvstore: kvs}
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

// TODO: replace SenderID with SenderContact
// some messages are sent from a different port

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
	ValueData []byte
	Contacts []contact.T
}

type RPCStore struct {
	RPCType int
	SenderID kademliaid.T
	Value kvstore.Value
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

func (nw *T) rpc(c *contact.T, msg interface{}, response interface{}) (error) {
	b, err := msgpack.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling FindNode RPC: %v\n", err)
		return err
	}
	conn, err := nw.send(c, b)
	if err != nil {
		return err
	}
	defer conn.Close()
	rb, err := nw.receive(conn)
	if err != nil {
		return err
	}
	err = msgpack.Unmarshal(rb, response)
	if err != nil {
		return err
	}
	// TODO: update routing table
	return nil
}

func (nw *T) Ping(c *contact.T) error {
	msg := RPCPing{RPCType: PING, SenderID: *nw.id}
	var res RPCPingResponse
	err := nw.rpc(c, msg, &res)
	if err != nil {
		return err
	}
	// routing table is updated as a side effect of receiving the response
	return nil
}

func (nw *T) FindNode(c *contact.T, findID *kademliaid.T) ([]contact.T, error) {
	msg := RPCFindNode{RPCType: FIND_NODE, SenderID: *nw.id, FindID: *findID}
	var res RPCFindNodeResponse
	err := nw.rpc(c, msg, &res)
	if err != nil {
		return nil, err
	}
	//contacts := response["Contacts"].([]contact.T)
	return res.Contacts, nil
}

func (nw *T) FindValue(c *contact.T, findID *kademliaid.T) ([]byte, []contact.T, bool, error) {
	msg := RPCFindValue{RPCType: FIND_VALUE, SenderID: *nw.id, FindID: *findID}
	var res RPCFindValueResponse
	err := nw.rpc(c, msg, &res)
	if err != nil {
		return nil, nil, false, err
	}
	if len(res.ValueData) == 0 {
		// node did not have the key
		return nil, res.Contacts, false, nil
	}
	return res.ValueData, nil, true, nil
}

func (nw *T) Store(c *contact.T, val *kvstore.Value) error {
	msg := RPCStore{RPCType: STORE, SenderID: *nw.id, Value: *val}
	// not using rpc() since this rpc doesn't need a response
	b, err := msgpack.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling Store RPC: %v\n", err)
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
		nw.findNodeResponse(message, raddr)
	case FIND_VALUE:
		nw.findValueResponse(message, raddr)
	case STORE:
		nw.storeResponse(message)
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

func (nw *T) storeResponse(b []byte) {
	var msg RPCStore
	err := msgpack.Unmarshal(b, &msg)
	if err != nil {
		log.Printf("Failed to unmarshal into struct")
		return
	}
	// can't deserialize into private fields of the struct
	log.Printf("Hey I am trying to store %v", msg.Value)
	nw.kvstore.Store(msg.Value)
	// no confirmation is sent
}

func (nw *T) pingResponse(raddr *net.UDPAddr) {
	msg := RPCPingResponse{RPCType: PING_RESPONSE, SenderID: *nw.id}
	err := nw.respond(msg, raddr)
	if err != nil {
		log.Println("Failed to respond to ping: %v\n", err)
	}
}

func (nw *T) findValueResponse(b []byte, raddr *net.UDPAddr) {
	var msg RPCFindValue
	err := msgpack.Unmarshal(b, &msg)
	if err != nil {
		log.Printf("Failed to unmarshal into struct")
		return
	}
	val, ok := nw.kvstore.Get(msg.FindID)
	if ok {
		msg := RPCFindValueResponse{RPCType: FIND_VALUE_RESPONSE, SenderID: *nw.id, ValueData: val.GetData()}
		err := nw.respond(msg, raddr)
		if err != nil {
			log.Println("Failed to respond with value: %v\n", err)
		}
	} else {
		// if we can't find it, just treat it as a FindNode
		contacts := nw.routingtable.FindKClosestContacts(&msg.FindID)
		response := RPCFindValueResponse{RPCType: FIND_VALUE_RESPONSE, SenderID: *nw.id, Contacts: contacts}
		err = nw.respond(response, raddr)
		if err != nil {
			log.Println("Failed to respond with contacts: %v\n", err)
		}
	}
}

func (nw *T) findNodeResponse(b []byte, raddr *net.UDPAddr) {
	var msg RPCFindNode
	err := msgpack.Unmarshal(b, &msg)
	if err != nil {
		log.Printf("Failed to unmarshal into struct")
		return
	}
	contacts := nw.routingtable.FindKClosestContacts(&msg.FindID)
	response := RPCFindNodeResponse{RPCType: FIND_NODE_RESPONSE, SenderID: *nw.id, Contacts: contacts}
	err = nw.respond(response, raddr)
	if err != nil {
		log.Println("Failed to respond with contacts: %v\n", err)
	}
}
