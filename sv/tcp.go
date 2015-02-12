package sv

import (
	"fmt"
	"github.com/ugorji/go/codec"
	"net"
	"reflect"
	"time"
)

type TCPServerConfig struct {
	Port             int
	CodecHandle      codec.Handle
	DeadLineMillisec time.Duration
}

func (config *TCPServerConfig) GetPort() int {
	return config.Port
}
func (config *TCPServerConfig) GetCodecHandle() codec.Handle {
	return config.CodecHandle
}

func (config *TCPServerConfig) GetDeadLineMillisec() time.Duration {
	return config.DeadLineMillisec
}

func (config *TCPServerConfig) LoadDefault() {
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

type TCPNetwork struct {
	Config *TCPServerConfig
	ln     *net.TCPListener
}

func (n *TCPNetwork) GetConfig() IConfig {
	return n.Config
}

func (n *TCPNetwork) Close() {
	n.ln.Close()
}

func (n *TCPNetwork) Listen(port int) error {
	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}
	n.ln = ln
	return nil
}

func (n *TCPNetwork) Accept() (net.Conn, error) {
	// TODO
	n.ln.SetDeadline(time.Now().Add(n.Config.GetDeadLineMillisec() * time.Millisecond))
	conn, err := n.ln.Accept()
	return conn, err
}
