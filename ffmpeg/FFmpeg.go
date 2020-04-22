package ffmpeg

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

type FFmpeg struct {
}

func (ffmpeg *FFmpeg) ExecFFmpegCommand(arguments []string) error  {
	cmd := exec.Command("ffmpeg", arguments...)
	log.Println(cmd.String())
	err := cmd.Run()

	if err != nil {
		fmt.Println("error: ")
		fmt.Println(err)
		return err
	}
	return nil
}

func (ffmpeg *FFmpeg) ExecFFprobeCommand(arguments []string) (string,error)  {
	cmd := exec.Command("ffprobe", arguments...);
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
