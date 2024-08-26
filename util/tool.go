package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/duke-git/lancet/v2/datetime"
	"strings"
	"time"
)

func SaltMd5(str string, salt string) string {
	b := []byte(str)
	s := []byte(salt)
	h := md5.New()
	h.Write(s) // 先写盐值
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1Sign(signStr, secretKey string) string {
	signStr = signStr + "&key=" + secretKey
	res := sha1.Sum([]byte(signStr))
	resStr := hex.EncodeToString(res[:])
	finalSign := strings.ToUpper(resStr)
	return finalSign
}

// 使用MD5的加解密
func Md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func StringToTime(timeStr string) time.Time {
	tm, _ := time.ParseInLocation(time.DateTime, timeStr, time.Local)
	return tm
}

func in() {
	datetime.AddDay(time.Now(), 24)
}
