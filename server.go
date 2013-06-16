package apns

import (
	"os"
	"path"
)

const (
	APNS_SERVER_SANDBOX_HOSTNAME     = "gateway.sandbox.push.apple.com"
	APNS_SERVER_HOSTNAME             = "gateway.push.apple.com"
	APNS_SERVER_PORT                 = 2195
	FEEDBACK_SERVER_SANDBOX_HOSTNAME = "feedback.sandbox.push.apple.com"
	FEEDBACK_SERVER_HOSTNAME         = "feedback.push.apple.com"
	FEEDBACK_SERVER_PORT             = 2196
)

type server struct {
	Host         string
	Port         int
	FeedbackHost string
	FeedbackPort int
}

// Service is responsible for certificate handling and implementing rw protocol
type Service struct {
	Certificate *File
}

type ServiceProtocol interface {
	Write() bool
	Read() string
}

func NewServer(environment string, path string) *server {
	host := getEnvironment(environment)

	return &server{Host: host, Port: APNS_SERVER_PORT}
}

func getEnvironment(environment string) (host string) {
	if environment == "production" {
		host = APNS_SERVER_HOSTNAME
	} else {
		host = APNS_SERVER_SANDBOX_HOSTNAME
	}

	return host
}
