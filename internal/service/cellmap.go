package service

import (
	"fmt"
	"github.com/buger/jsonparser"
	"seabattle/internal/interfaces"
	"strconv"
)

type MapCell struct {
	hasShip bool
	alive bool
}

func (c *MapCell) AliveCell() bool {
	return c.hasShip && c.alive
}
func (c *MapCell) HasShip() bool {
	return c.hasShip
}

func (c *MapCell) SetAliveCell(b bool) {
	c.alive = b
}
func (c *MapCell) SetHasShip(b bool) {
	c.hasShip = b
}

type CellMap struct {
	cells map[int]map[int]interfaces.MapCell `json:"cells"`
}

func (c *CellMap) From(json []byte) error {
	c.cells = make(map[int]map[int]interfaces.MapCell)
	cells, _, _, err := jsonparser.Get(json, "cells")
	if err != nil {
		return err
	}

	jsonparser.ObjectEach(cells, func(key []byte, value []byte, _ jsonparser.ValueType, _ int) error {
		x, err := strconv.Atoi(string(key))
		if err != nil {
			return err
		}

		yLayer := make(map[int]interfaces.MapCell)
		jsonparser.ObjectEach(value, func(key []byte, value []byte, _ jsonparser.ValueType, _ int) error {
			y, err := strconv.Atoi(string(key))
			if err != nil {
				return err
			}

			yLayer[y] = &MapCell{
				hasShip: string(value) == "true",
				alive: true,
			}

			return nil
		})

		c.cells[x] = yLayer

		return nil
	})

	return c.Validate()
}

func (c *CellMap) ExistsShip(x, y int) bool {
	return c.cells[x][y].HasShip()
}

func (c *CellMap) Validate() error {
	needShips := 20
	haveShips := 0

	for _, yLayer := range c.cells {
		for _, val := range yLayer {
			if val.HasShip() {
				haveShips++
			}
		}
	}

	if haveShips != needShips {
		return fmt.Errorf("invalid ships count")
	}

	return nil
}

func (c *CellMap) DamageCell(x, y int) error {
	if _, ok := c.cells[x]; ok {
		if _, ok = c.cells[x][y]; ok {
			c.cells[x][y].SetAliveCell(false)
			return nil
		}
	}

	return fmt.Errorf("cell not found")
}

func (c *CellMap) HasAliveShips() bool {
	for _, yLayer := range c.cells {
		for _, val := range yLayer {
			if val.AliveCell() {
				return true
			}
		}
	}

	return false
}