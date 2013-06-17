package apns

import (
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	ServerHost string
}

func (c *Client) Configure(host string) {
	c.ServerHost = host
}

func (c *Client) Provision(appId string, certificatePath string, environment string) {
	postParams := make(url.Values)
	postParams.Set("appId", appId)
	postParams.Set("certificatePath", certificatePath)
	postParams.Set("environment", environment)

	_, err := http.PostForm(c.ServerHost+"/provision/", postParams)

	if err != nil {
		log.Fatal("provisioning was unsuccessful")
	}
}

func (c *Client) Notify(appId string, notification *Notification) {
	// @TODO(mduvall): consolidate these param settings, not sure what to use here
	postParams := make(url.Values)
	postParams.Set("appId", appId)
	postParams.Set("token", notification.Token)
	postParams.Set("payload", string(notification.Payload))
	postParams.Set("identifier", string(notification.Identifier))

	_, err := http.PostForm(c.ServerHost+"/notify/", postParams)

	if err != nil {
		log.Fatal("notification was unsuccessful")
	}
}
