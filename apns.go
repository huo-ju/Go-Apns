package goapns

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Notification struct {
	Device_token string
	Alert        string
	Badge        uint
	Sound        string
	Identifier   uint32
	Expiry       time.Duration
	Args         interface{}
}

func (n *Notification) MarshalJSON() ([]byte, error) {
	alert, _ := json.Marshal(n.Alert)
	badge, _ := json.Marshal(n.Badge)
	sound, _ := json.Marshal(n.Sound)
	args, _ := json.Marshal(n.Args)
	return []byte(fmt.Sprintf("{\"aps\":{\"alert\":%s,\"badge\":%s,\"sound\":%s},\"args\":%s}", alert, badge, sound, args)), nil
}

type NotificationError struct {
	Command    uint8
	Status     uint8
	Identifier uint32
}

type Apn struct {
	cert tls.Certificate
	server string

	conn      *tls.Conn
	Errorchan chan NotificationError
}

func Connect(cert_filename string, key_filename string, server string) (*Apn, error) {
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
	payloadbyte, _ := json.Marshal(notification)

	tokenbin, err := hex.DecodeString(notification.Device_token)
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
