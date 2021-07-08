// +build darwin dragonfly freebsd netbsd openbsd

package lemonwork

import (
	"syscall"
)

type KqueuePoller struct {
	pollBufferSize int
	pollBuffer     []syscall.Kevent_t

	autoClose bool

	onCloseCallback func(fd int)
	onDataCallback  func(fd int)

	kqueuefd int
}

func NewKqueuePoller(PollSize int) (kp *KqueuePoller, err error) {
	kp = new(KqueuePoller)
	kp.pollBufferSize = PollSize
	kp.pollBuffer = make([]syscall.Kevent_t, PollSize)
	kp.kqueuefd, err = syscall.Kqueue()
	if err != nil {
		return nil, err
	}
	kp.autoClose = true
	kp.onCloseCallback = func(fd int) {}
	kp.onDataCallback = func(fd int) {}
	return
}

func (kp *KqueuePoller) PollSize() int {
	return kp.pollBufferSize
}

func (kp *KqueuePoller) Register(fd int) error {
	_, err := syscall.Kevent(kp.kqueuefd, []syscall.Kevent_t{
		{
			Ident:  uint64(fd),
			Filter: syscall.EVFILT_READ,
			Flags:  syscall.EV_ADD,
		},
	}, nil, nil)
	return err
}

func (kp *KqueuePoller) Poll() error {
	kp.pollBuffer = kp.pollBuffer[:kp.pollBufferSize]
	n, err := syscall.Kevent(kp.kqueuefd, nil, kp.pollBuffer, nil)
	if err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		event := kp.pollBuffer[i]
		fd := int(event.Ident)
		flags := event.Flags
		filter := event.Filter

		if flags&syscall.EV_EOF != 0 {
			if kp.autoClose {
				syscall.Close(fd)
			}
			kp.onCloseCallback(fd)
		} else if filter&syscall.EVFILT_READ != 0 {
			kp.onDataCallback(fd)
		}
	}
	return nil
}

func (kp *KqueuePoller) SetOnCloseCallback(f func(fd int)) {
	kp.onCloseCallback = f
}

func (kp *KqueuePoller) SetOnDataCallback(f func(fd int)) {
	kp.onDataCallback = f
}

func (kp *KqueuePoller) SetAutoClose(x bool) {
	kp.autoClose = x
}

func (kp *KqueuePoller) ClosePoller() error {
	return syscall.Close(kp.kqueuefd)
}

func NewPoller(PollSize int) (kp NetPoll, err error) {
	return NewKqueuePoller(PollSize)
}
