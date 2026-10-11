package singleton

// До лабораторной 4 (авторизация) пользователи заданы константами.
// Все методы API получают текущего пользователя только через эти функции.

// GetCreatorID — текущий представитель команды (создатель заявок)
func GetCreatorID() uint {
	return 1
}

// GetModeratorID — организатор соревнования (модератор)
func GetModeratorID() uint {
	return 2
}
