package logger

import (
	"fmt"
	cfg "github.com/tendermint/tendermint/config"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"sync"
)

type TestLogger struct {
	InfoMsgs  []string
	ErrorMsgs []string
	mtx       sync.Mutex
	// Context   *server.Context
	Logger log.Logger
}

func NewTestLogger() *TestLogger {
	tmLog := log.NewTMLogger(os.Stdout)
	tmLog, _ = tmflags.ParseLogLevel("main:info,state:info,*:error", tmLog, "error")
	return &TestLogger{Logger: tmLog}
	// return &TestLogger{}
}

func (tl *TestLogger) Debug(string, ...interface{}) {
}

func (tl *TestLogger) Info(msg string, keyvals ...interface{}) {
	tl.mtx.Lock()
	defer tl.mtx.Unlock()
	tl.InfoMsgs = append(tl.InfoMsgs, fmt.Sprintf(msg, keyvals...))
	tl.Logger.Info(msg, keyvals...)
}

func (tl *TestLogger) Error(msg string, keyvals ...interface{}) {
	tl.mtx.Lock()
	defer tl.mtx.Unlock()
	tl.ErrorMsgs = append(tl.ErrorMsgs, fmt.Sprintf(msg, keyvals...))
	tl.Logger.Error(msg, keyvals...)
}

func (tl *TestLogger) With(keyvals ...interface{}) log.Logger {
	// tl.Logger = tl.Logger.With(keyvals...)
	return tl
}

type Context struct {
	Config *cfg.Config
	Logger log.Logger
}
