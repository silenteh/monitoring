package server

import (
	"bytes"
	"github.com/silenteh/monitoring/config"
	"github.com/silenteh/monitoring/jsonutils"
	"github.com/silenteh/monitoring/queue"
	"github.com/silenteh/monitoring/sysinfo"
	"github.com/silenteh/monitoring/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// HTTP CLIENT
var tr = &http.Transport{}
var client = &http.Client{Transport: tr}

func Register(regChan chan config.ServerConfig) {

	if utils.FileExists("server.json") {
		serverConfig := config.LoadServerConfig()

		serverInfo := sysinfo.ServerInfo
		serverJson := jsonutils.ToJsonWithMap(serverInfo, queue.SYSINFO)
		url := utils.ConcatenateString("http://", config.IP, ":", config.PORT, "/updateserverinfo")
		resp, _ := client.Post(url, "application/json", bytes.NewBufferString(serverJson))
		defer resp.Body.Close()

		regChan <- serverConfig
	} else {

		serverInfo := sysinfo.ServerInfo
		serverJson := jsonutils.ToJsonWithMapForServerRegistration(serverInfo, queue.SYSINFO)

		url := utils.ConcatenateString("http://", config.IP, ":", config.PORT, "/registerserver")
		resp, err := client.Post(url, "application/json", bytes.NewBufferString(serverJson))

		if err != nil {
			// handle error
			log.Panic(err) //.(err)
		}
		defer resp.Body.Close()
		body, errBody := ioutil.ReadAll(resp.Body)
		if errBody != nil {
			log.Panic(err)
		}
		fileErr := utils.WriteFile("server.json", body, os.FileMode.Perm(0640))
		if fileErr != nil {
			log.Fatal(err)
		}
		serverConfig := config.LoadServerConfig()
		regChan <- serverConfig

	}

}
