package apns

import (
	"log"
	"net/rpc"
	"strconv"
)

type Client struct {
	Client *rpc.Client
}

func (c *Client) Configure(port int) {
	portString := strconv.Itoa(port)
	client, err := rpc.DialHTTP("tcp", ":"+portString)

	if err != nil {
		log.Fatal("unable to open client connection on localhost:" + portString)
	}

	c.Client = client
}

func (c *Client) Provision(appId string, certificatePath string, environment string) {
	if c.Client == nil {
		log.Fatal("configuration needs to be called first")
	}

	var reply int
	err := c.Client.Call("Server.Provision", certificatePath, &reply)

	if err != nil {
		log.Fatal("provisioning was unsuccessful")
	}
}

func (c *Client) Notify(appId string, notification *Notification) {
	if c.Client == nil {
		log.Fatal("configuration needs to be called first")
	}

	var reply int
	err := c.Client.Call("Server.Notify", notification, &reply)

	if err != nil {
		log.Fatal("notification was unsuccessful")
	}
}
