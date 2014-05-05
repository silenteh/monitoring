package config

import (
	"encoding/json"
	"github.com/silenteh/monitoring/utils"
	"log"
)

const (
	IP   = "127.0.0.1"
	PORT = "8082"
)

var serverConfig ServerConfig
var appConfig AppConfig

type AppConfig struct {
	Key          string `json:"key"`
	Secret       string `json:"secret"`
	LogFile      string `json:"logFile"`
	LogVerbosity string `json:"logVerbosity"`
	PollingTiker int    `json:"pollingTicker"`
}

type ServerConfig struct {
	ServerId string `json:"serverId"`
	CID      string `json:"cid"`
}

func LoadAppConfig() AppConfig {

	if appConfig.Key == "" && appConfig.Secret == "" {
		var localAppConfig AppConfig

		configFile := utils.ReadFile("config.json")

		//jsonParser := json.NewDecoder(configFile)
		if err := json.Unmarshal(configFile, &localAppConfig); err != nil {
			log.Fatal("parsing config file", err.Error())
		}
		appConfig = localAppConfig
		return appConfig
	} else {
		return appConfig
	}

}

func LoadServerConfig() ServerConfig {

	if serverConfig.CID == "" && serverConfig.ServerId == "" {
		var localServerConfig ServerConfig

		configFile := utils.ReadFile("server.json")

		//jsonParser := json.NewDecoder(configFile)
		if err := json.Unmarshal(configFile, &localServerConfig); err != nil {
			log.Fatal("parsing config file", err.Error())
		}
		serverConfig = localServerConfig
		return serverConfig
	} else {
		return serverConfig
	}

}
