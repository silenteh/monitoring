package jsonutils

import (
	"encoding/json"
	"github.com/silenteh/monitoring/config"
	"github.com/silenteh/monitoring/utils"
	"log"
)

func ToJson(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(b)
}

func ToJsonWithMap(data interface{}, key string) string {

	serverConfig := config.LoadServerConfig()
	appConfig := config.LoadAppConfig()

	jsonMap := make(map[string]interface{})
	jsonMap["serverId"] = serverConfig.ServerId
	jsonMap["cid"] = serverConfig.CID

	jsonMap["key"] = appConfig.Key
	jsonMap["secret"] = appConfig.Secret

	jsonMap["type"] = key
	jsonMap["data"] = data

	jsonMap["ts"] = utils.UTCTimeStamp() //int32(time.Now().Unix())

	b, err := json.Marshal(jsonMap)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(b)
}

func ToJsonWithMapForServerRegistration(data interface{}, key string) string {

	appConfig := config.LoadAppConfig()

	jsonMap := make(map[string]interface{})

	jsonMap["key"] = appConfig.Key
	jsonMap["secret"] = appConfig.Secret

	jsonMap["type"] = key
	jsonMap["data"] = data

	b, err := json.Marshal(jsonMap)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(b)
}
