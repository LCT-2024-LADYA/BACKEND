package errs

import "errors"

var (
	NeedToAuth = errors.New("Необходима авторизация")

	ErrNoRole           = errors.New("Роли с данным id не существует")
	ErrNoSpecialization = errors.New("Специализации с данным id не существует")
	ErrNoAchievement    = errors.New("Достижения с данным id не существует")
	ErrNoService        = errors.New("Достижения с данным id не существует")
	ErrNoUser           = errors.New("Пользователя с данным id не существует")
	ErrNoTrainer        = errors.New("Тренера с данным id не существует")
	InvalidEmail        = errors.New("Пользователя с такой почтой не существует")
	InvalidPassword     = errors.New("Пароль не верен")
	ErrAlreadyExist     = errors.New("Сущность уже существует")
)
