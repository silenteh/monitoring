package queue

import (
	"bytes"
	"github.com/silenteh/monitoring/config"
	"github.com/silenteh/monitoring/cpu"
	"github.com/silenteh/monitoring/diskinfo"
	"github.com/silenteh/monitoring/fdinfo"
	"github.com/silenteh/monitoring/loadinfo"
	"github.com/silenteh/monitoring/memory"
	"github.com/silenteh/monitoring/netinfo"
	"github.com/silenteh/monitoring/sysinfo"
	"github.com/silenteh/monitoring/utils"
	"log"
	"net/http"
)

const SYSINFO = "sysInfo"
const NETINFO = "netInfo"
const DISKINFO = "diskInfo"
const LOADINFO = "loadInfo"
const MEMINFO = "memInfo"
const CPULOAD = "cpuLoad"
const FD_INFO = "fdInfo"

type Request struct {
	RequestType    string
	RequestChannel chan string
}

type Response struct {
	ResponseData    string
	ResponseType    string
	ResponseChannel chan string
}

var pushChannel = make(chan string, 20)

// HTTP CLIENT
var tr = &http.Transport{}
var client = &http.Client{Transport: tr}

// func Start(queueChannel chan string) {

// }

// func Stop(queueChannel chan string) {
// 	queueChannel.Stop()
// }

func init() {
	push()
}

func Produce(request string, responseChannel chan string) {

	//typeOfRequest := request.RequestType

	switch request {

	case SYSINFO:
		sysinfo.GatherInfo(responseChannel)
		break

	case NETINFO:
		netinfo.GatherInfo(responseChannel)
		break

	case DISKINFO:
		diskinfo.GatherInfo(responseChannel)
		break

	case LOADINFO:
		loadinfo.GatherInfo(responseChannel)
		break

	case MEMINFO:
		memory.GatherInfo(responseChannel)
		break

	case CPULOAD:
		cpu.GatherInfo(responseChannel)
		break

	case FD_INFO:
		fdinfo.GatherInfo(responseChannel)
		break

	}

}

func Consume(responseChannel chan string) {
	go func() {
		for {
			response := <-responseChannel
			pushChannel <- response
			//fmt.Println(response)
			//fmt.Println("------------")
		}
	}()
}

func push() {
	go func() {
		for {

			toPush := <-pushChannel

			//r, err := client.Post("http://127.0.0.1:8082/push", "application/json", bytes.NewBufferString(toPush))
			url := utils.ConcatenateString("http://", config.IP, ":", config.PORT, "/v1/collector")
			resp, err := client.Post(url, "application/json", bytes.NewBufferString(toPush))

			if err != nil {
				// handle error
				log.Fatal(err) //.(err)
			}
			defer resp.Body.Close()
			//body, err := ioutil.ReadAll(resp.Body)

		}
	}()
}
