package client

import (
	. "github.com/tomyhero/go-tcp_server/context"
	"net"
)

type Client struct {
	conn         net.Conn
	CDataManager *CDataManager
}

func (c *Client) Disconnect() {
	c.conn.Close()
}
func (c *Client) Connect(to string) error {

	conn, err := net.Dial("tcp", to)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Send(req *CData) error {
	err := c.CDataManager.Send(c.conn, req.GetData())
	return err
}

func (c *Client) Receive() (*CData, error) {
	data, err := c.CDataManager.Receive(c.conn)
	if err != nil {
		return nil, err
	}

	// TODO validate data format

	res := &CData{
		Header: data["H"].(map[string]interface{}),
		Body:   data["B"].(map[string]interface{}),
	}

	return res, nil
}
