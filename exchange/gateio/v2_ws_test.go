package gateio

import (
	"os"
	"testing"
	"time"

	"github.com/dqner/fym"
	"github.com/stretchr/testify/require"
)

func TestWebsocketClient_Ping(t *testing.T) {
	l, err := fym.NewZapLogger(fym.LogConf{
		Level:   "debug",
		Outputs: []string{"stdout"},
		Errors:  []string{"stderr"},
	},
	)
	require.NoError(t, err)
	defer l.Sync()
	client := new(WebsocketClient).Init("ws.gateio.ws", "/v4", l.Sugar())
	client.SetHandler(
		func() {
			l.Sugar().Debug("successfully connected")
			client.Ping(0)
		},
		func(resp interface{}) {
			pong, ok := resp.(string)
			if ok {
				l.Sugar().Debugf("response: %s", pong)
			} else {
				t.Error("wrong response")
			}
		},
	)
	client.Connect(true)

	time.Sleep(10 * time.Second)

	client.Close()
}

// go test -c ./exchange/gateio
// host=ws.gateio.ws apiKey=xxx secretKey=yyy ./gateio.test -test.run=TestWebsocketClient_Auth
func TestWebsocketClient_Auth(t *testing.T) {
	l, err := fym.NewZapLogger(fym.LogConf{
		Level:   "debug",
		Outputs: []string{"stdout"},
		Errors:  []string{"stderr"},
	},
	)
	require.NoError(t, err)
	defer l.Sync()
	apiKey := os.Getenv("apiKey")
	secretKey := os.Getenv("secretKey")
	host := os.Getenv("host")

	client := new(WebsocketClient).Init(host, "/v4", l.Sugar())
	client.SetHandler(
		func() {
			l.Sugar().Debug("successfully connected")
			client.Auth(apiKey, secretKey)
		},
		func(resp interface{}) {
			l.Sugar().Debugf("auth response: %v", resp)
		},
	)
	client.Connect(true)

	time.Sleep(10 * time.Second)

	client.Close()
}
