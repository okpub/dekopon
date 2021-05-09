package cluster

type Handler func(interface{}) (interface{}, error)

type Server interface {
	Disconnect()
}
