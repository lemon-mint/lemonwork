package lemonwork

type NetPoll interface {
	PollSize() int
	Poll() error

	Register(fd int) error

	SetOnCloseCallback(func(fd int))
	SetOnDataCallback(func(fd int))

	SetAutoClose(bool)

	ClosePoller() error
}
