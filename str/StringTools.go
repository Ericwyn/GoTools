package str

import (
	"regexp"
	"strconv"
	"strings"
)

func ReplaceAll(str string, old string, new string) string {
	if strings.Index(str, old) < 0 {
		return str
	} else {
		return ReplaceAll(strings.ReplaceAll(str, old, new), old, new)
	}
}

func Continues(str string, continueStr string) bool {
	return strings.Contains(str, continueStr)
}

func StartWith(str string, startStr string) bool {
	return strings.Index(str, startStr) == 0
}

func GetAllStrByRegexp(str string, arrRegexp *regexp.Regexp) []string {
	return arrRegexp.FindAllString(str, -1)
}

func ConverToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

//func ConverToDouble()  {
//
//}

func Split(str string, sep string) []string {
	return strings.Split(str, sep)
}

func ToLowCase(str string) string {
	return strings.ToLower(str)
}

func ToUpper(str string) string {
	return strings.ToUpper(str)
}
