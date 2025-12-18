package tools

import "os"

func writeFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}
