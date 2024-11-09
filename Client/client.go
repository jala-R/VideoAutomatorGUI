package client

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
)

type MetaData struct {
	ProjectName      string
	ImagesFolder     string
	ReuseAudioFolder string
	ScriptLocation   string
	VoiceID          string
	SentenceGap      float64
	ParaGap          float64
	OutputFolder     string
}

func MakeRequest(metadata MetaData) error {
	//get all Image names
	fmt.Println(metadata)
	images, err := getAllImageNames(metadata.ImagesFolder)
	if err != nil {
		return err
	}
	if metadata.ReuseAudioFolder != "" {
		getAllAudioFiles(metadata.ReuseAudioFolder)
	}

	fmt.Println(images)
	return nil
}

func getAllAudioFiles(folder string) ([]string, error) {
	audioFiles := []string{}

	dirs, err := os.ReadDir(folder)
	if err != nil {
		slog.Error(fmt.Sprintf("Get Image Names: %s", err.Error()))
		return nil, err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			temp := strings.Split(dir.Name(), ".")
			if temp[len(temp)-1] == "mp3" {
				_, err := strconv.ParseInt(temp[0], 10, 32)
				if err != nil {
					err = nil
					continue
				}
				audioFiles = append(audioFiles, dir.Name())
			}
		}
	}

	sort.Slice(audioFiles, func(i, j int) bool {
		f1 := strings.Split(audioFiles[i], ".")
		f2 := strings.Split(audioFiles[j], ".")

		iName, _ := strconv.ParseInt(f1[0], 10, 32)
		jName, _ := strconv.ParseInt(f2[0], 10, 32)

		return iName < jName
	})

	return audioFiles, nil

}

func getAllImageNames(path string) (paths []string, err error) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		slog.Error(fmt.Sprintf("Get Image Names: %s", err.Error()))
		return
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			fileName := dir.Name()

			if isImageFormat(getFileExtension(fileName)) {
				paths = append(paths, fileName)
			}

		}
	}

	return
}

func getFileExtension(fileName string) string {
	temp := strings.Split(fileName, ".")
	return temp[len(temp)-1]
}

func isImageFormat(str string) bool {
	allowedImgExt := []string{
		"jpeg",
		"jpg",
		"png",
		"webp",
	}

	for _, val := range allowedImgExt {
		if val == str {
			return true
		}
	}

	return false
}
