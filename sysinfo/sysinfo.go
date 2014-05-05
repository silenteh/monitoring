package sysinfo

import (
	"github.com/silenteh/monitoring/config"
	"github.com/silenteh/monitoring/cpu"
	"github.com/silenteh/monitoring/jsonutils"
	"github.com/silenteh/monitoring/memory"
	"github.com/silenteh/monitoring/utils"
	"log"
	"os"
)

type ServerModel struct {
	Hostname string
	OS       string
	OSModel  string
	Cpu      cpu.CpuInfo
	Ram      float64
	Cid      string
	Ts       int32
}

var server_id = config.LoadServerConfig().ServerId
var ServerInfo = ServerModel{Hostname(), utils.DetectOS(), OsModel(), cpu.Cpus, memory.TotalRam(), server_id, utils.UTCTimeStamp()}

func GatherInfo(channel chan string) {
	go func() {
		channel <- jsonutils.ToJsonWithMap(&ServerInfo, "sysInfo")
	}()
}

// get the system hostname
func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return name
}

// get the OS model version
func OsModel() string {
	detectedOs := utils.DetectOS()

	version := "unknown - lsb_release missing ?"

	switch detectedOs {
	case utils.DARWIN:
		productName := utils.ExecCommand(true, "sw_vers", "-productName")       // osx architecture Ex: x86_64
		productVersion := utils.ExecCommand(true, "sw_vers", "-productVersion") // osx architecture Ex: x86_64
		productBuild := utils.ExecCommand(true, "sw_vers", "-buildVersion")     // osx architecture Ex: x86_64
		version = utils.ConcatenateString(productName, ": ", productVersion, " - build: ", productBuild)
		break
	case utils.LINUX:
		productName := utils.ExecCommand(true, "lsb_release", "-ds")
		codeName := utils.ExecCommand(true, "lsb_release", "-cs")
		version = utils.ConcatenateString(productName, " (", codeName, ")")
		break
	default:
		break
	}
	return version

}
