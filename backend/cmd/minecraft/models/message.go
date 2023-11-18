package models

type TurnMCRequest struct {
	TurnOnOff bool
}

type Message struct {
	Command string
	Body    interface{}
}
