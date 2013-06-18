package apns

import (
	"log"
	"net/rpc"
	"strconv"
)

type Client struct {
	RpcClient *rpc.Client
}

func (c *Client) Configure(port int) {
	portString := strconv.Itoa(port)
	client, err := rpc.DialHTTP("tcp", ":"+portString)

	if err != nil {
		log.Fatal("unable to open client connection on localhost:" + portString)
	}

	c.RpcClient = client
}

func (c *Client) Provision(appId string, certificatePath string, environment string) {
	if c.RpcClient == nil {
		log.Fatal("configuration needs to be called first")
	}

	var reply int
	err := c.RpcClient.Call("Server.Provision", certificatePath, &reply)

	if err != nil {
		log.Fatal("provisioning was unsuccessful")
	}
}

func (c *Client) Notify(appId string, notification *Notification) {
	if c.RpcClient == nil {
		log.Fatal("configuration needs to be called first")
	}

	var reply int
	err := c.RpcClient.Call("Server.Notify", notification, &reply)

	if err != nil {
		log.Fatal("notification was unsuccessful")
	}
}
