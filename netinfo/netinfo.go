package netinfo

import (
	"github.com/silenteh/monitoring/jsonutils"
	"github.com/silenteh/monitoring/utils"
	"strings"
)

type NetStat struct {
	InterfaceName string
	RXBytes       uint64 //
	RXPackets     uint64 //
	RXErrors      uint64 //
	RXDrops       uint64 //
	RXFifo        uint64 //
	RXFrame       uint64 //
	RXCompressed  uint64 //
	RXMulticast   uint64 //
	TXBytes       uint64 //
	TXPackets     uint64 //
	TXErrors      uint64 //
	TXDrops       uint64 //
	TXFifo        uint64 //
	TXCollisions  uint64 //
	TXCarrier     uint64 //
	TXCompressed  uint64 //
	Ts            int32
}

func GatherInfo(channel chan string) {
	go func() {
		netStats := Stats()
		channel <- jsonutils.ToJsonWithMap(&netStats, "netInfo")
	}()
}

func Stats() []NetStat {
	// cat /proc/net/dev

	var allStats []NetStat
	var allStatsTemp []NetStat

	detectedOs := utils.DetectOS()

	switch detectedOs {
	case utils.DARWIN:
		//output := utils.ExecCommand(false, "cat", "/proc/net/dev")
		resultBytes := utils.ReadFile("netinfo.txt")
		dfResult := string(resultBytes)
		//fmt.Printf("%s\n", dfResult)

		itemsArray := utils.ParseMultiValueFormat(dfResult, true, " ")
		allStatsTemp = make([]NetStat, len(itemsArray))
		added := 0
		for index, element := range itemsArray {
			if index == 0 { // /proc/net/dev has a double header !
				continue
			}
			size := len(element)
			if size >= 17 {
				l := NetStat{}
				l.InterfaceName = strings.Replace(element[0], ":", "", -1)
				l.RXBytes = utils.StringToUINT64(element[1], false)
				l.RXPackets = utils.StringToUINT64(element[2], false)
				l.RXErrors = utils.StringToUINT64(element[3], false)
				l.RXDrops = utils.StringToUINT64(element[4], false)
				l.RXFifo = utils.StringToUINT64(element[5], false)
				l.RXFrame = utils.StringToUINT64(element[6], false)
				l.RXCompressed = utils.StringToUINT64(element[7], false)
				l.RXMulticast = utils.StringToUINT64(element[8], false)

				l.TXBytes = utils.StringToUINT64(element[9], false)
				l.TXPackets = utils.StringToUINT64(element[10], false)
				l.TXErrors = utils.StringToUINT64(element[11], false)
				l.TXDrops = utils.StringToUINT64(element[12], false)
				l.TXFifo = utils.StringToUINT64(element[13], false)
				l.TXCollisions = utils.StringToUINT64(element[14], false)
				l.TXCarrier = utils.StringToUINT64(element[15], false)
				l.TXCompressed = utils.StringToUINT64(element[16], false)
				l.Ts = utils.UTCTimeStamp()
				allStatsTemp[added] = l
				added++
			}
		}

		allStats = make([]NetStat, added)
		copy(allStats, allStatsTemp)
		break
	case utils.LINUX:
		output := utils.ExecCommand(false, "cat", "/proc/net/dev")
		dfResult := string(output)
		itemsArray := utils.ParseMultiValueFormat(dfResult, true, " ")
		allStatsTemp = make([]NetStat, len(itemsArray))
		added := 0
		for index, element := range itemsArray {
			if index == 0 { // /proc/net/dev has a double header !
				continue
			}
			size := len(element)
			if size >= 17 {
				l := NetStat{}
				l.InterfaceName = strings.Replace(element[0], ":", "", -1)
				l.RXBytes = utils.StringToUINT64(element[1], false)
				l.RXPackets = utils.StringToUINT64(element[2], false)
				l.RXErrors = utils.StringToUINT64(element[3], false)
				l.RXDrops = utils.StringToUINT64(element[4], false)
				l.RXFifo = utils.StringToUINT64(element[5], false)
				l.RXFrame = utils.StringToUINT64(element[6], false)
				l.RXCompressed = utils.StringToUINT64(element[7], false)
				l.RXMulticast = utils.StringToUINT64(element[8], false)

				l.TXBytes = utils.StringToUINT64(element[9], false)
				l.TXPackets = utils.StringToUINT64(element[10], false)
				l.TXErrors = utils.StringToUINT64(element[11], false)
				l.TXDrops = utils.StringToUINT64(element[12], false)
				l.TXFifo = utils.StringToUINT64(element[13], false)
				l.TXCollisions = utils.StringToUINT64(element[14], false)
				l.TXCarrier = utils.StringToUINT64(element[15], false)
				l.TXCompressed = utils.StringToUINT64(element[16], false)
				l.Ts = utils.UTCTimeStamp()
				allStatsTemp[added] = l
				added++
			}
		}

		allStats = make([]NetStat, added)
		copy(allStats, allStatsTemp)
		break
	}
	return allStats
}
