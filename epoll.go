// +build linux

package lemonwork

import "syscall"

const epollet = 1 << 31

type EpollPoller struct {
	pollBufferSize int
	pollBuffer     []syscall.EpollEvent

	autoClose bool

	onCloseCallback func(fd int)
	onDataCallback  func(fd int)

	epollFd int
}

func NewEpollPoller(PollSize int) (ep *EpollPoller, err error) {
	ep = new(EpollPoller)
	ep.pollBufferSize = PollSize
	ep.pollBuffer = make([]syscall.EpollEvent, PollSize)
	ep.epollFd, err = syscall.EpollCreate1(0)
	ep.autoClose = true
	ep.onCloseCallback = func(fd int) {}
	ep.onDataCallback = func(fd int) {}
	return
}

func (ep *EpollPoller) Register(fd int) error {
	var event syscall.EpollEvent
	event.Fd = int32(fd)
	event.Events = syscall.EPOLLIN | syscall.EPOLLRDHUP | epollet
	return syscall.EpollCtl(ep.epollFd, syscall.EPOLL_CTL_ADD, fd, &event)
}

func (ep *EpollPoller) Poll() error {
	n, err := syscall.EpollWait(ep.epollFd, ep.pollBuffer, -1)
	if err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		event := ep.pollBuffer[i]
		fd := int(event.Fd)
		events := event.Events
		switch events {
		case syscall.EPOLLIN | syscall.EPOLLRDHUP:
			if ep.autoClose {
				syscall.Close(fd)
			}
			ep.onCloseCallback(fd)
		case syscall.EPOLLIN:
			ep.onDataCallback(fd)
		}
	}
	return nil
}

func (ep *EpollPoller) PollSize() int {
	return ep.pollBufferSize
}

func (ep *EpollPoller) SetOnCloseCallback(f func(fd int)) {
	ep.onCloseCallback = f
}

func (ep *EpollPoller) SetOnDataCallback(f func(fd int)) {
	ep.onDataCallback = f
}

func (ep *EpollPoller) SetAutoClose(x bool) {
	ep.autoClose = x
}

func (ep *EpollPoller) ClosePoller() error {
	return syscall.Close(ep.epollFd)
}

func NewPoller(PollSize int) (ep NetPoll, err error) {
	return NewEpollPoller(PollSize)
}
