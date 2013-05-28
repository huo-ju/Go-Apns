package apns

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
	DeviceToken        string
	Identifier         uint32
	ExpireAfterSeconds int

	Payload *Payload
}

// An Apn contain a ErrorChan channle when connected to apple server. When a notification sent wrong, you can get the error infomation from this channel.
type Apn struct {
	ErrorChan <-chan error

	server  string
	conf    *tls.Config
	conn    *tls.Conn
	timeout time.Duration

	sendChan  chan *sendArg
	errorChan chan error
}

// New Apn with cert_filename and key_filename.
func New(cert_filename, key_filename, server string, timeout time.Duration) (*Apn, error) {
	echan := make(chan error)

	cert, err := tls.LoadX509KeyPair(cert_filename, key_filename)
	if err != nil {
		return nil, err
	}

	certificate := []tls.Certificate{cert}
	conf := &tls.Config{
		Certificates: certificate,
	}

	ret := &Apn{
		ErrorChan: echan,
		server:    server,
		conf:      conf,
		timeout:   timeout,
		sendChan:  make(chan *sendArg),
		errorChan: echan,
	}

	go sendLoop(ret)
	return ret, err
}

func (a *Apn) GetErrorChan() <-chan error {
	return a.ErrorChan
}

// Send a notification to iOS
func (a *Apn) Send(notification *Notification) error {
	err := make(chan error)
	arg := &sendArg{
		n:   notification,
		err: err,
	}
	a.sendChan <- arg
	return <-err
}

type sendArg struct {
	n   *Notification
	err chan<- error
}

func (a *Apn) Close() error {
	if a.conn == nil {
		return nil
	}
	conn := a.conn
	a.conn = nil
	return conn.Close()
}

func (a *Apn) connect() (<-chan int, error) {
	// make sure last readError(...) will fail when reading.
	err := a.Close()
	if err != nil {
		return nil, fmt.Errorf("close last connection failed: %s", err)
	}

	conn, err := net.Dial("tcp", a.server)
	if err != nil {
		return nil, fmt.Errorf("connect to server error: %d", err)
	}

	var client_conn *tls.Conn = tls.Client(conn, a.conf)
	err = client_conn.Handshake()
	if err != nil {
		return nil, fmt.Errorf("handshake server error: %s", err)
	}

	a.conn = client_conn
	quit := make(chan int)
	go readError(client_conn, quit, a.errorChan)

	return quit, nil
}

func (a *Apn) send(notification *Notification) error {
	tokenbin, err := hex.DecodeString(notification.DeviceToken)
	if err != nil {
		return fmt.Errorf("convert token to hex error: %s", err)
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

	_, err = a.conn.Write(pushPackage)
	if err != nil {
		return fmt.Errorf("write socket error: %s", err)
	}
	return nil
}

func sendLoop(apn *Apn) {
	for {
		arg := <-apn.sendChan
		quit, err := apn.connect()
		if err != nil {
			arg.err <- err
			continue
		}
		arg.err <- apn.send(arg.n)

		for connected := true; connected; {
			select {
			case <-quit:
				connected = false
			case <-time.After(apn.timeout):
				connected = false
			case arg := <-apn.sendChan:
				arg.err <- apn.send(arg.n)
			}
		}

		err = apn.Close()
		if err != nil {
			e := NewNotificationError(nil, err)
			apn.errorChan <- e
		}
	}
}

func readError(conn *tls.Conn, quit chan<- int, c chan<- error) {
	p := make([]byte, 6, 6)
	for {
		n, err := conn.Read(p)
		e := NewNotificationError(p[:n], err)
		c <- e
		if err != nil {
			quit <- 1
			return
		}
	}
}
