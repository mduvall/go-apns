package apns

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type Notification struct {
	Token      string
	Payload    *Payload
	Identifier string
}

func (n *Notification) constructBytePackage() []byte {
	tokenbin, err := hex.DecodeString(n.Token)
	if err != nil {
		log.Fatal("invalid device token")
	}

	expiry := time.Now().Add(time.Duration(0) * time.Second).Unix()
	identifier, err := strconv.Atoi(n.Identifier)
	payloadJson, err := json.Marshal(n.Payload)

	if err != nil {
		log.Fatal("payload json is not valid")
	}

	buff := bytes.NewBuffer([]byte{})
	binary.Write(buff, binary.BigEndian, uint8(1))
	binary.Write(buff, binary.BigEndian, uint32(identifier))
	binary.Write(buff, binary.BigEndian, uint32(expiry))
	binary.Write(buff, binary.BigEndian, uint16(len(tokenbin)))
	binary.Write(buff, binary.BigEndian, tokenbin)
	binary.Write(buff, binary.BigEndian, uint16(len(payloadJson)))
	binary.Write(buff, binary.BigEndian, payloadJson)

	return buff.Bytes()
}
