package cpu

import (
	"bytes"
	"github.com/silenteh/monitoring/jsonutils"
	"github.com/silenteh/monitoring/utils"
	"log"
	"os/exec"
	"strings"
)

type CpuInfo struct {
	Number       int
	Architecture string
	Model        string
	Ts           int32
}

type CpuUsage struct {
	Core           string
	User           float64 // amount of CPU used by the users
	Nice           float64 // amount of CPU used by the users with low priority
	System         float64 // amount of CPU used by the OS
	Idle           float64 // amount of CPU IDLE
	IOWait         float64 // amount of CPU used by the I/O
	Steal          float64 // amount of CPU stolen
	VirtualCPU     float64 // amount of CPU used by the Virtual Machines
	NiceVirtualCPU float64 // amount of CPU used by the Virtual Machines with low priority
	Overall        bool    // is the overall of ALL cores ?
	Ts             int32
}

var Cpus = CpuInfo{Count(), Arch(), Description(), utils.UTCTimeStamp()}

func GatherInfo(channel chan string) {
	go func() {
		usage := CpuLoad()
		channel <- jsonutils.ToJsonWithMap(&usage, "cpuLoad")
	}()
}

// get CPU count
func Count() int {
	detectedOs := utils.DetectOS()

	cpuNumber := 1

	switch detectedOs {
	case utils.DARWIN:
		cpuNumberString := utils.ExecCommand(true, "/usr/sbin/sysctl", "-n", "hw.ncpu") // osx number of CPU Ex: x86_64
		cpuNumber = utils.StringToInt(cpuNumberString, true)
		break
	case utils.LINUX:
		cpuNumberString := utils.ExecCommand(true, "nproc") // linux
		cpuNumber = utils.StringToInt(cpuNumberString, true)
		break
	}

	return cpuNumber
}

// get CPU architecture
func Arch() string {
	detectedOs := utils.DetectOS()

	architecture := "unknown"

	switch detectedOs {
	case utils.DARWIN:
		architecture = utils.ExecCommand(true, "/usr/sbin/sysctl", "-n", "hw.machine") // osx architecture Ex: x86_64
		break
	case utils.LINUX:
		architecture = utils.ExecCommand(true, "uname", "-p") // linux architecture
		break
	default:
		//println("None detected !!!!")
		break
	}

	return architecture
}

func Description() string {

	detectedOs := utils.DetectOS()

	cpuDescription := "unknown"

	switch detectedOs {
	case utils.DARWIN:
		cmdResult := utils.ExecCommand(true, "/usr/sbin/sysctl", "machdep.cpu.brand_string")
		cpuDescription = strings.Replace(cmdResult, "machdep.cpu.brand_string: ", "", -1) // osx descr Ex: x86_64
		break
	case utils.LINUX:

		data, errData := exec.Command("cat", "/proc/cpuinfo").CombinedOutput()
		if errData != nil {
			log.Fatal(errData)
		}
		cpuInfo := string(data)

		cmd := exec.Command("grep", "-m1", "name")
		cmd.Stdin = strings.NewReader(cpuInfo)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		rawResult := out.String()
		replaceLabel := strings.Replace(rawResult, "model name\t: ", "", -1)
		replaceNewLine := strings.Replace(replaceLabel, "\n", "", -1)
		finalInfo := strings.Replace(replaceNewLine, "  ", "", -1)

		cpuDescription = finalInfo

		break
	}

	return cpuDescription

}

func CpuLoad() []CpuUsage {
	detectedOs := utils.DetectOS()
	numCpu := Count() + 1

	cpuInfo := make([]CpuUsage, numCpu, numCpu)

	switch detectedOs {
	case utils.DARWIN:
		resultBytes := utils.ReadFile("cpuinfo.txt")
		dfResult := string(resultBytes)
		//fmt.Printf("%s\n", dfResult)

		infoArray := utils.ParseMultiValueFormat(dfResult, false, " ")
		allStatsTemp := make([]CpuUsage, len(infoArray))
		added := 0
		for _, element := range infoArray {
			cpuLabel := element[0]
			size := len(element)
			if strings.Contains(cpuLabel, "cpu") {
				cpu := CpuUsage{}
				if cpuLabel == "cpu" {
					cpu.Overall = true
				}

				cpu.Core = cpuLabel
				cpu.User = utils.StringToFloat64(element[1], false)
				cpu.Nice = utils.StringToFloat64(element[2], false)
				cpu.System = utils.StringToFloat64(element[3], false)
				cpu.Idle = utils.StringToFloat64(element[4], false)
				cpu.IOWait = utils.StringToFloat64(element[5], false)
				if size > 5 {
					cpu.Steal = utils.StringToFloat64(element[6], false)
				}
				if size > 6 {
					cpu.VirtualCPU = utils.StringToFloat64(element[7], false)
				}
				if size > 7 {
					cpu.NiceVirtualCPU = utils.StringToFloat64(element[8], false)
				}
				cpu.Ts = utils.UTCTimeStamp()
				allStatsTemp[added] = cpu
				added++
			}
		}

		copy(cpuInfo, allStatsTemp)

		break
	case utils.LINUX:
		output := utils.ExecCommand(false, "cat", "/proc/stat")
		dfResult := string(output)
		//fmt.Printf("%s\n", dfResult)

		infoArray := utils.ParseMultiValueFormat(dfResult, false, " ")
		allStatsTemp := make([]CpuUsage, len(infoArray))
		added := 0
		for _, element := range infoArray {
			cpuLabel := element[0]
			size := len(element)
			if strings.Contains(cpuLabel, "cpu") {
				cpu := CpuUsage{}
				if cpuLabel == "cpu" {
					cpu.Overall = true
				}

				cpu.Core = cpuLabel
				cpu.User = utils.StringToFloat64(element[1], false)
				cpu.Nice = utils.StringToFloat64(element[2], false)
				cpu.System = utils.StringToFloat64(element[3], false)
				cpu.Idle = utils.StringToFloat64(element[4], false)
				cpu.IOWait = utils.StringToFloat64(element[5], false)
				if size > 5 {
					cpu.Steal = utils.StringToFloat64(element[6], false)
				}
				if size > 6 {
					cpu.VirtualCPU = utils.StringToFloat64(element[7], false)
				}
				if size > 7 {
					cpu.NiceVirtualCPU = utils.StringToFloat64(element[8], false)
				}
				cpu.Ts = utils.UTCTimeStamp()
				allStatsTemp[added] = cpu
				added++
			}
		}

		copy(cpuInfo, allStatsTemp)

		break
	}

	return cpuInfo
}
