package goapns

import (
    "fmt"
    "encoding/json"
    "encoding/hex"
    "net"
    "bytes"
    "time"
    "encoding/binary"
    "crypto/tls"
)

type Notification struct {
    Device_token string
    Alert string
    Sound string
    Identifier uint32
    Expiry time.Duration
}
type NotificationError struct{
    Command uint8
    Status uint8
    Identifier uint32
}

type Apn struct{
    Conn *tls.Conn
    Errorchan chan NotificationError
}

func Connect(cert_filename string,key_filename string, server string) (*Apn, error){
    rchan := make(chan NotificationError)

    cert,cert_err := tls.LoadX509KeyPair(cert_filename,key_filename)

    if cert_err != nil {
        return &Apn{nil,nil},cert_err
    }

    conn, err:= net.Dial("tcp",server)

    if err != nil {
        return &Apn{nil,nil},err
    }

    Certificate := []tls.Certificate{cert}
    conf := tls.Config{
        Certificates:Certificate,
    }

    var client_conn *tls.Conn = tls.Client(conn,&conf)
    err = client_conn.Handshake()
    if err != nil {
        return &Apn{nil,nil},err
    }

    go readError(client_conn, rchan)


    return &Apn{client_conn,rchan},nil
}

func (apnconn *Apn) SendNotification(notification *Notification) error{
        fmt.Println(notification.Alert);
        var payload map[string](map[string]string)
        payload = make(map[string](map[string]string))
        payload["aps"] =  make(map[string]string)
        payload["aps"]["alert"] = notification.Alert
        payload["aps"]["sound"]=notification.Sound
        payloadbyte, _:= json.Marshal(payload)

        tokenbin, err := hex.DecodeString(notification.Device_token)
        if err != nil {
            return err
        }

        buffer := bytes.NewBuffer([]byte{})
        binary.Write(buffer, binary.BigEndian, uint8(1))
        binary.Write(buffer, binary.BigEndian, uint32(notification.Identifier))
        binary.Write(buffer, binary.BigEndian, uint32(time.Second + notification.Expiry))
        binary.Write(buffer, binary.BigEndian, uint16(len(tokenbin)))
        binary.Write(buffer, binary.BigEndian, tokenbin)
        binary.Write(buffer, binary.BigEndian, uint16(len(payloadbyte)))
        binary.Write(buffer, binary.BigEndian, payloadbyte)
        pushPackage := buffer.Bytes()

        _, err = apnconn.Conn.Write(pushPackage)
        return err
}

func readError(client_conn *tls.Conn,c chan NotificationError){
    var readb []byte
    readb = make([]byte, 6, 6)
    n, _:= client_conn.Read(readb)
    if n > 0 {
        notificationerror := NotificationError{}
        notificationerror.Command = uint8(readb[0])
        notificationerror.Status = uint8(readb[1])
        notificationerror.Identifier = uint32(readb[2])<<24+uint32(readb[3])<<16+uint32(readb[4])<<8+uint32(readb[5])
        c <- notificationerror
    }
}

func init() {
}
