package tlog

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

func LogT() {
	fmt.Println(time.Now().Clock())
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("hello, world", "user", "/Name?test")
}
