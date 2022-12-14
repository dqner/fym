package huobi

import (
	"github.com/dqner/fym/exchange"
	"github.com/huobirdcenter/huobi_golang/logging/applogger"
	"github.com/huobirdcenter/huobi_golang/pkg/client/marketwebsocketclient"
	"github.com/huobirdcenter/huobi_golang/pkg/client/websocketclientbase"
	"github.com/huobirdcenter/huobi_golang/pkg/model/market"
)

func (c *Client) SubscribeTrade(symbol, clientId string, responseHandler exchange.TradeHandler) {
	hb := new(marketwebsocketclient.TradeWebSocketClient).Init(c.Host)
	hb.SetHandler(
		// Connected handler
		func() {
			hb.Subscribe(symbol, clientId)
		},
		tradeHandler(responseHandler),
	)

	hb.Connect(true)
}

func (c *Client) UnsubscribeTrade(symbol, clientId string) {
	hb := new(marketwebsocketclient.TradeWebSocketClient).Init(c.Host)
	hb.UnSubscribe(symbol, clientId)
}

func tradeHandler(responseHandler exchange.TradeHandler) websocketclientbase.ResponseHandler {
	return func(response interface{}) {
		depthResponse, ok := response.(market.SubscribeTradeResponse)
		if ok {
			if &depthResponse != nil {
				if depthResponse.Tick != nil && depthResponse.Tick.Data != nil {
					applogger.Info("WebSocket received trade update: count=%d", len(depthResponse.Tick.Data))

					t := depthResponse.Tick.Data[0]

					responseHandler(exchange.TradeDetail{
						Id:        t.TradeId,
						Price:     t.Price,
						Amount:    t.Amount,
						Timestamp: t.Timestamp,
						Direction: t.Direction,
					})
				}

				if depthResponse.Data != nil {
					applogger.Info("WebSocket received trade data: count=%d", len(depthResponse.Data))
					//for _, t := range depthResponse.Data {
					//	applogger.Info("Trade data, id: %d, price: %v, amount: %v", t.TradeId, t.Price, t.Amount)
					//}
				}
			}
		} else {
			applogger.Warn("Unknown response: %v", response)
		}
	}
}
