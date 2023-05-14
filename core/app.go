package core

import (
	"bode.fun/otp/log"
	"github.com/pocketbase/dbx"
)

type App interface {
	DB() *dbx.DB
	Logger() log.Logger
}
