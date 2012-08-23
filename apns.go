package goapns

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"net"
	"time"
)

type SimpleAps struct {
	Alert string `json:"alert"`
	Badge uint   `json:"badge"`
	Sound string `json:"sound"`
}

type Notification struct {
	DeviceToken string
	Identifier  uint32
	Aps         SimpleAps
	CustomFiels map[string]interface{}
	Expiry      time.Duration
}

func (n *Notification) MarshalJSON() ([]byte, error) {
	model := make(map[string]interface{})
	model["aps"] = n.Aps

	if n.CustomFiels != nil {
		for key, value := range n.CustomFiels {
			model[key] = value
		}
	}

	return json.Marshal(model)
}

type NotificationError struct {
	Command    uint8
	Status     uint8
	Identifier uint32
}

type Apn struct {
	cert   tls.Certificate
	server string

	conn      *tls.Conn
	Errorchan chan NotificationError
}

func Connect(cert_filename, key_filename, server string) (*Apn, error) {
	rchan := make(chan NotificationError)

	cert, cert_err := tls.LoadX509KeyPair(cert_filename, key_filename)
	if cert_err != nil {
		return nil, cert_err
	}

	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}

	certificate := []tls.Certificate{cert}
	conf := tls.Config{
		Certificates: certificate,
	}

	var client_conn *tls.Conn = tls.Client(conn, &conf)
	err = client_conn.Handshake()
	if err != nil {
		return nil, err
	}

	go readError(client_conn, rchan)

	return &Apn{cert, server, client_conn, rchan}, nil
}

func (apnconn *Apn) Reconnect() error {
	conn, err := net.Dial("tcp", apnconn.server)
	if err != nil {
		return err
	}

	certificate := []tls.Certificate{apnconn.cert}
	conf := tls.Config{
		Certificates: certificate,
	}

	var client_conn *tls.Conn = tls.Client(conn, &conf)
	err = client_conn.Handshake()
	if err != nil {
		return err
	}

	go readError(client_conn, apnconn.Errorchan)

	return nil
}

func (apnconn *Apn) SendNotification(notification *Notification) error {
	payloadbyte, _ := notification.MarshalJSON()

	tokenbin, err := hex.DecodeString(notification.DeviceToken)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, uint8(1))
	binary.Write(buffer, binary.BigEndian, uint32(notification.Identifier))
	binary.Write(buffer, binary.BigEndian, uint32(time.Second+notification.Expiry))
	binary.Write(buffer, binary.BigEndian, uint16(len(tokenbin)))
	binary.Write(buffer, binary.BigEndian, tokenbin)
	binary.Write(buffer, binary.BigEndian, uint16(len(payloadbyte)))
	binary.Write(buffer, binary.BigEndian, payloadbyte)
	pushPackage := buffer.Bytes()

	_, err = apnconn.conn.Write(pushPackage)
	return err
}

func readError(client_conn *tls.Conn, c chan NotificationError) {
	var readb []byte
	readb = make([]byte, 6, 6)
	n, _ := client_conn.Read(readb)
	if n > 0 {
		notificationerror := NotificationError{}
		notificationerror.Command = uint8(readb[0])
		notificationerror.Status = uint8(readb[1])
		notificationerror.Identifier = uint32(readb[2])<<24 + uint32(readb[3])<<16 + uint32(readb[4])<<8 + uint32(readb[5])
		c <- notificationerror
	}
}

func init() {
}
