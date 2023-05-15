package core

import (
	"bode.fun/2fa/log"
	"github.com/pocketbase/dbx"
)

type App interface {
	DB() *dbx.DB
	Logger() log.Logger
}
