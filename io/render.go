package io

type Renderer interface {
	Run()
	Close()
	Set(position Position, change Cell)
	Draw()
}
