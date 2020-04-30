package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Structs and their fields must be exported for json parsing

// FFProbeJSON is for json parsing
type FFProbeJSON struct {
	Packets []*FFProbePacketJSON `json:"packets"`
	Format  *FFProbeFormatJSON   `json:"format"`
}

// FFProbePacketJSON is for json parsing
type FFProbePacketJSON struct {
	PacketTime     string `json:"dts_time"`
	PacketDuration string `json:"duration_time"`
}

// FFProbeFormatJSON is for json parsing
type FFProbeFormatJSON struct {
	FormatName string `json:"format_name"`
}

func parseAudioDuration(timeStr string) (time.Duration, error) {
	timeParts := strings.Split(timeStr, ":")
	hour, minute := timeParts[0], timeParts[1]
	secUs1 := strings.Split(timeParts[2], ".")
	sec := secUs1[0]
	microsec := secUs1[1]

	return time.ParseDuration(hour + "h" + minute + "m" + sec + "s" + microsec + "us")
}

func parseAudioTime(timeStr string) (time.Time, error) {
	// We must have a double digit hour value
	if timeStr[1] == ':' {
		timeStr = "0" + timeStr
	}
	// timeStr goes down to microseconds, but datetime only supports nanoseconds
	return time.Parse("15:04:05.000000", timeStr)
}

// ProbeAudio returns the audio's codec and duration in seconds.
// ffmpegDir is the directory that ffmpeg lives in, and can be left blank if ffmpeg is in your PATH
func ProbeAudio(input io.Reader, ffmpegDir string) (string, int) {
	// ffprobe -show_format -pretty -loglevel quiet -print_format json -show_packets pipe:
	cmd := exec.Command(ffmpegDir+"ffprobe", "-show_format", "-pretty", "-loglevel", "quiet",
		"-print_format", "json", "-show_packets", "pipe:")

	cmd.Stdin = input

	jsonBytes, err := cmd.Output()
	if err != nil {
		log.Fatalln(err)
	}

	var jsonStruct FFProbeJSON
	json.Unmarshal(jsonBytes, &jsonStruct)

	lastPacket := jsonStruct.Packets[len(jsonStruct.Packets)-1]
	totalDuration, err := parseAudioTime(lastPacket.PacketTime)
	if err != nil {
		log.Fatalln(err)
	}
	packetDuration, err := parseAudioDuration(lastPacket.PacketDuration)
	if err != nil {
		log.Fatalln(err)
	}
	// total duration = last packet time + last packet duration
	totalDuration.Add(packetDuration)

	duration := totalDuration.Hour()*3600 + totalDuration.Minute()*60 + totalDuration.Second()

	return jsonStruct.Format.FormatName, duration
}

//EncodeAudio encodes audio from input to output. Bitrate is in kbps
func EncodeAudio(bitrate int, inCodec string, outCodec string, input io.Reader, output io.Writer, ffmpegDir string) {
	fmt.Println(">> Encoding", inCodec, "file to", outCodec)
	encoders := map[string]string{
		"opus": "libopus",
		"mp3":  "libmp3lame",
	}

	cmd := exec.Command(ffmpegDir+"ffmpeg", "-f", inCodec, "-i", "pipe:", "-y", "-c:a", encoders[outCodec],
		"-ac", "1", "-b:a", strconv.Itoa(bitrate)+"k", "-f", outCodec, "pipe:")

	fmt.Println(">> Executing: " + strings.Join(cmd.Args, " "))

	cmd.Stdin = input
	cmd.Stdout = output
	// cmd.Stderr = os.Stderr

	// Wait for command to complete.
	err := cmd.Run()

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(">> Encoded file.")
	return
}
