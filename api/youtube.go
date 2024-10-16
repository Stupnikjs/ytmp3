package api

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func ExampleClient(videoId string) string {

	client := youtube.Client{}
	video, err := client.GetVideo(videoId)
	if err != nil {
		panic(err)
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}
	defer stream.Close()
	filename := ""

	filename = video.Title
	maxLength := 20
	if len(video.Title) > maxLength {
		filename = video.Title[0:maxLength]
	}

	file, err := os.Create(filename + ".mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}
	return file.Name()

}

func FFmpegWrap(arg string) string {
	filename := ExampleClient(arg)
	defer os.Remove(filename)
	woutMp4 := strings.Split(filename, ".")[0]

	cmd := exec.Command("ffmpeg", "-i", filename, "-q:a", "0", "-map", "a", woutMp4+".mp3")
	cmd.Run()
	return woutMp4 + ".mp3"
}
