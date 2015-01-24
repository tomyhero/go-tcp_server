package server

/*
import (
	"bufio"
	"fmt"
	"github.com/golang/glog"
	"net"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	defer func() {
		glog.Infof("Stop Test Client")
		//close(quit)

		// ln.Close()
	}()

	sv := Server{Port: 8080}
	go sv.Run()

	time.Sleep(10 * time.Millisecond)

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		glog.Error(err.Error())
		return
	}
	defer conn.Close()
	fmt.Fprintf(conn, "client")
	res, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		glog.Error(err.Error())
		return
	}

	glog.Infof("Received: %s", res)
}
*/
