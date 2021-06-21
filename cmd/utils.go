package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var regTms = regexp.MustCompile(`^(\d+)(\w{1,2})$`)

func ParseTimed(tms string, defs ...time.Duration) time.Duration {
	bs := time.Hour
	if len(defs) > 0 && defs[0] > 0 {
		bs = defs[0]
	}
	if !regTms.MatchString(tms) {
		return bs
	}
	s := regTms.FindAllStringSubmatch(tms, -1)[0]
	switch s[2] {
	case "h":
		bs = time.Hour
	case "m":
		bs = time.Minute
	case "s":
		bs = time.Second
	case "ms":
		bs = time.Millisecond
	}
	n, _ := strconv.ParseInt(s[1], 10, 64)
	if n <= 0 {
		return bs
	}
	return bs * time.Duration(n)
}

func Infof(s string, args ...interface{}) {
	tms := time.Now().Format("2006-01-02 15:04:05")
	if len(args) > 0 {
		fmt.Println(tms + " [info] " + fmt.Sprintf(s, args...))
	} else {
		fmt.Println(tms + " [info] " + s)
	}
}
func Errorf(s string, args ...interface{}) {
	tms := time.Now().Format("2006-01-02 15:04:05")
	if len(args) > 0 {
		println(tms + " [err] " + fmt.Sprintf(s, args...))
	} else {
		println(tms + " [err] " + s)
	}
}
func Debugf(s string, args ...interface{}) {
	if !Debug {
		return
	}
	tms := time.Now().Format("2006-01-02 15:04:05")
	if len(args) > 0 {
		println(tms + " [debug] " + fmt.Sprintf(s, args...))
	} else {
		println(tms + " [debug] " + s)
	}
}
