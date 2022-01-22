package service

import (
	"github.com/gorilla/websocket"
	"log"
	"seabattle/internal/interfaces"
	"seabattle/internal/utils/mylogger"
)

type client struct {
	ws *websocket.Conn
	pair interfaces.GamePair
	isFirst bool
	msgCh chan interfaces.GameMessage
	doneCh chan bool
}

func NewClient(ws *websocket.Conn) interfaces.GameClient {
	msgCh := make(chan interfaces.GameMessage)
	doneCh := make(chan bool, 1)
	c := &client{
		ws: ws,
		msgCh: msgCh,
		doneCh: doneCh,
		isFirst: true,
	}

	mylogger.PrintGoroutines("NewClient")

	c.Listen()

	return c
}

func (c *client) Send(m interfaces.GameMessage) {
	c.msgCh <- m
}

func (c *client) AddPair(pair interfaces.GamePair)  {
	c.pair = pair
}

func (c *client) SetFirst(isFirst bool) {
	c.isFirst = isFirst
}

func (c *client) onMessage(m interfaces.GameMessage) {
	log.Printf("%+v", m)
}

func (c *client) Listen() {
	go c.readMessages()
	go c.writeMessages()
}

func (c *client) readMessages() {
	defer c.ws.Close()

	mylogger.PrintGoroutines("readMessages")

	for {
		m := NewEmptyMessage()
		err := c.ws.ReadJSON(m)

		if err != nil {
			log.Println(err)
			c.doneCh <- true
			log.Println("Closed readMessages #1")
			return
		} else {
			c.onMessage(m)
		}
	}
}

func (c *client) writeMessages() {

	mylogger.PrintGoroutines("writeMessages")

	for {
		select {
		case <- c.doneCh:
			c.pair.Remove(c.isFirst)
			c.doneCh <- true // for readMessages
			log.Println("Closed writeMessages #1")
			return
		case m := <- c.msgCh:
			if err := c.ws.WriteJSON(m); err != nil {
				log.Println("writeMessages err: ", err)
			}
		}
	}
}


