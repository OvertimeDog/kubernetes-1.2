package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/glog"
	"net"
)

type dataHead struct {
	msgType uint32
	dataID  uint32
	utctime uint32
	bodylen uint32
	resv    [2]uint32
}

// Client gse data pipe client
type Client struct {

	// The domain socket address
	Endpoint string
	Conn     net.Conn
}

func (dhead *dataHead) packageData(data []byte) ([]byte, error) {

	buffer := new(bytes.Buffer)

	// package head
	if err := binary.Write(buffer, binary.BigEndian, dhead); nil != err {
		return nil, fmt.Errorf("can not package the data head,%v", err)
	}

	// package body
	if _, err := buffer.Write(data); nil != err {
		return nil, fmt.Errorf("can not package the data body ,%v", err)
	}

	return buffer.Bytes(), nil
}

// Connect connect to gse data pipe
func (gsec *Client) Connect() error {

	conn, err := net.Dial("unix", gsec.Endpoint)

	if err != nil {
		return fmt.Errorf("no gse data pipe  available, maybe gseagent is not running")
	}

	glog.V(3).Infof("current endpoint of gsedatapipe: %s", gsec.Endpoint)
	gsec.Conn = conn
	return nil
}

// Send write data to data pipe
func (gsec *Client) Send(dataid uint32, data []byte) error {

	if gsec.Conn == nil {
		if error := gsec.Connect(); error != nil {
			return error
		}
	}

	dhead := dataHead{0xc01, dataid, 0, uint32(len(data)), [2]uint32{0, 0}}

	packageData, err := dhead.packageData(data)
	if nil != err {
		return err
	}

	if _, error := gsec.Conn.Write(packageData); error != nil {
		glog.Errorf("fail send data to data pipe: %v", error)
		gsec.Close()
		return error
	}

	return nil
}

// Close close gse data pipe
func (gsec *Client) Close() error {
	gsec.Conn.Close()
	gsec.Conn = nil
	return nil
}

// New create a new client of gse data pipe
func New(endpoint string) (*Client, error) {

	return &Client{Endpoint: endpoint}, nil
}
