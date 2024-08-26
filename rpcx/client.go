package rpcx

import (
	"github.com/kawasoki/gowork/accerror"
	"github.com/kawasoki/gowork/gconf"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"sync"
	"time"
)

var (
	baseRpcxOnce sync.Once
)
var (
	baseRpcxClient client.XClient
)

func GetBaseXClient() client.XClient {
	baseRpcxOnce.Do(func() {
		d, _ := client.NewPeer2PeerDiscovery(gconf.GConf.MobileBaseRpcxUrl, "")
		baseRpcxClient = client.NewXClient("BaseService", client.Failtry, client.RandomSelect, d, client.DefaultOption)
		client.ClientErrorFunc = func(res *protocol.Message, e string) client.ServiceError {
			err, e1 := accerror.MewErrorString(e)
			if e1 == nil {
				return err
			}
			return client.NewServiceError(e)
		}

		// baseRpcxClient.Auth("bearer abcdefg1234567")

	})
	return baseRpcxClient
}

// 熔断器    CircuitBreaker Breaker = circuit.NewRateBreaker(0.95, 100)
// if failed 5 times, return error immediately, and will try to connect after 30 seconds
var breaker = func() client.Breaker {
	return client.NewConsecCircuitBreaker(5, 30*time.Second)

}
