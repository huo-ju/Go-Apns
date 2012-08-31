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

// An Apn contain a ErrorChan channle when connected to apple server. When a notification sent wrong, you can get the error infomation from this channel.
type Apn struct {
	ErrorChan <-chan NotificationError

	server    string
	conf      *tls.Config
	conn      *tls.Conn
	errorChan chan<- NotificationError
}

// Connect to apple server, with cert_filename and key_filename.
func Connect(cert_filename, key_filename, server string) (*Apn, error) {
	rchan := make(chan NotificationError)

	cert, err := tls.LoadX509KeyPair(cert_filename, key_filename)
	if err != nil {
		return nil, err
	}

	certificate := []tls.Certificate{cert}
	conf := &tls.Config{
		Certificates: certificate,
	}

	ret := &Apn{
		ErrorChan: rchan,
		server:    server,
		conf:      conf,
		errorChan: rchan,
	}

	err = ret.Reconnect()
	return ret, err
}

// Reconnect if connection break.
func (apnconn *Apn) Reconnect() error {
	// make sure last readError(...) will fail when reading.
	if apnconn.conn != nil {
		apnconn.conn.Close()
	}

	conn, err := net.Dial("tcp", apnconn.server)
	if err != nil {
		return err
	}

	var client_conn *tls.Conn = tls.Client(conn, apnconn.conf)
	err = client_conn.Handshake()
	if err != nil {
		return err
	}

	apnconn.conn = client_conn
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

func readError(conn *tls.Conn, c chan<- NotificationError) {
	p := make([]byte, 6, 6)
	for {
		n, err := conn.Read(p)
		e := NewNotificationError(p[:n], err)
		c <- e
		if err != nil {
			return
		}
	}
}
