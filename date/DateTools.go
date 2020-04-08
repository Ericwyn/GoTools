package date

import (
	"github.com/Ericwyn/GoTools/str"
	"time"
)

var dateFormatMap = map[string]string{}

func Format(time time.Time, formatStr string) string {
	formatStr = getFormatString(formatStr)
	return time.Format(getFormatString(formatStr))
}

// 传入 yyyyMMdd 之类的
func getFormatString(formatStr string) string {
	// 原始 format string
	formatStrOld := formatStr
	// 构建一个 formatStr
	if formatStr := dateFormatMap[formatStrOld]; formatStr == "" {
		formatStr = str.ReplaceAll(formatStrOld, "yyyy", "2006")
		formatStr = str.ReplaceAll(formatStr, "yy", "06")
		formatStr = str.ReplaceAll(formatStr, "MM", "01")
		formatStr = str.ReplaceAll(formatStr, "dd", "02")
		formatStr = str.ReplaceAll(formatStr, "HH", "15")
		formatStr = str.ReplaceAll(formatStr, "hh", "3")
		formatStr = str.ReplaceAll(formatStr, "mm", "04")
		formatStr = str.ReplaceAll(formatStr, "ss", "05")
		dateFormatMap[formatStrOld] = formatStr
	}
	return dateFormatMap[formatStrOld]
}
