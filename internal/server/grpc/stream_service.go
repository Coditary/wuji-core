package grpc

import (
	"encoding/json"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (s *CoreService) GenerateTextStream(req *wujiv1.CallRequest, stream wujiv1.WujiCore_GenerateTextStreamServer) error {
	ctx := stream.Context()
	var body struct {
		DriverID string             `json:"driver_id"`
		Request  driver.TextRequest `json:"request"`
	}
	if err := json.Unmarshal(req.GetPayloadJson(), &body); err != nil {
		return stream.Send(&wujiv1.TextStreamChunk{Error: err.Error(), Done: true})
	}
	body.Request.Stream = true
	body.Request.OnDelta = func(delta string) error {
		return stream.Send(&wujiv1.TextStreamChunk{Delta: delta})
	}
	body.Request.OnThinkDelta = func(delta string) error {
		return stream.Send(&wujiv1.TextStreamChunk{ThinkDelta: delta})
	}

	resp, err := s.core.GenerateText(ctx, body.DriverID, body.Request)
	if err != nil {
		return stream.Send(&wujiv1.TextStreamChunk{Error: err.Error(), Done: true})
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		return stream.Send(&wujiv1.TextStreamChunk{Error: err.Error(), Done: true})
	}
	return stream.Send(&wujiv1.TextStreamChunk{PayloadJson: payload, Done: true})
}
