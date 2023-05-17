package core

import (
	"github.com/charmbracelet/charm/kv"
	"github.com/charmbracelet/log"
)

type App interface {
	DB() *kv.KV
	Logger() *log.Logger
}
