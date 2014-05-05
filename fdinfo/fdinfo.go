package fdinfo

import (
	"github.com/silenteh/monitoring/jsonutils"
	"github.com/silenteh/monitoring/utils"
	"strings"
)

type FileDescriptor struct {
	Open int
	Max  int
}

func GatherInfo(channel chan string) {
	go func() {
		fdInfo := Stats()
		channel <- jsonutils.ToJsonWithMap(&fdInfo, "fdInfo")
	}()
}

func Stats() FileDescriptor {
	// cat /proc/sys/fs/file-nr

	detectedOs := utils.DetectOS()

	fd := FileDescriptor{}

	switch detectedOs {
	case utils.DARWIN:

		break
	case utils.LINUX:
		result := utils.ExecCommand(true, "cat", "/proc/sys/fs/file-nr")
		itemsArrayTop := utils.ParseMultiValueFormat(result, false, " ")
		if len(itemsArrayTop) > 0 {
			itemsArray := itemsArrayTop[0]
			fd.Open = utils.StringToInt(itemsArray[0], false)
			fd.Max = utils.StringToInt(itemsArray[2], false)
		}

		break
	default:
		break
	}
	return fd

}

func ProcessOpenFiles(pid string) int {

	// ls -la /proc/5155/fd/ | wc -la

	detectedOs := utils.DetectOS()

	openFiles := 0

	switch detectedOs {
	case utils.DARWIN:

		break
	case utils.LINUX:
		path := utils.ConcatenateString("/proc/", pid, "/fd/")
		result := utils.ExecCommand(false, "ls", path)
		itemsArray := strings.Split(result, "/n")
		openFiles = len(itemsArray) - 2 // the 2 is the . and ..

		break
	default:
		break
	}
	return openFiles

}

func ProcessThreads(pid string) int {

	detectedOs := utils.DetectOS()

	threadNumber := 0

	switch detectedOs {
	case utils.DARWIN:

		break
	case utils.LINUX:
		path := utils.ConcatenateString("/proc/", pid, "/task/")
		result := utils.ExecCommand(false, "ls", path)
		itemsArray := strings.Split(result, "/n")
		threadNumber = len(itemsArray) - 2 // the 2 is the . and ..

		break
	default:
		break
	}
	return threadNumber

}
