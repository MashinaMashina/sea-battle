package service

import (
	"fmt"
	"log"
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
		p.SendAll(NewMessage("wait_opponent", "Ожидание подключения второго игрока"))
	} else {
		p.SendAll(NewMessage("game_start", "Начало игры"))
	}

	return nil
}

func (p *pair) SendOpponent(m interfaces.GameMessage) error {
	return nil
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

func (p *pair) Remove(first bool) error {
	if first {
		p.user1 = nil
	} else {
		p.user2 = nil
	}

	log.Println("User disconnected")

	return nil
}