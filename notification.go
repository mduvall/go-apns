package apns

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"time"
)

type Notification struct {
	Token      string
	Payload    []byte
	Identifier int
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
