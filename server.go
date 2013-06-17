package apns

import (
	"crypto/tls"
	"log"
)

const (
	APNS_SERVER_SANDBOX_HOSTNAME     = "gateway.sandbox.push.apple.com"
	APNS_SERVER_HOSTNAME             = "gateway.push.apple.com"
	APNS_SERVER_PORT                 = "2195"
	FEEDBACK_SERVER_SANDBOX_HOSTNAME = "feedback.sandbox.push.apple.com"
	FEEDBACK_SERVER_HOSTNAME         = "feedback.push.apple.com"
	FEEDBACK_SERVER_PORT             = "2196"
)

// Provides generic configuration, provision of cert, and notification API
type server struct {
	FeedbackHost string
	FeedbackPort int
	APNSService  *service
}

// Service is responsible for certificate handling
// and implementing rw protocol for notification
type service struct {
	Host        string
	Port        string
	Certificate string
	Connection  *tls.Conn
}

func NewServer(environment string, filePath string) (createdServer *server, err error) {
	createdServer = &server{}
	host := getEnvironment(environment)
	createdService, err := createdServer.newService(filePath, host)
	if err != nil {
		return nil, err
	}
	createdServer.APNSService = createdService

	return createdServer, nil
}

// Opens a TLS connection with the certificate
func (*server) initializeConnection(s *service) {
	cert, err := tls.LoadX509KeyPair(s.Certificate, s.Certificate)
	if err != nil {
		log.Fatal("service is unable to load key at path ", s.Certificate)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}

	conn, err := tls.Dial("tcp", s.Host+":"+s.Port, &config)
	if err != nil {
		log.Fatal("service is unable to connect to host at ", s.Host)
	}

	s.Connection = conn
}

func (s *server) Write(notification *Notification) (err error) {
	notificationByteSlice := notification.constructBytePackage()
	conn := s.APNSService.Connection

	if _, err = conn.Write(notificationByteSlice); err != nil {
		log.Fatal("service was unable to send notification")
		return err
	}

	return nil
}

func (s *server) newService(filePath string, host string) (createdService *service, err error) {
	createdService = &service{Certificate: filePath, Host: host, Port: APNS_SERVER_PORT}
	s.initializeConnection(createdService)

	return createdService, nil
}

func getEnvironment(environment string) (host string) {
	if environment == "production" {
		host = APNS_SERVER_HOSTNAME
	} else {
		host = APNS_SERVER_SANDBOX_HOSTNAME
	}

	return host
}
