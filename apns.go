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

type Notification struct {
	DeviceToken        string
	Identifier         uint32
	ExpireAfterSeconds int

	Payload *Payload
}

type NotificationError struct {
	Command    uint8
	Status     uint8
	Identifier uint32
}

// An Apn contain a ErrorChan channle when connected to apple server. When a notification sent wrong, you can get the error infomation from this channel.
type Apn struct {
	ErrorChan <-chan NotificationError

	cert      tls.Certificate
	server    string
	conn      *tls.Conn
	errorChan chan<- NotificationError
}

// Connect to apple server, with cert_filename and key_filename.
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

	return &Apn{rchan, cert, server, client_conn, rchan}, nil
}

// Reconnect if connection break.
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

	go readError(client_conn, apnconn.errorChan)

	return nil
}

// Send a notificatioin to iOS
func (apnconn *Apn) SendNotification(notification *Notification) error {
	tokenbin, err := hex.DecodeString(notification.DeviceToken)
	if err != nil {
		return err
	}

	payloadbyte, _ := json.Marshal(notification.Payload)
	expiry := time.Now().Add(time.Duration(notification.ExpireAfterSeconds) * time.Second).Unix()

	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, uint8(1))
	binary.Write(buffer, binary.BigEndian, uint32(notification.Identifier))
	binary.Write(buffer, binary.BigEndian, uint32(expiry))
	binary.Write(buffer, binary.BigEndian, uint16(len(tokenbin)))
	binary.Write(buffer, binary.BigEndian, tokenbin)
	binary.Write(buffer, binary.BigEndian, uint16(len(payloadbyte)))
	binary.Write(buffer, binary.BigEndian, payloadbyte)
	pushPackage := buffer.Bytes()

	_, err = apnconn.conn.Write(pushPackage)
	return err
}

func readError(client_conn *tls.Conn, c chan<- NotificationError) {
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
