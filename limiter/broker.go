package limiter

import (
	"errors"
	"fmt"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"log"
	"time"
)

// RequestLimiter 包装了限流器和熔断器的结构
type RequestLimiter struct {
	limiter *rate.Limiter
	circuit *gobreaker.CircuitBreaker
}

type UserFunc func() (interface{}, error)

// An Option configures a mutex.
type Option interface {
	Apply(*gobreaker.Settings)
}

// OptionFunc is a function that configures a mutex.
type OptionFunc func(*gobreaker.Settings)

// Apply calls f(mutex)
func (f OptionFunc) Apply(mutex *gobreaker.Settings) {
	f(mutex)
}

// WithTimeOut 熔断持续时间
func WithTimeOut(expiry time.Duration) Option {
	return OptionFunc(func(m *gobreaker.Settings) {
		m.Timeout = expiry
	})
}

// WithInterval 熔断失败率持续时间
func WithInterval(expiry time.Duration) Option {
	return OptionFunc(func(m *gobreaker.Settings) {
		m.Interval = expiry
	})
}

// WithMaxRequests 熔断后半开启 请求数
func WithMaxRequests(rq uint32) Option {
	return OptionFunc(func(m *gobreaker.Settings) {
		m.MaxRequests = rq
	})
}

func NewRequestLimier(name string, options ...Option) *RequestLimiter {
	st := gobreaker.Settings{
		Name: name,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Println(name, from.String(), to.String())
		},
	}
	if len(options) > 0 {
		for _, o := range options {
			o.Apply(&st)
		}
	}
	circuit := gobreaker.NewCircuitBreaker(st)
	limiter := rate.NewLimiter(rate.Limit(10), 5)
	return &RequestLimiter{
		circuit: circuit,
		limiter: limiter,
	}
}

var count int

func (rl *RequestLimiter) Broker(req UserFunc) (interface{}, error) {
	// 使用速率限制器进行限流
	count++
	if !rl.limiter.Allow() {
		log.Println("Request limited by rate limiter")
		return nil, errors.New("request limited by rate limiter")
	}
	data, err := rl.circuit.Execute(req)
	fmt.Println(count, "请求数:", rl.circuit.Counts().Requests, ",失败数:", rl.circuit.Counts().TotalFailures)
	return data, err
}
