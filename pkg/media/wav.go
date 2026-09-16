package media

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

// WritePlaceholderWAV writes a mono 16-bit PCM WAV file filled with silence.
// Used by test drivers until a real audio backend is connected.
func WritePlaceholderWAV(path string, durationSec float32, sampleRate int) error {
	if durationSec <= 0 {
		durationSec = 1
	}
	if sampleRate <= 0 {
		sampleRate = 44100
	}

	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}
	}

	numSamples := int(float64(durationSec) * float64(sampleRate))
	if numSamples < 1 {
		numSamples = 1
	}
	dataSize := numSamples * 2 // 16-bit mono

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create wav file: %w", err)
	}
	defer f.Close()

	header := make([]byte, 44)
	copy(header[0:4], "RIFF")
	binary.LittleEndian.PutUint32(header[4:8], uint32(36+dataSize))
	copy(header[8:12], "WAVE")
	copy(header[12:16], "fmt ")
	binary.LittleEndian.PutUint32(header[16:20], 16)
	binary.LittleEndian.PutUint16(header[20:22], 1) // PCM
	binary.LittleEndian.PutUint16(header[22:24], 1) // mono
	binary.LittleEndian.PutUint32(header[24:28], uint32(sampleRate))
	byteRate := sampleRate * 2
	binary.LittleEndian.PutUint32(header[28:32], uint32(byteRate))
	binary.LittleEndian.PutUint16(header[32:34], 2) // block align
	binary.LittleEndian.PutUint16(header[34:36], 16) // bits per sample
	copy(header[36:40], "data")
	binary.LittleEndian.PutUint32(header[40:44], uint32(dataSize))

	if _, err := f.Write(header); err != nil {
		return fmt.Errorf("write wav header: %w", err)
	}

	chunk := make([]byte, 4096)
	for written := 0; written < dataSize; written += len(chunk) {
		n := len(chunk)
		if remaining := dataSize - written; remaining < n {
			n = remaining
		}
		if _, err := f.Write(chunk[:n]); err != nil {
			return fmt.Errorf("write wav data: %w", err)
		}
	}
	return nil
}
