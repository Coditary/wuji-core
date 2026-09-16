package data

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MarshalJSONShape encodes a WDD shape as JSON bytes.
func MarshalJSONShape(shape DataShape) ([]byte, error) {
	if shape == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(DocumentFromShape(shape))
}

// ParseJSONShape decodes a WDD JSON document into a shape.
func ParseJSONShape(raw []byte) (DataShape, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty shape json")
	}
	return ParseJSON(bytes.NewReader(raw))
}
