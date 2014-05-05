package memory

import (
	"github.com/silenteh/monitoring/jsonutils"
	"github.com/silenteh/monitoring/utils"
	"strings"
)

type MemoryUsage struct {
	Total         float64
	Available     float64 // amount of RAM available
	App           float64 // amount of RAM used by Apps
	OS            float64 // amount of RAM used by ther OS ! cannot be borrowed
	Used          float64
	TotalSwap     float64
	AvailableSwap float64
	UsedSwap      float64
	Ts            int32
}

func GatherInfo(channel chan string) {
	go func() {
		mem := UsedRam()
		channel <- jsonutils.ToJsonWithMap(&mem, "memInfo")
	}()
}

func TotalRam() float64 {
	detectedOs := utils.DetectOS()
	availableRam := float64(0)
	switch detectedOs {
	case utils.DARWIN:
		availableRamString := utils.ExecCommand(true, "/usr/sbin/sysctl", "-n", "hw.memsize")
		availableRamStringLower := strings.ToLower(availableRamString)
		availableRam = convertStringBytesToGB(availableRamStringLower)
		break
	case utils.LINUX:
		output := utils.ExecCommand(false, "cat", "/proc/meminfo") // linux
		m := utils.ParseOutputLabelValueFormat(output, false, ":")
		availableRam = convertKiloBytesToGB(m["MemTotal"])
		break
	}
	return availableRam
}

func UsedRam() MemoryUsage {

	memUsage := MemoryUsage{}
	memUsage.Total = TotalRam()

	detectedOs := utils.DetectOS()
	switch detectedOs {
	case utils.DARWIN:
		output := utils.ExecCommand(false, "vm_stat")
		memUsage = realTimeMemoryFromCommandOutput(output, detectedOs)
		memUsage.Ts = utils.UTCTimeStamp()
		break
	case utils.LINUX:
		output := utils.ExecCommand(false, "cat", "/proc/meminfo")
		memUsage = realTimeMemoryFromCommandOutput(output, detectedOs)
		memUsage.Ts = utils.UTCTimeStamp()
		break
	}
	return memUsage
}

func realTimeMemoryFromCommandOutput(output string, detectedOs string) MemoryUsage {
	memUsage := MemoryUsage{}

	switch detectedOs {
	case utils.DARWIN:
		pageSize := utils.PageSize()

		m := utils.ParseOutputLabelValueFormat(output, true, ":")
		pagesOccupiedByCompressor := (float64(m["Pages occupied by compressor"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		//pagesStoredInCompressor := (float64(m["Pages stored in compressor"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		fileCache := (float64(m["File-backed pages"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		pagesFree := (float64(m["Pages free"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		//pagesActive := (float64(m["Pages active"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		pagesInactive := (float64(m["Pages inactive"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		pagesWired := (float64(m["Pages wired down"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		pagesSpeculative := (float64(m["Pages speculative"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		//anonymousPages := (float64(m["Anonymous pages"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		//pagesPurgeable := (float64(m["Pages purgeable"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		pagesReactivated := (float64(m["Pages reactivated"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))

		//fmt.Printf("File Cache: %s GB\n", strconv.FormatFloat(fileCache, 'f', 2, 64))
		//fmt.Printf("Wired Memory: %s GB\n", strconv.FormatFloat(pagesWired, 'f', 2, 64))
		//fmt.Printf("Compressed: %s GB\n", strconv.FormatFloat(pagesOccupiedByCompressor, 'f', 3, 64))

		memUsage.Available = pagesFree
		memUsage.App = pagesInactive + pagesSpeculative + pagesReactivated //- pagesInactive - pagesPurgeable
		memUsage.OS = pagesWired
		//memUsage.Used = pagesActive + pagesInactive + pagesSpeculative + pagesWired + pagesOccupiedByCompressor
		memUsage.Used = memUsage.App + fileCache + pagesWired + pagesOccupiedByCompressor

		break
	case utils.LINUX:
		m := utils.ParseOutputLabelValueFormat(output, false, ":")
		totalMem := convertKiloBytesToGB(m["MemTotal"])
		freeMem := convertKiloBytesToGB(m["MemFree"])

		buffersMem := convertKiloBytesToGB(m["Buffers"])
		cachedMem := convertKiloBytesToGB(m["Cached"])

		totalSwap := convertKiloBytesToGB(m["SwapTotal"])
		freeSwap := convertKiloBytesToGB(m["SwapFree"])

		memUsage.OS = buffersMem + cachedMem

		memUsage.Available = freeMem + memUsage.OS
		memUsage.Used = totalMem - freeMem

		memUsage.App = memUsage.Used - memUsage.OS

		memUsage.TotalSwap = totalSwap
		memUsage.AvailableSwap = freeSwap
		memUsage.UsedSwap = totalSwap - freeSwap

		break
	}

	return memUsage
}

func convertBytesToGB(value int) float64 {
	return (float64(value) / float64(1024.00) / float64(1024.00) / float64(1024.00))
}

func convertStringBytesToGB(value string) float64 {
	i := utils.StringToFloat64(value, false)
	floatValue := (i / float64(1024.00) / float64(1024.00) / float64(1024.00))
	return floatValue
}

func convertKiloBytesToGB(value int) float64 {
	return (float64(value) / float64(1024.00) / float64(1024.00))
}
