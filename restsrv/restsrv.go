package restsrv

import (
	"github.com/gin-gonic/gin"
	"github.com/mjolnir92/kdfs/restmsg"
	"fmt"
	"net/http"
)

func main() {
	router := gin.Default()
	// prefix everything with /v1
	v1 := router.Group("/v1")
	{
		v1.POST("/store", storeEndpoint)
		v1.GET("/store/:id", getEndpoint)
		v1.POST("/pin/:id", pinEndpoint)
		v1.DELETE("/pin/:id", unpinEndpoint)
	}
}

// POST /store
func storeEndpoint(c *gin.Context) {
	var req restmsg.StoreRequest
	err := c.MustBindWith(&req, gin.binding.MsgPack)
	if err != nil {
		b, err := msgpack.Marshal(restmsg.GenericResponse{Status: http.StatusBadRequest, Message: "Can't read the data"})
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal response: %v", err))
		}
		c.Data(http.StatusOK, gin.binding.MIMEMSGPACK2, b)
		return
	}
	id := kd.KademliaStore(&req.File)
	b, err := msgpack.Marshal(StoreResponse{Status: http.StatusOK, Message: "Success", ID: id.String()})
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal response: %v", err))
	}
	c.Data(http.StatusOK, gin.binding.MIMEMSGPACK2, b)
}

// GET /store/:id
func getEndpoint(c *gin.Context) {
	var id string = c.Param("id")
	kid := kademliaid.New(id)
	file := kd.Cat(kid)
	b, err := msgpack.Marshal(restmsg.CatResponse{Status: http.StatusOK, Message: "Success", File: file})
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal response: %v", err))
	}
	c.Data(http.StatusOK, gin.binding.MIMEMSGPACK2, b)
}

// POST /pin/:id
func pinEndpoint(c *gin.Context) {
	var id string = c.Param("id")
	kid := kademliaid.New(id)
	kd.Pin(kid)
	b, err := msgpack.Marshal(restmsg.GenericResponse{Status: http.StatusOK, Message: "Success"})
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal response: %v", err))
	}
	c.Data(http.StatusOK, gin.binding.MIMEMSGPACK2, b)
}

// DELETE /pin/:id
func unpinEndpoint(c *gin.Context) {
	var id string = c.Param("id")
	kid := kademliaid.New(id)
	kd.Unpin(kid)
	b, err := msgpack.Marshal(restmsg.GenericResponse{Status: http.StatusOK, Message: "Success"})
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal response: %v", err))
	}
	c.Data(http.StatusOK, gin.binding.MIMEMSGPACK2, b)
}
