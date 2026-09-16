package remotemedia

import (
	"fmt"

	"github.com/klauspost/compress/zstd"
)

var (
	zstdEncoder, _ = zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedDefault))
	zstdDecoder, _ = zstd.NewReader(nil)
)

// CompressZstd returns zstd-compressed bytes.
func CompressZstd(raw []byte) ([]byte, error) {
	if len(raw) == 0 {
		return raw, nil
	}
	if zstdEncoder == nil {
		return nil, fmt.Errorf("zstd encoder unavailable")
	}
	return zstdEncoder.EncodeAll(raw, make([]byte, 0, len(raw)/2)), nil
}

// DecompressZstd returns decompressed bytes.
func DecompressZstd(compressed []byte) ([]byte, error) {
	if len(compressed) == 0 {
		return compressed, nil
	}
	if zstdDecoder == nil {
		return nil, fmt.Errorf("zstd decoder unavailable")
	}
	return zstdDecoder.DecodeAll(compressed, nil)
}

// MaybeDecompress returns data, decompressing when zstd is set.
func MaybeDecompress(data []byte, zstdFlag bool) ([]byte, error) {
	if !zstdFlag {
		return data, nil
	}
	return DecompressZstd(data)
}

// BlobData prepares payload bytes for transfer (compress when worthwhile).
func BlobData(raw []byte) ([]byte, bool, error) {
	if len(raw) < 512 {
		return raw, false, nil
	}
	compressed, err := CompressZstd(raw)
	if err != nil {
		return nil, false, err
	}
	if len(compressed) >= len(raw) {
		return raw, false, nil
	}
	return compressed, true, nil
}
