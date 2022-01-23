package service

import (
	"github.com/buger/jsonparser"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"seabattle/internal/interfaces"
	"time"
)

type client struct {
	ws *websocket.Conn
	pair interfaces.GamePair
	isFirst bool
	msgCh chan interfaces.GameMessage
	doneCh chan bool
	stage uint8
	CellMap *CellMap
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

	c.Listen()

	return c
}

func (c *client) GetCellMap() interfaces.GameCellMap {
	return c.CellMap
}

func (c *client) onShipsInstalled(m interfaces.GameMessage) {
	c.CellMap = &CellMap{}

	if err := c.CellMap.From(m.GetMessage()); err != nil {
		mess := NewMessage("error_ship_install", err.Error())
		c.Send(mess)
		return
	}

	c.SetStage(interfaces.StageReady)

	// TODO check if opponent's stage is valid
	if c.GetOpponent().GetStage() != interfaces.StageReady {
		c.Send(NewMessage("wait_opponent_ready", "Ожидание второго участника"))
	} else {
		c.SendAll(NewMessage("game_start", nil))
		c.gameStart()
	}
}

func (c *client) SuccessShoot(x, y int, m interfaces.GameMessage) {
	c.Send(NewMessage("success_shoot", m.GetMessage()))

	c.GetOpponent().GotSuccessShoot(x, y, m)
	c.GetOpponent().GetCellMap().DamageCell(x, y)

	if !c.GetOpponent().GetCellMap().HasAliveShips() {
		c.Win()
	}
}

func (c *client) GotSuccessShoot(x, y int, m interfaces.GameMessage) {
	c.Send(NewMessage("got_success_shoot", m.GetMessage()))
}

func (c *client) FailShoot(x, y int, m interfaces.GameMessage) {
	c.GetOpponent().GotFailShoot(x, y, m)

	c.Send(NewMessage("fail_shoot", m.GetMessage()))
}

func (c *client) GotFailShoot(x, y int, m interfaces.GameMessage) {
	c.Send(NewMessage("got_fail_shoot", m.GetMessage()))
}

func (c *client) onShoot(m interfaces.GameMessage) {
	if c.GetStage() != interfaces.StageShoot {
		c.Send(NewMessage("error", "You are not allowed to shoot"))
		return
	}

	x64, err := jsonparser.GetInt(m.GetMessage(), "x")
	if err != nil {
		log.Println(err)
		return
	}

	y64, err := jsonparser.GetInt(m.GetMessage(),"y")
	if err != nil {
		log.Println(err)
		return
	}

	x := int(x64)
	y := int(y64)

	if c.GetOpponent().GetCellMap().ExistsShip(x, y) {
		c.SuccessShoot(x, y, m)
	} else {
		c.FailShoot(x, y, m)

		c.GetOpponent().ShootStage()
	}
}

func (c *client) Win() {
	c.GetOpponent().Defeat()

	c.SetStage(interfaces.StageEnd)
	c.Send(NewMessage("win", nil))
}

func (c *client) Defeat() {
	c.SetStage(interfaces.StageEnd)
	c.Send(NewMessage("defeat", nil))
}

func (c *client) gameStart() {
	rand.Seed(time.Now().Unix())

	if rand.Intn(1) == 1 && c.isFirst {
		c.ShootStage()
	} else {
		c.GetOpponent().ShootStage()
	}
}

func (c *client) ShootStage() {
	c.SetStage(interfaces.StageShoot)
	c.Send(NewMessage("you_shoot", nil))

	c.GetOpponent().SetStage(interfaces.StageWaitOpponentFair)
	c.GetOpponent().Send(NewMessage("wait_opponent_shoot", nil))
}

func (c *client) onMessage(m interfaces.GameMessage) {
	c.Send(m)

	switch m.GetCode() {
	case "ships_installed":
		c.onShipsInstalled(m)
	case "shoot":
		c.onShoot(m)
	default:
		log.Println("Got unhandled message with code:", m.GetCode())
	}
}

func (c *client) Listen() {
	go c.readMessages()
	go c.writeMessages()
}

func (c *client) SetStage(stage uint8) {
	c.stage = stage
}

func (c *client) GetStage() uint8 {
	return c.stage
}

func (c *client) CloseConn() {
	c.doneCh <- true
}

func (c *client) Send(m interfaces.GameMessage) {
	c.msgCh <- m
}

func (c *client) SendAll(m interfaces.GameMessage) error {
	return c.pair.SendAll(m)
}

func (c *client) SendOpponent(m interfaces.GameMessage) error {
	return c.pair.SendOpponent(c.isFirst, m)
}

func (c *client) GetOpponent() interfaces.GameClient {
	return c.pair.GetOpponent(c.isFirst)
}

func (c *client) AddPair(pair interfaces.GamePair)  {
	c.pair = pair
}

func (c *client) SetFirst(isFirst bool) {
	c.isFirst = isFirst
}

func (c *client) readUserMessage(resMess interfaces.GameMessage) error {
	_, mess, err := c.ws.ReadMessage()

	//log.Println(string(mess))

	if err != nil {
		return err
	}

	code, err := jsonparser.GetString(mess, "code")

	if err != nil {
		return err
	}

	outerMess, _, _, err := jsonparser.Get(mess, "message")
	if err != nil {
		return err
	}

	resMess.SetMessage(outerMess)
	resMess.SetCode(code)

	return nil
}

func (c *client) readMessages() {
	defer c.ws.Close()

	for {
		m := NewEmptyMessage()
		err := c.readUserMessage(m)

		if err != nil {
			log.Println(err)
			c.CloseConn()
			log.Println("Closed readMessages #1")
			return
		} else {
			c.onMessage(m)
		}
	}
}

func (c *client) writeMessages() {
	for {
		select {
		case <- c.doneCh:
			c.pair.Remove(c.isFirst)
			c.doneCh <- true // for readMessages
			log.Println("Closed writeMessages #1")
			return
		case m := <- c.msgCh:
			json, err := m.GetJSON()
			if err != nil {
				log.Println("writeMessages err #1: ", err)
			}

			if err := c.ws.WriteMessage(websocket.TextMessage, json); err != nil {
				log.Println("writeMessages err #2: ", err)
			}
		}
	}
}


