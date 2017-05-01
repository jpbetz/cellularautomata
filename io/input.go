package io

type Input struct {
	Click chan Position
	Quit chan bool
	PausePlay chan bool
}
