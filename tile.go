package main

type TileType uint8

const (
	Empty TileType = iota
	Brick
	Steel
	Trees
	Water
	Ice
	Base
)
