package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/kkdai/youtube/v2"
)

func main() {
	// YouTube video URL
	videoURL := os.Args[1]

	// Initialize YouTube client
	client := youtube.Client{}

	// Fetch video details
	video, err := client.GetVideo(videoURL)
	if err != nil {
		log.Fatalf("Error fetching video info: %v", err)
	}

	// Look for a format that contains both video and audio
	var combinedFormat *youtube.Format
	for _, format := range video.Formats {
		if format.AudioChannels > 0 { // This ensures it has both video and audio
			combinedFormat = &format
			break // Exit loop as soon as a suitable format is found
		}
	}

	// Check if a combined video+audio format was found
	if combinedFormat == nil {
		log.Fatalf("Error: Could not find a format with both video and audio.")
	}

	// Download the combined video and audio stream
	downloadWithProgress(client, video, combinedFormat, "output_combined.mp4")

	fmt.Println("Download complete with video and audio!")
}

// downloadWithProgress downloads a stream with a progress bar
func downloadWithProgress(client youtube.Client, video *youtube.Video, format *youtube.Format, outputFileName string) {
	// Open output file
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer outputFile.Close()

	// Get the video/audio stream
	resp, length, err := client.GetStream(video, format)
	if err != nil {
		log.Fatalf("Error fetching stream: %v", err)
	}
	defer resp.Close()

	// Initialize progress bar
	bar := pb.Full.Start64(length)
	barReader := bar.NewProxyReader(resp)

	// Write stream to file with progress
	_, err = io.Copy(outputFile, barReader)
	if err != nil {
		log.Fatalf("Error saving file: %v", err)
	}

	bar.Finish()
	fmt.Printf("Downloaded %s\n", outputFileName)
}
