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

func (c *Client) Notify(appId string, tokens []string, notifications map[string]string) {

}
