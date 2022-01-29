package main

import (
	"context"
	"errors"
	"io"
	"net"
	"runtime/debug"

	"github.com/mgr9525/go-ruisutil/ruisIo"
)

type Proxy struct {
	runing bool
	conn1  net.Conn
	conn2  net.Conn
	//endfn func()
}

func StartProxy(conn1 net.Conn, conn2 net.Conn) (*Proxy, error) {
	if conn1 == nil || conn2 == nil {
		return nil, errors.New("param err")
	}
	c := &Proxy{
		runing: true,
		conn1:  conn1,
		conn2:  conn2,
	}
	go func() {
		var ln1, ln2 int64
		buf := make([]byte, 1024*100)
		for c.runing {
			err := c.read1(buf, &ln1, &ln2)
			if err != nil {
				if err != io.EOF {
					Debugf("Proxy read1 err:%v", err)
				}
				c.conn1.Close()
				break
				//c.Stop()
			}
		}
		Debugf("Proxy read1 read:%d,write:%d", ln1, ln2)
	}()
	go func() {
		var ln1, ln2 int64
		buf := make([]byte, 1024*100)
		for c.runing {
			err := c.read2(buf, &ln1, &ln2)
			if err != nil {
				if err != io.EOF {
					Debugf("Proxy read2:%v", err)
				}
				c.conn2.Close()
				break
				//c.Stop()
			}
		}
		Debugf("Proxy read2 read:%d,write:%d", ln1, ln2)
	}()
	return c, nil
}
func (c *Proxy) Stop() {
	c.close()
	c.runing = false
}
func (c *Proxy) close() {
	c.conn1.Close()
	c.conn2.Close()
}
func (c *Proxy) read1(buf []byte, ln1, ln2 *int64) error {
	defer func() {
		if err := recover(); err != nil {
			Debugf("Proxy read1 recover:%+v", err)
			Debugf("%s", string(debug.Stack()))
		}
	}()
	n, err := c.conn1.Read(buf)
	if n > 0 {
		*ln1 += int64(n)
		errw := ruisIo.TcpWrite(context.Background(), c.conn2, buf[:n])
		// nw, errw := c.conn2.Write(buf[:n])
		*ln2 += int64(n)
		if errw != nil {
			return errw
		}
		/* if n != nw {
			return errors.New("write len err")
		} */
	}
	if err != nil {
		return err
	}
	return nil
}
func (c *Proxy) read2(buf []byte, ln1, ln2 *int64) error {
	defer func() {
		if err := recover(); err != nil {
			Debugf("Proxy read1 recover:%+v", err)
			Debugf("%s", string(debug.Stack()))
		}
	}()
	n, err := c.conn2.Read(buf)
	if n > 0 {
		*ln1 += int64(n)
		errw := ruisIo.TcpWrite(context.Background(), c.conn1, buf[:n])
		// nw, errw := c.conn1.Write(buf[:n])
		*ln2 += int64(n)
		if errw != nil {
			return errw
		}
		/* if n != nw {
			return errors.New("write len err")
		} */
	}
	if err != nil {
		return err
	}
	return nil
}
