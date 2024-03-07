package main

import (
	"context"
	"encoding/hex"
	rpc "github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/celestiaorg/celestia-node/share"
	"github.com/gin-gonic/gin"
	"github.com/rollkit/celestia-da/celestia"
	"net/http"
)

var (
	rpcAddr = ":8080"
)

func startRPCInfoServer(authToken, rpcToken, rpcAddress, nsString string, gasPrice float64) {
	r := gin.Default()
	r.Use(Cors())
	da := NewRPCCelestiaDa(authToken, rpcAddress, nsString, gasPrice)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/nodeID", func(c *gin.Context) {
		nodeID, _ := da.GetPeerInfo(c)
		c.JSON(http.StatusOK, gin.H{
			"nodeID": nodeID,
		})
	})

	//todo: protect the auth-token
	r.GET("/rpcToken", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"rpcToken": rpcToken,
		})
	})

	r.Run(rpcAddr)
}

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}
}

func NewRPCCelestiaDa(authToken, rpcAddress, nsString string, gasPrice float64) *celestia.CelestiaDA {
	client, err := rpc.NewClient(context.Background(), rpcAddress, authToken)
	if err != nil {
		log.Fatalln("failed to create celestia-node RPC client:", err)
	}
	nsBytes := make([]byte, len(nsString)/2)
	_, err = hex.Decode(nsBytes, []byte(nsString))
	if err != nil {
		log.Fatalln("invalid hex value of a namespace:", err)
	}
	namespace, err := share.NewBlobNamespaceV0(nsBytes)
	if err != nil {
		log.Fatalln("invalid namespace:", err)
	}

	da := celestia.NewCelestiaDA(client, namespace, gasPrice, context.Background())
	return da
}
