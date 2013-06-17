package apns

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"log"
	"time"
)

const (
	APNS_SERVER_SANDBOX_HOSTNAME     = "gateway.sandbox.push.apple.com"
	APNS_SERVER_HOSTNAME             = "gateway.push.apple.com"
	APNS_SERVER_PORT                 = "2195"
	FEEDBACK_SERVER_SANDBOX_HOSTNAME = "feedback.sandbox.push.apple.com"
	FEEDBACK_SERVER_HOSTNAME         = "feedback.push.apple.com"
	FEEDBACK_SERVER_PORT             = "2196"
)

type Notification struct {
	Token      string
	Payload    []byte
	Identifier int
}

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

func (n *Notification) constructBytePackage() []byte {
	tokenbin, err := hex.DecodeString(n.Token)
	if err != nil {
		log.Fatal("invalid device token")
	}

	expiry := time.Now().Add(time.Duration(0) * time.Second).Unix()

	buff := bytes.NewBuffer([]byte{})
	binary.Write(buff, binary.BigEndian, uint8(1))
	binary.Write(buff, binary.BigEndian, uint32(n.Identifier))
	binary.Write(buff, binary.BigEndian, uint32(expiry))
	binary.Write(buff, binary.BigEndian, uint16(len(tokenbin)))
	binary.Write(buff, binary.BigEndian, tokenbin)
	binary.Write(buff, binary.BigEndian, uint16(len(n.Payload)))
	binary.Write(buff, binary.BigEndian, n.Payload)

	return buff.Bytes()
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
