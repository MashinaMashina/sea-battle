package service

import (
	json2 "encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"seabattle/internal/interfaces"
)

type message struct {
	Code string `json:"code"`
	Msg []byte `json:"-"`
}

func NewMessage(code string, msg interface{}) interfaces.GameMessage {
	mess := &message{
		Code: code,
	}

	mess.SetMessage(msg)

	return mess
}

func NewEmptyMessage() interfaces.GameMessage {
	return &message{}
}

func (m *message) GetCode() string {
	return m.Code
}
func (m *message) GetMessage() []byte {
	return m.Msg
}
func (m *message) SetCode(s string) {
	m.Code = s
}
func (m *message) SetMessage(b interface{}) {
	switch b.(type) {
	case string:
		m.Msg = []byte(fmt.Sprintf("\"%s\"", b.(string)))
	case []byte:
		m.Msg = b.([]byte)
	default:
		//log.Println(fmt.Sprintf("%T", b))
	}
}
func (m *message) GetJSON() ([]byte, error) {
	json, err := json2.Marshal(m)

	if err != nil {
		return nil, err
	}

	if m.Msg == nil {
		return json, nil
	}

	json, err = jsonparser.Set(json, m.Msg, "message")

	if err != nil {
		return nil, err
	}

	return json, nil
}