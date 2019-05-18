package camera

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"sync"
)

const (
	pictureExecutable = "/usr/bin/raspistill"
)

var latestPicture []byte
var latestPictureMu sync.RWMutex

func StartTakingPictures(width, height, interval int) {
	args := []string{
		"--width", strconv.Itoa(width),
		"--height", strconv.Itoa(height),
		"--output", "-",
		"--timeout", "0",
		"--timelapse", strconv.Itoa(interval),
	}

	cmd := exec.Command(pictureExecutable, args...)
	stdout, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		log.Fatalf("failed to start raspistill process: %s", err)
	}

	defer func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("failed to end raspistill process: %s", err)
		}
	}()

	readBuffer := make([]byte, 2048)
	currentFile := new(bytes.Buffer)

	for {
		n, err := stdout.Read(readBuffer)
		if err != nil {
			log.Printf("error when reading data: %s", err)
			break
		}

		currentFile.Write(readBuffer)

		if n != 2048 { // todo: fix this to analyze content of the file
			latestPictureMu.Lock()
			latestPicture = currentFile.Bytes()
			latestPictureMu.Unlock()

			currentFile.Reset()

			log.Println("took a picture")
		}
	}
}

func LatestPicture() []byte {
	// todo: what to do if there is no picture yet
	latestPictureMu.RLock()
	defer latestPictureMu.RUnlock()

	return latestPicture
}
