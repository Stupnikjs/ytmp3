package api

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kkdai/youtube/v2"
)

var DowloadDir = "/static/download"

func ExampleClient(videoId string) (string, error) {

	client := youtube.Client{}
	video, err := client.GetVideo(videoId)
	if err != nil {
		return "", err
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return "", err
	}
	defer stream.Close()
	filename := ""

	filename = video.Title

	// ADD THE PATH TO STATIC FILES

	maxLength := 20
	if len(video.Title) > maxLength {
		filename = video.Title[0:maxLength]
	}

	file, err := os.Create(filename + ".mp4")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return "", err
	}
	return file.Name(), nil

}

func FFmpegWrap(arg string) (error, string) {
	filename, err := ExampleClient(arg)
	if err != nil {
		return err, ""
	}
	woutMp4 := strings.Split(filename, ".")[0]

	cmd := exec.Command("ffmpeg", "-i", filename, "-q:a", "0", "-map", "a", woutMp4+".mp3")
	cmd.Run()

	// not working
	cmd = exec.Command("rm", "-r", filename)
	cmd.Run()
	return nil, woutMp4 + ".mp3"
}
