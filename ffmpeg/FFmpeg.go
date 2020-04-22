package ffmpeg

import (
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
