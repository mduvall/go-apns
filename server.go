package apns

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
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
type Server struct {
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

func StartServer(environment string, port int) (err error) {
	createdServer := &Server{}
	host := getEnvironment(environment)
	createdService, err := createdServer.newService(host)
	if err != nil {
		return err
	}

	createdServer.APNSService = createdService
	createdServer.setupRPC(port)

	return nil
}

func (s *Server) setupRPC(port int) {
	rpc.Register(s)
	rpc.HandleHTTP()

	portString := strconv.Itoa(port)
	listener, err := net.Listen("tcp", ":"+portString)
	if err != nil {
		log.Fatal("binding rpc to port 8080 failed")
	}

	http.Serve(listener, nil)
}

func (s *Server) Provision(certificatePath string, reply *int) error {
	s.APNSService.Certificate = certificatePath
	s.initializeConnection()

	return nil
}

func (s *Server) Notify(notification *Notification, reply *int) error {
	if s.APNSService.Certificate == "" {
		log.Fatal("the certificate needs to be provisioned by the client")
	}

	s.write(notification)

	return nil
}

// Opens a TLS connection with the certificate
func (s *Server) initializeConnection() {
	service := s.APNSService
	cert, err := tls.LoadX509KeyPair(service.Certificate, service.Certificate)
	if err != nil {
		log.Fatal("service is unable to load key at path ", service.Certificate)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}

	conn, err := tls.Dial("tcp", service.Host+":"+service.Port, &config)
	if err != nil {
		log.Fatal("service is unable to connect to host at ", service.Host)
	}

	service.Connection = conn
}

func (s *Server) write(notification *Notification) (err error) {
	// Arg is of type notification
	notificationByteSlice := notification.constructBytePackage()
	conn := s.APNSService.Connection

	if _, err = conn.Write(notificationByteSlice); err != nil {
		log.Fatal("service was unable to send notification")
		return err
	}

	return nil
}

func (s *Server) newService(host string) (createdService *service, err error) {
	createdService = &service{Host: host, Port: APNS_SERVER_PORT}

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
