package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"seabattle/internal/interfaces"
)

type pair struct {
	user1 interfaces.GameClient
	user2 interfaces.GameClient
}

func NewPair() *pair {
	return &pair{}
}

func (p *pair) AddClient(client interfaces.GameClient) error {
	if ! p.IsFree() {
		return fmt.Errorf("is not free pair")
	}

	client.AddPair(p)

	if p.user1 == nil {
		client.SetFirst(true)
		p.user1 = client
	} else if p.user2 == nil {
		client.SetFirst(false)
		p.user2 = client
	}

	if p.IsFree() {
		log.Infoln("User connected. Waiting opponent")
		p.SendAll(NewMessage("wait_opponent_connect", "Ожидание подключения второго игрока"))
	} else {
		log.Infoln("User connected. Start a game")
		p.SendAll(NewMessage("begin_ships_install", "Начало игры"))
	}

	return nil
}

func (p *pair) GetOpponent(iamIsFirst bool) interfaces.GameClient {
	if iamIsFirst {
		return p.user2
	} else {
		return p.user1
	}
}

func (p *pair) SendOpponent(iamIsFirst bool, m interfaces.GameMessage) error {
	if iamIsFirst && p.user2 != nil {
		p.user2.Send(m)
		return nil
	} else if !iamIsFirst && p.user1 != nil {
		p.user1.Send(m)
		return nil
	} else {
		return fmt.Errorf("dont have connected users")
	}
}

func (p *pair) SendAll(m interfaces.GameMessage) error {
	if p.user1 != nil {
		p.user1.Send(m)
	}
	if p.user2 != nil {
		p.user2.Send(m)
	}
	return nil
}

func (p *pair) IsFree() bool {
	return p.user1 == nil || p.user2 == nil
}

func (p *pair) HasOpponent() bool {
	return !p.IsFree()
}

func (p *pair) Remove(first bool) error {
	if first {
		p.user1 = nil
	} else {
		p.user2 = nil
	}

	log.Println("User disconnected")

	return nil
}