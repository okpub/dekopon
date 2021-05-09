package cluster

type Server interface {
	Start() error
	Disconnect()
	Publish(interface{}) error
	Request(interface{}) (interface{}, error)
}
