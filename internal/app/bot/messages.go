package bot

import ()

type messager interface {
	GetMessage()
}

var (
	messages = map[string]string{
		"TYPENAME":    "Введите ваше имя",
		"INVALIDNAME": "Некорректное имя",
	}
)
