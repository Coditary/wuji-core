package data

import (
	"bytes"
	"fmt"
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

// DocumentFromShape builds a WDD document map for serialization.
func DocumentFromShape(s DataShape) map[string]any {
	switch v := s.(type) {
	case *VectorBatch:
		doc := envelopeFromRecordSet(v.ToRecordSet())
		doc["kind"] = string(ShapeVector)
		doc["tensor"] = map[string]any{
			"dims":   v.Dims,
			"count":  v.Count(),
			"dtype":  "f32",
			"layout": "row_major",
			"data":   v.Vectors,
		}
		return doc
	case *Table:
		doc := envelopeFromRecordSet(v.ToRecordSet())
		doc["kind"] = string(ShapeTable)
		return doc
	case *Graph:
		doc := envelopeFromRecordSet(v.ToRecordSet())
		doc["kind"] = string(ShapeGraph)
		return doc
	case *RecordSet:
		return envelopeFromRecordSet(v)
	default:
		return map[string]any{"wuji": "1"}
	}
}

// WriteMsgpack serializes a shape as MessagePack.
func WriteMsgpack(w io.Writer, s DataShape) error {
	return msgpack.NewEncoder(w).Encode(DocumentFromShape(s))
}

// ParseMsgpack reads a WDD document from MessagePack.
func ParseMsgpack(r io.Reader) (DataShape, error) {
	var doc map[string]any
	if err := msgpack.NewDecoder(r).Decode(&doc); err != nil {
		return nil, err
	}
	return parseDocumentMap(doc)
}

// MarshalMsgpack returns a WDD document encoded as MessagePack bytes.
func MarshalMsgpack(shape DataShape) ([]byte, error) {
	var buf bytes.Buffer
	if err := WriteMsgpack(&buf, shape); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalMsgpack decodes MessagePack bytes into a shape.
func UnmarshalMsgpack(payload []byte) (DataShape, error) {
	if len(payload) == 0 {
		return nil, fmt.Errorf("empty msgpack payload")
	}
	return ParseMsgpack(bytes.NewReader(payload))
}
