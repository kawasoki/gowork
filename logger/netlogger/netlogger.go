package netlogger

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"
)

type EncoderType int

const (
	headerLen      = uint32(4)
	logLevelLen    = uint32(2)
	Timeformat     = "2006-01-02 15:04:05.000"
	ConsoleEncoder = EncoderType(0)
	JsonEncoder    = EncoderType(1)
)

type LogAgent interface {
	Write(p []byte) (n int, err error)
	LogLevelStrToUint(text []byte) (uint16, error)
	EnCode(payload []byte) []byte
}
type LogAgentConf struct {
	ServerName  string
	AgentAddr   string
	ChanBuffer  int
	EncoderConf *zapcore.EncoderConfig
	EncoderType EncoderType
}
type pkg struct {
	buff  []byte
	total uint32
}
type ZapLoggerAgent struct {
	conf           *LogAgentConf
	logger         *zap.Logger
	bufferChan     chan *pkg
	c              net.Conn
	connAtomic     uint32 //1 connect 0:not connect
	tryConnRunning uint32 //1 running 0:not running
}

func (l *ZapLoggerAgent) Daemon() *ZapLoggerAgent {

	go func() {
		for pg := range l.bufferChan {
			_, err := l.c.Write(pg.buff[:pg.total])
			if err != nil {
				fmt.Printf(BytesToString(pg.buff[:pg.total]))
			}
			BUFFERPOOL.Put(pg.buff)
		}
	}()
	return l
}
func (l *ZapLoggerAgent) Conn() *ZapLoggerAgent {
	c, err := net.Dial("udp", l.conf.AgentAddr)
	if err != nil {
		log.Println(l.conf.AgentAddr, "udp连接失败:", err)
		atomic.CompareAndSwapUint32(&l.connAtomic, 1, 0)

		l.tryConn()
		return l
	}
	atomic.CompareAndSwapUint32(&l.connAtomic, 0, 1)
	l.c = c
	return l
}
func (l *ZapLoggerAgent) tryConn() {
	if atomic.LoadUint32(&l.tryConnRunning) == 1 {
		return
	}
	go func() {
		for {
			if atomic.LoadUint32(&l.connAtomic) == 1 {
				atomic.CompareAndSwapUint32(&l.tryConnRunning, 1, 0)
				return
			}
			atomic.CompareAndSwapUint32(&l.tryConnRunning, 0, 1)
			l.Conn()
			time.Sleep(time.Second * 5)
		}
	}()
}
func (l *ZapLoggerAgent) initLogger() *ZapLoggerAgent {
	if l.conf.EncoderConf == nil {
		l.conf.EncoderConf = &zapcore.EncoderConfig{
			MessageKey:       "message",
			LevelKey:         "level",
			EncodeLevel:      zapcore.CapitalLevelEncoder, // INFO
			TimeKey:          "time",
			EncodeTime:       zapcore.TimeEncoderOfLayout(Timeformat),
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: " ",
			CallerKey:        "f",
		}
	}
	encoder := zapcore.NewConsoleEncoder(*l.conf.EncoderConf)
	l.conf.EncoderConf.EncodeLevel = zapcore.CapitalLevelEncoder
	if l.conf.EncoderType == JsonEncoder {
		encoder = zapcore.NewJSONEncoder(*l.conf.EncoderConf)
	}
	w := zapcore.AddSync(l)
	core := zapcore.NewCore(encoder, w, zapcore.DebugLevel)
	l.logger = zap.New(core, zap.AddCaller())
	return l
}
func (l *ZapLoggerAgent) Init(config *LogAgentConf) *ZapLoggerAgent {
	if config == nil {
		panic("config nil")
	}
	if config.ServerName == "" {
		panic("ServerName invalid")
	}
	if config.AgentAddr == "" {
		panic("AgentAddr invalid")
	}
	if config.EncoderType < ConsoleEncoder || config.EncoderType > JsonEncoder {
		panic("invalid EncoderType")
	}
	l.conf = config
	if l.conf.ChanBuffer == 0 {
		l.conf.ChanBuffer = 1024
	}
	l.bufferChan = make(chan *pkg, l.conf.ChanBuffer)
	l.initLogger()
	return l
}
func (l *ZapLoggerAgent) Write(p []byte) (n int, err error) {
	if atomic.LoadUint32(&l.connAtomic) == 0 {
		fmt.Printf(BytesToString(p))
		return len(p), nil
	}
	buff, total := l.EnCode(p)
	pg := &pkg{
		buff:  buff,
		total: total,
	}
	select {

	case l.bufferChan <- pg:
	default:
		fmt.Printf(BytesToString(p))
	}
	return len(p), nil
}

func (l *ZapLoggerAgent) EnCode(payload []byte) ([]byte, uint32) {
	if l.conf.ServerName == "" {
		panic("ServerName invalid")
	}

	hl := uint32(len(l.conf.ServerName)) + logLevelLen
	total := uint32(len(payload)) + hl + headerLen
	ptr := BUFFERPOOL.Get(total)
	buf := *ptr
	//buf := make([]byte, uint32(len(payload))+hl+headerLen)
	binary.LittleEndian.PutUint32(buf, hl)
	binary.LittleEndian.PutUint16(buf[headerLen:], l.logLevelStrToUint(payload))
	copy(buf[headerLen+logLevelLen:], l.conf.ServerName)
	copy(buf[headerLen+hl:], payload)
	return buf, total
}
func (l *ZapLoggerAgent) logLevelStrToUint(text []byte) uint16 {
	if l.conf.EncoderConf.LevelKey == "" {
		return 0 //
	}
	if l.conf.EncoderConf.EncodeLevel == nil {
		return 0
	}

	var lvl zapcore.Level
	if l.conf.EncoderType == JsonEncoder {
		str := BytesToString(text)
		start := strings.Index(str, ":") + 2
		end := start
		for i := start; i < start+16; i++ {
			if str[i] == '"' {
				end = i
				break
			}
		}

		_ = lvl.UnmarshalText(text[start:end])
		return uint16(lvl) + 1
	}
	separator := StringToBytes(l.conf.EncoderConf.ConsoleSeparator)
	if l.conf.EncoderConf.TimeKey != "" {
		idx := bytes.Index(text, separator)
		idx++
		start := bytes.Index(text[idx:], separator)
		start++
		end := bytes.Index(text[start+idx:], separator)
		_ = lvl.UnmarshalText(text[start+idx : start+idx+end])
	} else {
		idx := bytes.Index(text, separator)
		_ = lvl.UnmarshalText(text[:idx])
	}
	//fmt.Println(lvl)
	return uint16(lvl) + 1
}

func (l *ZapLoggerAgent) Logger() *zap.Logger {
	return l.logger
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
