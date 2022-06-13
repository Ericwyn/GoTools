package path

import (
	"os"
	"path/filepath"
	"strings"
)

func GetExecPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "."
	}

	if strings.Contains(dir, "\\Local\\Temp") || strings.Contains(dir, "/tmp") {
		return "."
	}
	return dir
}
