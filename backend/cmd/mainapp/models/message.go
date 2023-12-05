package models

type Message struct {
	Command string
	Body    interface{}
}

type MQResponse struct {
	Success bool
	Output  interface{}
}

type ErrorLog struct {
	Err string
}

type OutputLog struct {
	Output string
}
