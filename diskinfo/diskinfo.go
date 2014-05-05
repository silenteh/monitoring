package diskinfo

import (
	"github.com/silenteh/monitoring/jsonutils"
	"github.com/silenteh/monitoring/utils"
	"strings"
)

type DiskLayoutInfo struct {
	Device    string //
	Size      string //
	Used      string //
	Available string
	Usage     string
	Mounted   string
	Ts        int32
}

type DiskStat struct {
	DeviceName     string
	ReadsN         uint64 // number of reads
	ReadsSectorN   uint64 // number of sectors read
	ReadsMS        uint64 // number of MS spent in reads
	WritesN        uint64 // number of Writes
	WritesSectorN  uint64 // number of sectors written
	WritesMS       uint64 // number of MS spent in writes
	IOOperationsN  uint64 // number of I/O operations in progress
	IOOperationsMS uint64 // number of MS spent in I/O operations
	IOBacklog      uint64 // complicated ! https://www.kernel.org/doc/Documentation/iostats.txt
	Ts             int32
}

func GatherInfo(channel chan string) {
	go func() {
		diskLayout := Layout()
		channel <- jsonutils.ToJsonWithMap(&diskLayout, "diskLayout")
	}()

	go func() {
		diskIO := IO()
		channel <- jsonutils.ToJsonWithMap(&diskIO, "diskIO")
	}()
}

func Layout() []DiskLayoutInfo {

	var layout []DiskLayoutInfo

	detectedOs := utils.DetectOS()
	switch detectedOs {
	case utils.DARWIN:
		dfResult := utils.ExecCommand(false, "df", "-h")
		dfResultNoMap := strings.Replace(dfResult, "map ", "", -1)
		itemsArray := utils.ParseMultiValueFormat(dfResultNoMap, true, " ")
		layout = make([]DiskLayoutInfo, len(itemsArray))
		for index, element := range itemsArray {
			size := len(element)
			if size >= 9 {
				l := DiskLayoutInfo{}
				l.Device = element[0]
				l.Size = element[1]
				l.Used = element[2]
				l.Available = element[3]
				l.Usage = element[4]
				l.Mounted = element[8]
				l.Ts = utils.UTCTimeStamp()
				layout[index] = l
			}

		}

		break
	case utils.LINUX:
		dfResult := utils.ExecCommand(false, "df", "-h")
		itemsArray := utils.ParseMultiValueFormat(dfResult, true, " ")
		layout = make([]DiskLayoutInfo, len(itemsArray))
		for index, element := range itemsArray {
			size := len(element)
			if size >= 5 {
				l := DiskLayoutInfo{}
				l.Device = element[0]
				l.Size = element[1]
				l.Used = element[2]
				l.Available = element[3]
				l.Usage = element[4]
				l.Mounted = element[5]
				l.Ts = utils.UTCTimeStamp()
				layout[index] = l
			}
		}
		break
	}

	return layout
}

func IO() []DiskStat {

	var iostats []DiskStat
	var iostatsTemp []DiskStat

	detectedOs := utils.DetectOS()
	switch detectedOs {
	case utils.DARWIN:

		resultBytes := utils.ReadFile("linux_diskstats.txt")
		dfResult := string(resultBytes)

		itemsArray := utils.ParseMultiValueFormat(dfResult, false, " ")
		iostatsTemp = make([]DiskStat, len(itemsArray))
		added := 0
		for _, element := range itemsArray {

			t := element[2]
			if strings.Contains(t, "sd") || strings.Contains(t, "hd") || strings.Contains(t, "xv") || strings.Contains(t, "md") {
				l := DiskStat{}
				l.DeviceName = t
				l.ReadsN = utils.StringToUINT64(element[3], false)
				l.ReadsSectorN = utils.StringToUINT64(element[5], false)
				l.ReadsMS = utils.StringToUINT64(element[6], false)
				l.WritesN = utils.StringToUINT64(element[7], false)
				l.WritesSectorN = utils.StringToUINT64(element[9], false)
				l.WritesMS = utils.StringToUINT64(element[10], false)
				l.IOOperationsN = utils.StringToUINT64(element[11], false)
				l.IOOperationsMS = utils.StringToUINT64(element[12], false)
				l.IOBacklog = utils.StringToUINT64(element[13], false)
				l.Ts = utils.UTCTimeStamp()

				iostatsTemp[added] = l
				added++
			}

		}

		iostats = make([]DiskStat, added)
		copy(iostats, iostatsTemp)

		break
	case utils.LINUX:
		dfResult := utils.ExecCommand(false, "cat", "/proc/diskstats")
		itemsArray := utils.ParseMultiValueFormat(dfResult, false, " ")
		iostatsTemp = make([]DiskStat, len(itemsArray))
		added := 0
		for _, element := range itemsArray {

			t := element[2]
			if strings.Contains(t, "sd") || strings.Contains(t, "hd") || strings.Contains(t, "xv") || strings.Contains(t, "md") {
				l := DiskStat{}
				l.DeviceName = t
				l.ReadsN = utils.StringToUINT64(element[3], false)
				l.ReadsSectorN = utils.StringToUINT64(element[5], false)
				l.ReadsMS = utils.StringToUINT64(element[6], false)
				l.WritesN = utils.StringToUINT64(element[7], false)
				l.WritesSectorN = utils.StringToUINT64(element[9], false)
				l.WritesMS = utils.StringToUINT64(element[10], false)
				l.IOOperationsN = utils.StringToUINT64(element[11], false)
				l.IOOperationsMS = utils.StringToUINT64(element[12], false)
				l.IOBacklog = utils.StringToUINT64(element[13], false)
				l.Ts = utils.UTCTimeStamp()

				iostatsTemp[added] = l
				added++
			}

		}

		iostats = make([]DiskStat, added)
		copy(iostats, iostatsTemp)
		break
	}

	return iostats

}
