package errs

import "errors"

var (
	NeedToAuth = errors.New("Необходима авторизация")

	ErrNoUser       = errors.New("Пользователя с данным id не существует")
	ErrNoTrainer    = errors.New("Тренера с данным id не существует")
	InvalidEmail    = errors.New("Пользователя с такой почтой не существует")
	InvalidPassword = errors.New("Пароль не верен")
	ErrAlreadyExist = errors.New("Сущность уже существует")
)
