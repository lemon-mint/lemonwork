package lemonwork

import (
	"net"
	"reflect"
)

func GetFdFromTCPConn(c *net.TCPConn) int {
	return int(reflect.Indirect(
		reflect.Indirect(
			reflect.Indirect(
				reflect.Indirect(
					reflect.Indirect(
						reflect.ValueOf(c),
					).FieldByName("conn"),
				).FieldByName("fd"),
			).FieldByName("pfd"),
		).FieldByName("Sysfd"),
	).Int())
}
