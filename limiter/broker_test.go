package limiter

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestName(t *testing.T) {

}

func LimiterService(appid int) {
	rl := NewRequestLimier("HTTP GET", WithTimeOut(2*time.Second), WithMaxRequests(1))
	data, err := rl.Broker(GetGameData(appid))
	fmt.Println(data, err)
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetGameData(appid int) func() (interface{}, error) {
	return func() (interface{}, error) {
		// 在这里使用 appid 进行操作
		td := r.Intn(100)
		if td < appid {
			fmt.Println("error")
			return nil, errors.New("error")
		}
		return []byte("success"), nil
	}
}
