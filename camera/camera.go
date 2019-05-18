package camera

import (
	"log"
	"os/exec"
	"strconv"
)

const (
	pictureExecutable = "/usr/bin/raspistill"
)

func TakePicture(width, height int) (bytes []byte, err error) {
	args := []string{
		"--width", strconv.Itoa(width),
		"--height", strconv.Itoa(height),
		"--output", "-",
		"--timeout", "1000",
	}

	log.Println("taking picture")
	if photo, err := exec.Command(pictureExecutable, args...).CombinedOutput(); err != nil {
		return []byte{}, err
	} else {
		return photo, nil
	}
}
