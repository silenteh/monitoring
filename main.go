package main

import (
	"flag"
	"fmt"
	"github.com/silenteh/monitoring/config"
	"github.com/silenteh/monitoring/queue"
	"github.com/silenteh/monitoring/server"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

func main() {

	//read the command line parameters
	flag.Parse()

	// READ THE CONFIG FILES
	appConfig := config.LoadAppConfig()
	if appConfig.Key == "" && appConfig.Secret == "" {
		log.Fatal("Cannot open the config.json !!!")
		return
	}
	// registration channel
	// the app does not continue if the server is not registered with the backend

	registrationChannel := make(chan config.ServerConfig, 1)
	server.Register(registrationChannel)
	serverConfig := <-registrationChannel // here it blocks !

	if serverConfig.ServerId == "" {
		log.Fatal("The server failed to register...")
		return
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// run with multiple cores
	//nproc := cpu.Count()
	//runtime.GOMAXPROCS(nproc)
	runtime.MemProfileRate = 1

	//var responseChannel chan string
	responseChannel := make(chan string)

	queue.Consume(responseChannel)

	//httpserver.Start()

	// this is NOT needed because we store it with the server registration and this does not change unless
	// the server is rebooted
	// on monitoring start the server updates the informations

	//queue.Produce(queue.SYSINFO, responseChannel)
	queue.Produce(queue.LOADINFO, responseChannel)
	queue.Produce(queue.CPULOAD, responseChannel)
	queue.Produce(queue.DISKINFO, responseChannel)
	queue.Produce(queue.NETINFO, responseChannel)

	tickerInterval := appConfig.PollingTiker
	if tickerInterval < 15 {
		tickerInterval = 15
	}
	ticker := time.NewTicker(time.Duration(tickerInterval) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:

				queue.Produce(queue.LOADINFO, responseChannel)
				queue.Produce(queue.CPULOAD, responseChannel)
				queue.Produce(queue.DISKINFO, responseChannel)
				queue.Produce(queue.NETINFO, responseChannel)

				fmt.Printf("%s \n", time.Now())

			case <-quit:
				ticker.Stop()
				fmt.Println("Closed channel")

				for i := 0; i < 10; i++ {
					fmt.Println("OK")
				}
				return
			}
		}
	}()

	if *memprofile != "" {
		time.Sleep(60 * time.Second)
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		defer f.Close()
		//return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	s := <-c
	close(quit)

	close(c)

	fmt.Println("Got signal:", s)

}
