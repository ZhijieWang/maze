package common

type Observer interface {
	OnNotify(data interface{})
}
type Event interface {
}
type Notifier interface {
	Register(Observer)
	Deregister(Observer)
	Notify(Event)
}
type Simulation interface {
	Run(obs Observer) error
	Stop() bool
}
