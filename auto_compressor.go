package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// check number of args
func checkArgs() {
	if len(os.Args) != 3 {
		fmt.Println("usage: go run auto_compressor.go <input_file.mp4> <desired_output_size_MB>")
		os.Exit(1)
	}
}

// run ffmpeg to get video duration and bitrates
func getVideoInfo(inputFile string) (float64, int, int, error) {
	cmd := exec.Command("ffmpeg", "-i", inputFile)

	// ffmpeg sends the info to stderr, not stdout
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 0, 0, 0, err
	}

	// start the ffmpeg command
	if err := cmd.Start(); err != nil {
		return 0, 0, 0, err
	}

	// read the stderr output where ffmpeg writes the info
	outputBytes, err := io.ReadAll(stderr)
	if err != nil {
		return 0, 0, 0, err
	}

	// wait for the command to finish
	if err := cmd.Wait(); err != nil {
		// ffmpeg returns a non-zero exit code even when just printing info
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() != 0 {
			// ignore the error since it's expected
		} else {
			return 0, 0, 0, err
		}
	}

	outputStr := string(outputBytes)

	// extract duration and bitrates from the output
	duration := parseDuration(outputStr)
	if duration == 0 {
		return 0, 0, 0, fmt.Errorf("could not parse duration")
	}

	videoBitrate := parseVideoBitrate(outputStr)
	if videoBitrate == 0 {
		return 0, 0, 0, fmt.Errorf("could not parse video bitrate")
	}

	audioBitrate := parseAudioBitrate(outputStr)
	if audioBitrate == 0 {
		// default audio bitrate if not found
		audioBitrate = 128
	}

	return duration, videoBitrate, audioBitrate, nil
}

// parse duration from ffmpeg output
func parseDuration(output string) float64 {
	if strings.Contains(output, "Duration:") {
		start := strings.Index(output, "Duration:") + 10
		end := start + 11
		durationStr := output[start:end]
		timeParts := strings.Split(durationStr, ":")
		hours, _ := strconv.ParseFloat(timeParts[0], 64)
		minutes, _ := strconv.ParseFloat(timeParts[1], 64)
		seconds, _ := strconv.ParseFloat(timeParts[2], 64)
		return hours*3600 + minutes*60 + seconds
	}
	return 0
}

// parse video bitrate from ffmpeg output
func parseVideoBitrate(output string) int {
	if strings.Contains(output, "bitrate:") {
		start := strings.Index(output, "bitrate:") + 9
		end := strings.Index(output[start:], " kb/s") + start
		bitrateStr := output[start:end]
		bitrate, err := strconv.Atoi(strings.TrimSpace(bitrateStr))
		if err != nil {
			return 0
		}
		return bitrate
	}
	return 0
}

// parse the audio bitrate from ffmpeg output
func parseAudioBitrate(output string) int {
	// look for the specific audio stream bitrate (ie 128 kb/s under the audio stream)
	if strings.Contains(output, "Audio:") && strings.Contains(output, " kb/s") {
		audioIdx := strings.Index(output, "Audio:")
		kbpsIdx := strings.Index(output[audioIdx:], " kb/s") + audioIdx
		start := kbpsIdx - 4
		audioBitrateStr := output[start:kbpsIdx]
		audioBitrate, err := strconv.Atoi(strings.TrimSpace(audioBitrateStr))
		if err != nil {
			return 0
		}
		return audioBitrate
	}
	return 0
}

// calculate the desired video bitrate based on desired output size and audio bitrate
func calculateDesiredBitrate(duration float64, desiredSizeMB int, audioBitrate int) (int, error) {
	// convert mb to kb
	desiredSizeKB := float64(desiredSizeMB * 1024)
	// calculate total bitrate in kbps
	totalBitrate := (desiredSizeKB * 8) / duration
	// subtract the audio bitrate to get video bitrate
	videoBitrate := totalBitrate - float64(audioBitrate)

	// set a minimum video bitrate at 100 kbps
	const minVideoBitrate = 100
	if videoBitrate < minVideoBitrate {
		return 0, fmt.Errorf("calculated video bitrate (%d kbps) is too low, desired output size may be too small", int(videoBitrate))
	}

	return int(math.Round(videoBitrate)), nil
}

// compress the video with the calculated bitrate
func compressVideo(inputFile string, desiredBitrate int) {
	baseName := filepath.Base(inputFile)
	outputFile := "compressed_" + baseName

	// run ffmpeg with the new video bitrate
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-b:v", fmt.Sprintf("%dk", desiredBitrate), outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("ffmpeg compression failed: %v", err)
	}

	fmt.Printf("video compressed successfully to %s with bitrate %d kbps\n", outputFile, desiredBitrate)
}

func main() {
	checkArgs()

	inputFile := os.Args[1]
	desiredSizeMB, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("invalid output size: %v", err)
	}

	// get video info like duration and bitrates
	duration, videoBitrate, audioBitrate, err := getVideoInfo(inputFile)
	if err != nil {
		log.Fatalf("failed to get video info: %v", err)
	}

	fmt.Printf("video duration: %.2f seconds\n", duration)
	fmt.Printf("original video bitrate: %d kbps\n", videoBitrate)
	fmt.Printf("audio bitrate: %d kbps\n", audioBitrate)

	// calculate what the new video bitrate should be
	desiredBitrate, err := calculateDesiredBitrate(duration, desiredSizeMB, audioBitrate)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("calculated video bitrate for desired output size: %d kbps\n", desiredBitrate)

	// run the compression with the new bitrate
	compressVideo(inputFile, desiredBitrate)
}
