package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/toolkits/pkg/file"
	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/runner"

	"github.com/leezer3379/sms-sender/config"
	"github.com/leezer3379/sms-sender/cron"
	"github.com/leezer3379/sms-sender/redisc"
)

var (
	vers *bool
	help *bool
	conf *string
	test *string
)

func init() {
	vers = flag.Bool("v", false, "display the version.")
	help = flag.Bool("h", false, "print this help.")
	conf = flag.String("f", "", "specify configuration file.")
	test = flag.String("t", "", "test configuration.")
	flag.Parse()

	if *vers {
		fmt.Println("version:", config.Version)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	runner.Init()
	fmt.Println("runner.cwd:", runner.Cwd)
	fmt.Println("runner.hostname:", runner.Hostname)
}

func main() {
	aconf()
	pconf()

	if *test != "" {
		config.Test(strings.Split(*test, ","))
		os.Exit(0)
	}

	config.InitLogger()
	redisc.InitRedis()

	go cron.SendSms()

	ending()
}

func ending() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-c:
		fmt.Printf("stop signal caught, stopping... pid=%d\n", os.Getpid())
	}

	logger.Close()
	redisc.CloseRedis()
	fmt.Println("sender stopped successfully")
}

// auto detect configuration file
func aconf() {
	if *conf != "" && file.IsExist(*conf) {
		return
	}

	*conf = path.Join(runner.Cwd, "etc", "sms-sender.local.yml")
	if file.IsExist(*conf) {
		return
	}

	*conf = path.Join(runner.Cwd, "etc", "sms-sender.yml")
	if file.IsExist(*conf) {
		return
	}

	fmt.Println("no configuration file for sender")
	os.Exit(1)
}

// parse configuration file
func pconf() {
	if err := config.ParseConfig(*conf); err != nil {
		fmt.Println("cannot parse configuration file:", err)
		os.Exit(1)
	} else {
		fmt.Println("parse configuration file:", *conf)
	}
}