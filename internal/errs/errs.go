package errs

import "errors"

var (
	NeedToAuth = errors.New("Необходима авторизация")
)
