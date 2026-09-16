package driver

import wujiv1 "github.com/coditary/wuji-core/api/proto/v1"

func Video2AudioRequestFromProto(req *wujiv1.Video2AudioRequest) Video2AudioRequest {
	if req == nil {
		return Video2AudioRequest{}
	}
	return Video2AudioRequest{VideoPath: req.GetVideoPath()}
}

func Video2AudioRequestToProto(req Video2AudioRequest) *wujiv1.Video2AudioRequest {
	return &wujiv1.Video2AudioRequest{VideoPath: req.VideoPath}
}

func Video2AudioResponseFromProto(resp *wujiv1.Video2AudioResponse) *Video2AudioResponse {
	if resp == nil {
		return nil
	}
	return &Video2AudioResponse{
		AudioPath: resp.GetAudioPath(),
		Temporary: resp.GetTemporary(),
	}
}

func Video2AudioResponseToProto(resp *Video2AudioResponse) *wujiv1.Video2AudioResponse {
	if resp == nil {
		return &wujiv1.Video2AudioResponse{}
	}
	return &wujiv1.Video2AudioResponse{
		AudioPath: resp.AudioPath,
		Temporary: resp.Temporary,
	}
}
