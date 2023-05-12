package core

// INFO: 	This file is a direct copy of the file linked below.
//			The original source code is licensed under the MIT License linked below.
// Source:	https://github.com/pocketbase/pocketbase/blob/c6d599244239ed17b2f2f7ce892b1279ddabf5ac/core/db.go
// License:	https://github.com/pocketbase/pocketbase/blob/c6d599244239ed17b2f2f7ce892b1279ddabf5ac/LICENSE.md

import "github.com/pocketbase/dbx"

func initPragmas(db *dbx.DB) error {
	// note: the busy_timeout pragma must be first because
	// the connection needs to be set to block on busy before WAL mode
	// is set in case it hasn't been already set by another connection
	_, err := db.NewQuery(`
		PRAGMA busy_timeout       = 10000;
		PRAGMA journal_mode       = WAL;
		PRAGMA journal_size_limit = 200000000;
		PRAGMA synchronous        = NORMAL;
		PRAGMA foreign_keys       = TRUE;
	`).Execute()

	return err
}
