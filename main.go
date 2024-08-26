package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Candlestick struct {
	OpenTime                 int64   `json:"openTime"`
	Open                     float64 `json:"open"`
	High                     float64 `json:"high"`
	Low                      float64 `json:"low"`
	Close                    float64 `json:"close"`
	Volume                   float64 `json:"volume"`
	CloseTime                int64   `json:"closeTime"`
	QuoteAssetVolume         float64 `json:"quoteAssetVolume"`
	NumberOfTrades           int     `json:"numberOfTrades"`
	TakerBuyBaseAssetVolume  float64 `json:"takerBuyBaseAssetVolume"`
	TakerBuyQuoteAssetVolume float64 `json:"takerBuyQuoteAssetVolume"`
}

func main1() {
	// 设置API地址
	apiURL := "https://api.binance.com/api/v3/klines?symbol=BTCUSDT&interval=1h&startTime=%d&endTime=%d"

	endTime := time.Now().Unix() * 1000
	startTime := endTime - (365 * 24 * 60 * 60 * 1000)

	// 格式化API URL
	requestURL := fmt.Sprintf(apiURL, startTime, endTime)
	client := &http.Client{Timeout: time.Second * 5}
	proxyURL, _ := url.Parse("http://127.0.0.1:1087")
	client.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	req, _ := http.NewRequest("GET", requestURL, nil)
	// 发送HTTP请求
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	// 读取响应数据
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 解析JSON响应
	var candlesticks [][]interface{}
	err = json.Unmarshal(body, &candlesticks)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 输出数据
	for _, c := range candlesticks {
		fmt.Println(c)
		//fmt.Printf("Time: %s, Open: %.2f, High: %.2f, Low: %.2f, Close: %.2f, Volume: %.2f\n",
		//	time.Unix(int64(c.[0]).(float64)/1000), 0).Format("2006-01-02 15:04:05"),
		//	c.Open, c.High, c.Low, c.Close, c.Volume)
	}
}

const (
	alphabet = "almnopdefZ012ghijktuvwxyzABCDLMqrcNOPFGHIJsbKUVWXY78QRSTE93456"
	base     = int64(len(alphabet))
)

func main() {
	//primes := findPrimes(10000, 99999)
	//fmt.Println("5位数的质数：", primes)
	userID := int64(9999999999999)
	//userID *= 95279
	shortCode := encode(userID)
	fmt.Println("Short code:", shortCode)

	decodedUserID := decode(shortCode)
	fmt.Println("Decoded user ID:", decodedUserID)
}

// 将用户 ID 编码成短码
func encode(userID int64) string {
	code := ""
	for userID > 0 {
		code = string(alphabet[userID%base]) + code
		userID /= base
	}
	return code
}

// 将短码解码成用户 ID
func decode(code string) int64 {
	userID := int64(0)
	for _, char := range code {
		index := strings.IndexRune(alphabet, char)
		userID = userID*base + int64(index)
	}
	return userID
}

func findPrimes(start, end int) []int {
	var primes []int

	for num := start; num <= end; num++ {
		if isPrime(num) {
			primes = append(primes, num)
		}
	}

	return primes
}

func isPrime(num int) bool {
	if num <= 1 {
		return false
	}
	if num <= 3 {
		return true
	}
	if num%2 == 0 || num%3 == 0 {
		return false
	}
	i := 5
	for i*i <= num {
		if num%i == 0 || num%(i+2) == 0 {
			return false
		}
		i += 6
	}
	return true
}
