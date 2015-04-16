package handler

import (
	"fmt"
	"github.com/tomyhero/go-tcp_server/authorizer"
	"github.com/tomyhero/go-tcp_server/context"
	"net"
)

type EchoHandler struct {
	Authorizer context.IAuthorizer
}

// 他のハンドラーとぶつからないように、ユーニークのプレフィックス
func (h *EchoHandler) Prefix() string {
	return "echo"
}

// 認証インスタンス取得
func (h *EchoHandler) GetAuthorizer() context.IAuthorizer {
	return h.Authorizer
}

// コンストラクター。クライアントは、利用する際に、パスワード 1111 を必須にするサンプル
func NewEchoHandler() *EchoHandler {
	return &EchoHandler{Authorizer: authorizer.PlainPassword{Password: "1111"}}
}

// サーバ初期化時のフック
func (h *EchoHandler) HookInitialize(database map[string]interface{}) {
}

// サーバ停止時のフック
func (h *EchoHandler) HookDestroy(database map[string]interface{}) {
}

func (h *EchoHandler) HookDisconnect(conn net.Conn) {

}

// アクション実行前に実行されるフック。セッションに数を足しているサンプル
func (h *EchoHandler) HookBeforeExecute(c *context.Context) {
	//fmt.Println("Called BeforeExecuteHandler", c.Session, c.Database["echo"])
	session := c.Session
	_, ok := session["num"]
	if !ok {
		session["num"] = 0
	}
	session["num"] = session["num"].(int) + 1
}

// アクション実行後に実行されるフック。
func (h *EchoHandler) HookAfterExecute(c *context.Context) {
	//fmt.Println("Called AfterExecuteHandler")
}

//-----------------------
// 関数名を Action**** にすることにより、実行アクションを簡単に追加できる。
//------------------------

// echo_Echo コマンドで実行される。
func (h *EchoHandler) ActionEcho(c *context.Context) (*context.Context, error) {
	c.Res.Body = c.Req.Body
	fmt.Println(c.Req.Body["text"], c.Session["num"].(int))
	//fmt.Println(c.Req.Body["text"])

	for conn, _ := range c.Conns {
		err := c.CDataManager.Send(conn, c.Res.GetData())
		if err != nil {
			fmt.Println("Fail to send : %s", err)
			//glog.Warningf("Fail to send : %s", err)
		}
	}
	c.OnSendResponse = false
	return c, nil
}
