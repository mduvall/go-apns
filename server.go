package apns

import (
	"os"
)

const (
	APNS_SERVER_SANDBOX_HOSTNAME     = "gateway.sandbox.push.apple.com"
	APNS_SERVER_HOSTNAME             = "gateway.push.apple.com"
	APNS_SERVER_PORT                 = 2195
	FEEDBACK_SERVER_SANDBOX_HOSTNAME = "feedback.sandbox.push.apple.com"
	FEEDBACK_SERVER_HOSTNAME         = "feedback.push.apple.com"
	FEEDBACK_SERVER_PORT             = 2196
)

// Provides generic configuration, provision of cert, and notification API
type server struct {
	Host         string
	Port         int
	FeedbackHost string
	FeedbackPort int
	APNSService  *service
}

// Service is responsible for certificate handling
// and implementing rw protocol for notification
type service struct {
	Certificate *os.File
}

type serviceProtocol interface {
	Write() bool
	Read() string
}

func NewServer(environment string, filePath string) (createdServer *server, err error) {
	host := getEnvironment(environment)
	createdService, err := newService(filePath)
	if err != nil {
		return nil, err
	}

	return &server{Host: host, Port: APNS_SERVER_PORT, APNSService: createdService}, nil
}

func newService(filePath string) (createdService *service, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &service{Certificate: file}, nil
}

func getEnvironment(environment string) (host string) {
	if environment == "production" {
		host = APNS_SERVER_HOSTNAME
	} else {
		host = APNS_SERVER_SANDBOX_HOSTNAME
	}

	return host
}
