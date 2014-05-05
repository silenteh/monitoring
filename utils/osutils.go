package utils

import (
	"strings"
)

const DARWIN = "darwin"
const LINUX = "linux"

func DetectOS() string {
	out := ExecCommand(true, "uname")
	ostype := strings.ToLower(out)
	return ostype
}

func PageSize() int {
	detectedOs := DetectOS()
	pageSize := 1 // this will get multiplied !!!! DO NOT PUT IT TO ZERO !!!
	switch detectedOs {
	case DARWIN:
		pageSizeString := ExecCommand(true, "pagesize")
		pageSize = StringToInt(pageSizeString, true)
		break
	case LINUX:
		//getconf PAGESIZE
		pageSizeString := ExecCommand(true, "getconf", "PAGESIZE")
		pageSize = StringToInt(pageSizeString, true)
		break
	}
	return pageSize
}
