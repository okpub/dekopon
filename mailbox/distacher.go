package mailbox

var defaultThroughput = 0

type Dispatcher interface {
	Schedule(fn func())
	Throughput() int
}

type goroutineDispatcher int

func (goroutineDispatcher) Schedule(fn func()) {
	go fn()
}

func (d goroutineDispatcher) Throughput() int {
	return int(d)
}

func NewDefaultDispatcher() Dispatcher {
	return goroutineDispatcher(defaultThroughput)
}

type synchronizedDispatcher int

func (synchronizedDispatcher) Schedule(fn func()) {
	fn()
}

func (d synchronizedDispatcher) Throughput() int {
	return int(d)
}

func NewSynchronizedDispatcher() Dispatcher {
	return synchronizedDispatcher(defaultThroughput)
}
