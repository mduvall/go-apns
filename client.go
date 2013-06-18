package apns

import (
	"log"
	"net/rpc"
)

type Client struct {
	ServerHost string
	Client     *rpc.Client
}

func (c *Client) Configure(host string) {
	c.ServerHost = host
	client, err := rpc.DialHTTP("tcp", ":8080")

	if err != nil {
		log.Fatal("unable to open client connection on localhost:8080")
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
