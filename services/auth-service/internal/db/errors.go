package db

import "errors"

var ErrMissingDSN = errors.New("postgres DSN is empty")
