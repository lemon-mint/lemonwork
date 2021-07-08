package main

import (
	"log"
	"net"
	"syscall"

	"github.com/lemon-mint/lemonwork"
)

const PORT = ":8081"

func main() {
	np, err := lemonwork.NewPoller(16)
	if err != nil {
		log.Fatalln(err)
	}
	defer np.ClosePoller()
	np.SetOnDataCallback(func(fd int) {
		var buffer [1024 * 32]byte
		log.Println("Data received, fd:", fd)
		n, err := syscall.Read(fd, buffer[:])
		if err != nil {
			log.Println("Read Error, fd:", fd, "Error:", err)
			syscall.Close(fd)
		}
		log.Println("Read, fd:", fd, "Size:", n)
		n, err = syscall.Write(fd, buffer[:n])
		if err != nil {
			log.Println("Write Error, fd:", fd, "Error:", err)
			syscall.Close(fd)
		}
		log.Println("Write, fd:", fd, "Size:", n)
	})
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()
	log.Println("Server Started", PORT)

	go func() {
		for {
			np.Poll()
		}
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Accept Error")
			continue
		}
		fd := lemonwork.GetFdFromTCPConn(conn.(*net.TCPConn))
		err = syscall.SetNonblock(fd, true)
		if err != nil {
			log.Println("SetNonblock Error")
			continue
		}
		err = np.Register(fd)
		if err != nil {
			log.Println("Register Error")
			continue
		}
		log.Println("New Connection", conn.RemoteAddr(), "fd:", fd)
	}
}
