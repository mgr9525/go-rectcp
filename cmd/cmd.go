package cmd

import (
	"errors"
	"net"
	"os"
	"os/exec"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

const Version = "0.1.0"

var (
	app      = kingpin.New("rectcp", "A go tcp proxy application.")
	Timeout  = time.Hour
	Timeouts = ""
	Debug    = false

	HostBind   = ""
	HostTarget = ""
)

func Run() {
	regs()
	kingpin.Version(Version)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func regs() {
	cmd := app.Command("run", "run process").Default().
		Action(func(pc *kingpin.ParseContext) error {
			return run()
		})
	cmd.Flag("debug", "debug log show").BoolVar(&Debug)
	cmd.Flag("timeout", "timeout time").Short('t').StringVar(&Timeouts)
	cmd.Arg("bindhost", "bind srouce host").StringVar(&HostBind)
	cmd.Arg("targethost", "target host").StringVar(&HostTarget)

	cmd = app.Command("daemon", "run process background").
		Action(func(pc *kingpin.ParseContext) error {
			return start()
		})
	cmd.Flag("timeout", "timeout time").Short('t').StringVar(&Timeouts)
	cmd.Arg("bindhost", "bind srouce host").StringVar(&HostBind)
	cmd.Arg("targethost", "target host").StringVar(&HostTarget)
}
func run() error {
	if HostBind == "" || HostTarget == "" {
		return errors.New("param err,please check [bindhost,targethost]")
	}
	Timeout = ParseTimed(Timeouts, time.Second*30)

	lsr, err := net.Listen("tcp", HostBind)
	if err != nil {
		return err
	}
	Infof("Start Listen %s", HostBind)
	for {
		runAcp(lsr)
	}
}

func getArgs() []string {
	args := make([]string, 0)
	args = append(args, "run")
	if Timeouts != "" {
		args = append(args, "--timeout")
		args = append(args, Timeouts)
	}
	args = append(args, HostBind)
	args = append(args, HostTarget)
	return args
}
func start() error {
	if HostBind == "" || HostTarget == "" {
		return errors.New("param err,please check [bindhost,targethost]")
	}
	args := getArgs()
	fullpth, err := os.Executable()
	if err != nil {
		return err
	}
	println("start process")
	cmd := exec.Command(fullpth, args...)
	err = cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
