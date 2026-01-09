package utils

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

var MAX_URL_LEN = 2097152 // chromiums cap - albeit unlikely
var MAX_HOSTNAME = 253

func CheckDiskspace(targetlist string) (requiredBytes int64, availableBytes uint64, exceeds bool) {
	// Determine largest possible output based on the target list. Warn if overriding users diskspace.
	// so calculate something like:
	// num of targets * len(statuscode + url + host + scheme + success + insecure)
	// TODO: bufio here Probably bottlenecks - i should pass in the pre-existing bufio from a func that already reads the file.
	file, err := os.Open(targetlist)
	if err != nil {
		fmt.Println("Error: ", err)
		return 0, 0, false
	}
	defer file.Close()

	lines := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines += 1
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Fileread error: ", err)
		return 0, 0, false
	}

	// Get available disk space
	var stat unix.Statfs_t
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory: ", err)
		return 0, 0, false
	}
	if err := unix.Statfs(currDir, &stat); err != nil {
		fmt.Println("Error getting disk stats: ", err)
		return 0, 0, false
	}
	availableBytes = stat.Bavail * uint64(stat.Bsize)

	// 2. Do the calculation
	// Worst case: each line could produce MAX_URL_LEN + statuscode + host + scheme + success + insecure
	// statuscode (3-5 bytes) + url (MAX_URL_LEN) + host (MAX_HOSTNAME) + scheme (~8) + success (~6) + insecure (~6) + separators/newlines (~10)
	maxBytesPerLine := MAX_URL_LEN + MAX_HOSTNAME + 100
	requiredBytes = int64(lines) * int64(maxBytesPerLine)

	//3. Check if possibly override current memory.
	exceeds = uint64(requiredBytes) > availableBytes

	//4. Ask user if ok.
	if exceeds {
		fmt.Printf("WARNING: Estimated output size (%d MB) may exceed available disk space (%d MB)\n",
			requiredBytes/(1024*1024), availableBytes/(1024*1024))
	}

	return requiredBytes, availableBytes, exceeds
}
