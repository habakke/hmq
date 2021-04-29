package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/habakke/hmq/broker"
)

func init() {
	ConfigureMaxProcs()
}

// Configures the GOMAXPROCS to number CPU cores
func ConfigureMaxProcs() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	config, err := broker.ConfigureConfig(os.Args[1:])
	if err != nil {
		log.Fatal("configure broker config error: ", err)
	}

	b, err := broker.NewBroker(config)
	if err != nil {
		log.Fatal("New Broker error: ", err)
	}
	b.Start()

	s := waitForSignal()
	log.Println("signal received, broker closed.", s)
}

func waitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-signalChan
	signal.Stop(signalChan)
	return s
}
