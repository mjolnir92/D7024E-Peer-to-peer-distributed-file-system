package network

import (
	"net"
	"fmt"
	"log"
	"time"
	"github.com/mjolnir92/kdfs/kademliaid"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/routingtable"
	"github.com/vmihailenco/msgpack"
)

type T struct {
	id_local *kademliaid.T
	routingtable *routingtable.T
	timeout time.Duration
}

func New(timeoutms int64, id_local *kademliaid.T) T {
	return T{timeout: timeoutms * time.Millisecond, id_local: id_local}
}

const (
	PING = 0
	FIND_NODE = 1
	FIND_VALUE = 2
	STORE = 3
)

type RPC struct {
	rpc_type int
	sender_id kademliaid.T
}

type RPCPing struct {
	RPC
}

type RPCFindNode struct {
	RPC
	find_id kademliaid.T
}

type RPCFindValue struct {
	RPC
	find_id kademliaid.T
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

func (nw *T) send(c *contact.T, msg []byte) (net.UDPConn, error) {
	// TODO: contact should probably store the resolved address already
	// right now it's a string (who wrote this sample code?!)
	raddr := net.ResolveUDPAddr("udp", c.Address)
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

func (nw *T) receive(conn net.UDPConn) ([]byte, error) {
	// TODO: make this buffer size configurable somewhere
	p := make([]byte, 2048)
	conn.SetReadDeadline(time.Now().Add(nw.timeout))
	_, err = bufio.NewReader(conn).Read(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (nw *T) rpc(c *contact.T, msg *RPC) (map[string]interface{}, error) {
	b, err := msgpack.Marshal(&msg)
	if err != nil {
		log.Printf("Error marshalling FindNode RPC: %v\n", err)
		return nil, err
	}
	conn, err := send(c, b)
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

func (nw *T) Ping(c *contact.T) error {
	msg := &RPCPing{rpc_type: PING, sender_id: nw.id}
	response, err := rpc(c, msg)
	if err != nil {
		return err
	}
	// routing table is updated as a side effect of receiving the response
	return nil
}

func (nw *T) FindNode(c *contact.T, findID *kademliaid.T) ([]contact.T, error) {
	msg := &RPCFindNode{rpc_type: FIND_NODE, sender_id: nw.id, find_id: *findID}
	response, err := rpc(c, msg)
	if err != nil {
		return err
	}
	// TODO: do something with response and return
}

func (nw *T) FindValue(c *contact.T, findID *kademliaid.T) {
	msg := &RPCFindValue{rpc_type: FIND_VALUE, sender_id: nw.id, find_id: *findID}
	response, err := rpc(c, msg)
	if err != nil {
		return err
	}
	// TODO: do something with response and return
}

func (nw *T) Store(data []byte) error {
	msg := &RPCStore{rpc_type: STORE, sender_id: nw.id, store_value: data}
	// not using rpc() since this rpc doesn't need a response
	b, err := msgpack.Marshal(&msg)
	if err != nil {
		log.Printf("Error marshalling FindNode RPC: %v\n", err)
		return err
	}
	conn, err := send(c, b)
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

func (nw *T) pingResponse(raddr *net.UDPAddr) {
}

func (nw *T) findValueResponse(raddr *net.UDPAddr)) {
	// TODO: try to find it in the kv store
	// return
	// if we can't find it, just treat it as a FindNode
	nw.findNodeResponse(c, findID)
}

func (nw *T) findNodeResponse(raddr *net.UDPAddr) {
	target := (*args)["find_id"].(kademliaid.T)
	nw.routingtable.FindClosestContacts(&target)
}
