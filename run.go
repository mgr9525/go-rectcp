package main

import (
	"net"
	"runtime/debug"
)

func runAcp(lsr net.Listener) {
	defer func() {
		if err := recover(); err != nil {
			Debugf("runAcp recover:%+v", err)
			Debugf("%s", string(debug.Stack()))
		}
	}()
	conn, err := lsr.Accept()
	if err != nil {
		Errorf("runAcp AcceptTCP err:%+v", err)
		return
	}
	go handleConn(conn)
}
func handleConn(connSrc net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			connSrc.Close()
			Debugf("handleConn recover:%+v", err)
			Debugf("%s", string(debug.Stack()))
		}
	}()
	connTrg, err := net.DialTimeout("tcp", HostTarget, Timeout)
	if err != nil {
		connSrc.Close()
		Errorf("connect %s failed:%v", HostTarget, err)
		return
	}
	Infof("Start Proxy %s->%s(%s)", connSrc.LocalAddr().String(), connTrg.RemoteAddr().String(), HostTarget)
	StartProxy(connSrc, connTrg)
}
