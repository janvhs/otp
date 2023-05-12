//go:build !cgo

package core

// INFO: 	This file is a slightly modified copy of the file linked below.
//			The original source code is licensed under the MIT License linked below.
// Source:	https://github.com/pocketbase/pocketbase/blob/c6d599244239ed17b2f2f7ce892b1279ddabf5ac/core/db_nocgo.go
// License:	https://github.com/pocketbase/pocketbase/blob/c6d599244239ed17b2f2f7ce892b1279ddabf5ac/LICENSE.md

import (
	"github.com/pocketbase/dbx"
	_ "modernc.org/sqlite"
)

func ConnectDB(dbPath string) (*dbx.DB, error) {
	db, err := dbx.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := initPragmas(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
