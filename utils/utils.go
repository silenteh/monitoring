package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func UTCTimeStamp() int32 {
	return int32(time.Now().UTC().Unix())
}

func ConcatenateString(src string, additionals ...string) string {
	var buffer bytes.Buffer
	buffer.WriteString(src)
	for _, element := range additionals {
		buffer.WriteString(element)
	}

	return buffer.String()
}

func ExecCommand(stripNewLines bool, cmd string, args ...string) string {

	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	var finalString string

	if stripNewLines {
		finalString = strings.Replace(string(out), "\n", "", -1)
	} else {
		finalString = string(out)
	}
	return finalString
}

// check if a file exists
func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ReadFile(fileName string) []byte {
	buf := bytes.NewBuffer(nil)

	f, err := os.Open(fileName)
	//defer f.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		io.Copy(buf, f)
		f.Close()
	}
	//s := string(buf.Bytes())
	return buf.Bytes()
}

func WriteFile(fileName string, data []byte, permission os.FileMode) error {
	err := ioutil.WriteFile(fileName, data, permission)
	return err
}

func ParseOutputLabelValueFormat(output string, stripHeader bool, splitLabelBy string) map[string]int {
	// define the map container
	m := make(map[string]int)

	// remove all the multiple spaces or tabs
	removedTabs := strings.Replace(output, "\t", " ", -1)
	removedSpaces := strings.Replace(removedTabs, "  ", " ", -1)

	//spli the new lines
	outArray := strings.Split(removedSpaces, "\n")
	//fmt.Println("-----------------------------------------------")
	for index, element := range outArray {
		// skip the header line !
		if stripHeader && index == 0 {
			continue
		}

		labelValue := strings.Split(element, splitLabelBy)
		// we need to make sure we have only 2 elements in the labelValue !
		if len(labelValue) >= 2 {
			// Trim space from label
			labelNoSPace := strings.TrimSpace(labelValue[0])
			// Trim : from label
			labelNoColumn := strings.Trim(labelNoSPace, ":")

			valueNoSpace := strings.TrimSpace(labelValue[1])
			valueNoColumn := strings.Trim(valueNoSpace, ":")
			valueNoDot := strings.Trim(valueNoColumn, ".")
			valueLower := strings.ToLower(valueNoDot)
			valueNoUnitKb := strings.Replace(valueLower, "kb", "", -1)
			valueNoUnitKbNoSpace := strings.TrimSpace(valueNoUnitKb)

			// convert value to Integer !
			finalValue := StringToInt(valueNoUnitKbNoSpace, false)

			//fmt.Printf("%s - %d \n", labelNoColumn, finalValue)
			m[labelNoColumn] = finalValue
		}
	}
	//fmt.Println("-----------------------------------------------")

	return m
}

func ParseMultiValueFormat(output string, stripHeader bool, splitElementBy string) [][]string {

	// remove all the multiple spaces or tabs
	removedTabs := strings.Replace(output, "\t", " ", -1)

	rp := regexp.MustCompile("\\s{2,}")

	removedSpaces := strings.Replace(removedTabs, "  ", " ", -1)

	//split the new lines
	outArray := strings.Split(removedSpaces, "\n")

	var arraySize int
	if stripHeader {
		arraySize = len(outArray) - 1
	} else {
		arraySize = len(outArray)
	}

	// define the map container
	m := make([][]string, arraySize)
	added := 0
	//fmt.Println("-----------------------------------------------")
	for index, element := range outArray {
		// skip the header line !
		if stripHeader && index == 0 {
			continue
		}

		// sanitize element because of possible multiple spaces
		elementTrimmed := strings.TrimSpace(element)
		elementTrimmedMultipleSPaces := rp.ReplaceAllString(elementTrimmed, " ")
		// start to parse the elements
		lineElements := strings.Split(elementTrimmedMultipleSPaces, splitElementBy)
		size := len(lineElements)
		//fmt.Printf("%s - %d\n", lineElements, size)
		// we need to make sure we have only 2 elements in the labelValue !
		if size > 0 && lineElements[0] != "" {
			//fmt.Printf("%s\n", elementTrimmedMultipleSPaces)
			n := make([]string, size)
			for nIndex, nElement := range lineElements {

				noSpace := strings.TrimSpace(nElement)
				n[nIndex] = noSpace

			}
			m[added] = n
			added++
		}
	}
	//fmt.Println("-----------------------------------------------")
	copied := make([][]string, added)
	copy(copied, m)
	return copied
}

func StringToFloat64(value string, onErrorOne bool) float64 {
	i, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatal(err)
		if onErrorOne {
			return float64(1)
		}
		return float64(0)
	}
	return i
}

func StringToInt(value string, onErrorOne bool) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
		if onErrorOne {
			return 1
		}
		return 0
	}
	return i
}

func StringToUINT64(value string, onErrorOne bool) uint64 {
	i, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		log.Fatal(err)
		if onErrorOne {
			return 1
		}
		return 0
	}
	return i
}

func IntToString(i int64, base int) string {
	s := strconv.FormatInt(i, base)
	return s
}
