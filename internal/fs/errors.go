package fs

import "errors"

//Ошибки
var (
	ErrCopy          = errors.New("can not copy")
	ErrBlocked       = errors.New("blocked by another process")
	ErrAlreadyExists = errors.New("already exists")
	ErrNotExists     = errors.New("not exists")
	ErrNotRegular    = errors.New("not regular")
	ErrNotDirectory  = errors.New("not directory")
)
