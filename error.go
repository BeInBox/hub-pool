package POOL

import "errors"

var ErrorPoolKeyStillExist = errors.New("key still exist into the pool")
var ErrorPoolKeyNotExist = errors.New("key doesn't exist into the pool")
var ErrorPoolEntryMismatch = errors.New("unable to cast obj to pool entry")
