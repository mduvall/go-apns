package apns

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Client struct {
	ServerHost string
	Config     provisionConfig
}

type provisionConfig struct {
	AppId           string
	CertificatePath string
	Environment     string
}

func (c *Client) Configure(host string) {
	c.ServerHost = host
}

func (c *Client) Provision(appId string, certificatePath string, environment string) {
	c.Config = provisionConfig{AppId: appId, CertificatePath: certificatePath, Environment: environment}
	configByteSlice, err := json.Marshal(c.Config)
	if err != nil {
		log.Fatal("invalid provisioning configuration")
	}

	reader := bytes.NewBufferString(string(configByteSlice))
	_, err = http.Post(c.ServerHost+"/provision/", "application/json", reader)

	if err != nil {
		log.Fatal("provisioning was unsuccessful")
	}
}

func (c *Client) Notify(appId string, tokens []string, notifications map[string]string) {

}
