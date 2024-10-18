package api

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kkdai/youtube/v2"
)

var DownloadDir = "static/download/"

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

	maxLength := 50
	if len(video.Title) > maxLength {
		filename = video.Title[0:maxLength]
	}

	tempDir, err := os.MkdirTemp(DownloadDir, videoId)
	filePath := filepath.Join(tempDir, filename+".mp4")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("here", err)
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return "", err
	}
	fmt.Println("full", filePath)
	return filePath, nil

}

func FFmpegWrap(arg string) (string, error) {
	filename, err := ExampleClient(arg)
	if err != nil {
		return "", err
	}
	woutMp4 := strings.Split(filename, ".")[0]
	cmd := exec.Command("ffmpeg", "-i", filename, "-q:a", "0", "-map", "a", woutMp4+".mp3")
	cmd.Run()
	fmt.Println(woutMp4)
	return woutMp4 + ".mp3", nil
}
