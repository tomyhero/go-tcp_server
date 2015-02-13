package server

import (
	"fmt"
	"github.com/ugorji/go/codec"
	"net"
	"reflect"
	"time"
)

type UDPServerConfig struct {
	Port             int
	CodecHandle      codec.Handle
	DeadLineMillisec time.Duration
}

func (config *UDPServerConfig) GetPort() int {
	return config.Port
}
func (config *UDPServerConfig) GetCodecHandle() codec.Handle {
	return config.CodecHandle
}

func (config *UDPServerConfig) GetDeadLineMillisec() time.Duration {
	return config.DeadLineMillisec
}

func (config *UDPServerConfig) LoadDefault() {
	if config.Port == 0 {
		config.Port = 8081
	}

	if config.CodecHandle == nil {
		var h = new(codec.MsgpackHandle)
		h.MapType = reflect.TypeOf(map[string]interface{}{})
		h.RawToString = true
		config.CodecHandle = h
	}

	if config.DeadLineMillisec == 0 {
		config.DeadLineMillisec = 1000 // 1 sec
	}
}

// --------

type UDPNetwork struct {
	Config *UDPServerConfig
	conn   *net.UDPConn
}

func (n *UDPNetwork) GetConfig() IConfig {
	return n.Config
}

func (n *UDPNetwork) Close() {
	//
}

func (n *UDPNetwork) Listen(port int) error {
	laddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}
	n.conn = conn
	return nil
}

func (n *UDPNetwork) Accept() (net.Conn, error) {
	return n.conn, nil
}
