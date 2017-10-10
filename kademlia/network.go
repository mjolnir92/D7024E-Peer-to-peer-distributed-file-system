package kademlia

import (
	"net"
	"fmt"
	"log"
	"time"
	"bufio"
	"github.com/mjolnir92/kdfs/kademliaid"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/kvstore"
	"github.com/mjolnir92/kdfs/constants"
	"github.com/vmihailenco/msgpack"
)

const (
	PING = 0
	PING_RESPONSE = 1
	FIND_NODE = 2
	FIND_NODE_RESPONSE = 3
	FIND_VALUE = 4
	FIND_VALUE_RESPONSE = 5
	STORE = 6
)

type RPCHeader struct {
	RPCType int
	Sender contact.T
}

type RPCPing struct {
	RPCType int
	Sender contact.T
}

type RPCPingResponse struct {
	RPCType int
	Sender contact.T
}

type RPCFindNode struct {
	RPCType int
	Sender contact.T
	FindID kademliaid.T
}

type RPCFindNodeResponse struct {
	RPCType int
	Sender contact.T
	Contacts []contact.T
}

type RPCFindValue struct {
	RPCType int
	Sender contact.T
	FindID kademliaid.T
}

type RPCFindValueResponse struct {
	RPCType int
	Sender contact.T
	ValueData []byte
	Contacts []contact.T
}

type RPCStore struct {
	RPCType int
	Sender contact.T
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
	conn.SetReadDeadline(time.Now().Add(constants.TIMEOUT))
	_, err := bufio.NewReader(conn).Read(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (nw *T) rpc(c *contact.T, msg interface{}, response interface{}) (error) {
	header, err := nw.rpcNoRefresh(c, msg, response)
	if err != nil {
		return err
	}
	if header != nil {
		nw.routingtable.AddContact(header.Sender)
	}
	return nil
}

func (nw *T) rpcNoRefresh(c *contact.T, msg interface{}, response interface{}) (*RPCHeader, error) {
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
	err = msgpack.Unmarshal(rb, response)
	if err != nil {
		return nil, err
	}
	// TODO: avoid unmarshalling twice somehow
	// this one is only used to update our routing table
	var header RPCHeader
	err = msgpack.Unmarshal(rb, &header)
	if err != nil {
		return nil, err
	}
	return &header, nil
}

func (nw *T) Ping(c *contact.T) error {
	msg := RPCPing{RPCType: PING, Sender: *nw.contactMe}
	var res RPCPingResponse
	err := nw.rpc(c, msg, &res)
	if err != nil {
		return err
	}
	// routing table is updated as a side effect of receiving the response
	return nil
}

func (nw *T) PingNoRefresh(c *contact.T) error {
	msg := RPCPing{RPCType: PING, Sender: *nw.contactMe}
	var res RPCPingResponse
	_, err := nw.rpcNoRefresh(c, msg, &res)
	if err != nil {
		return err
	}
	return nil
}

func (nw *T) FindNode(c *contact.T, findID *kademliaid.T) ([]contact.T, error) {
	msg := RPCFindNode{RPCType: FIND_NODE, Sender: *nw.contactMe, FindID: *findID}
	var res RPCFindNodeResponse
	err := nw.rpc(c, msg, &res)
	if err != nil {
		return nil, err
	}
	return res.Contacts, nil
}

// FindValue returns the value as []byte if it was found or some []contacts if it wasn't.
// The third return value is a bool that is true if the value was found.
func (nw *T) FindValue(c *contact.T, findID *kademliaid.T) ([]byte, []contact.T, bool, error) {
	msg := RPCFindValue{RPCType: FIND_VALUE, Sender: *nw.contactMe, FindID: *findID}
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
	msg := RPCStore{RPCType: STORE, Sender: *nw.contactMe, Value: *val}
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
	// We have to unmarshal the rest of the message after we know what type it is
	// TODO: find a way to unmarshal to the right type immediately
	var header RPCHeader
	err := msgpack.Unmarshal(message, &header)
	if err != nil {
		log.Printf("Unable to unpack message from %v: %v\n", raddr, err)
		return
	}
	switch header.RPCType {
	case PING:
		nw.pingResponse(raddr)
	case FIND_NODE:
		nw.findNodeResponse(message, raddr)
	case FIND_VALUE:
		nw.findValueResponse(message, raddr)
	case STORE:
		nw.storeResponse(message)
	default:
		log.Printf("Unknown RPC: %v\n", header.RPCType)
		// garbage message, don't update routing table
		return
	}
	nw.routingtable.AddContact(header.Sender)
}

func (nw *T) storeResponse(b []byte) {
	var msg RPCStore
	err := msgpack.Unmarshal(b, &msg)
	if err != nil {
		log.Printf("Failed to unmarshal into struct")
		return
	}
	
	id := kademliaid.NewHash(msg.Value.GetData())
	repub := func() {
		contacts := nw.LookupContact(id)
		for i := 0; i < len(contacts); i++ {
			go nw.Store(&contacts[i], &msg.Value)
		}
	}
	expire := func() {
		nw.eventmanager.DeleteEvent(*id, constants.REPUBLISH)
		nw.kvstore.Remove(msg.Value)
		nw.eventmanager.DeleteEvent(*id, constants.EXPIRE) //removes some garbage
	}

	//msg.Value will only be inserted if the timestamp is newer
	ok := nw.kvstore.Store(msg.Value)
	if ok {
		if msg.Value.GetPin() == true {
			nw.eventmanager.DeleteEvent(*id, constants.EXPIRE)
			nw.eventmanager.InsertEvent(*id, constants.REPUBLISH, repub, constants.REPUBLISH_TIME)
		} else {
			nw.eventmanager.InsertEvent(*id, constants.EXPIRE, expire, constants.EXPIRE_TIME)
			nw.eventmanager.InsertEvent(*id, constants.REPUBLISH, repub, constants.REPUBLISH_TIME)
		}
	} else {
		//If we didn't insert a new value, should we reset the republish time (efficient republishing?)
		//Perhaps compare time of current value and msg.Value, only reset if the message had the same or a newer timestamp
		nw.eventmanager.ResetEvent(*id, constants.REPUBLISH, constants.REPUBLISH_TIME) 
	}
}

func (nw *T) pingResponse(raddr *net.UDPAddr) {
	msg := RPCPingResponse{RPCType: PING_RESPONSE, Sender: *nw.contactMe}
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
		msg := RPCFindValueResponse{RPCType: FIND_VALUE_RESPONSE, Sender: *nw.contactMe, ValueData: val.GetData()}
		err := nw.respond(msg, raddr)
		if err != nil {
			log.Println("Failed to respond with value: %v\n", err)
		}
	} else {
		// if we can't find it, treat it like a FindNode RPC
		contacts := nw.routingtable.FindKClosestContacts(&msg.FindID)
		response := RPCFindValueResponse{RPCType: FIND_VALUE_RESPONSE, Sender: *nw.contactMe, Contacts: contacts}
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
	response := RPCFindNodeResponse{RPCType: FIND_NODE_RESPONSE, Sender: *nw.contactMe, Contacts: contacts}
	err = nw.respond(response, raddr)
	if err != nil {
		log.Println("Failed to respond with contacts: %v\n", err)
	}
}
