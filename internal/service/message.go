package service

import (
	"seabattle/internal/interfaces"
)

type message struct {
	Code string `json:"code"`
	Msg interface{} `json:"message"`
}

func NewMessage(code string, msg interface{}) interfaces.GameMessage {
	return &message{
		Code: code,
		Msg: msg,
	}
}

func NewEmptyMessage() interfaces.GameMessage {
	return &message{}
}

func (m *message) GetCode() string {
	return m.Code
}