package core

import (
	"bode.fun/2fa/log"
	"github.com/charmbracelet/charm/kv"
)

type App interface {
	DB() *kv.KV
	Logger() log.Logger
}
