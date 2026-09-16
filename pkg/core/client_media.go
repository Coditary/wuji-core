package core

import (
	"context"

	"github.com/coditary/wuji-core/pkg/core/remotemedia"
	"github.com/coditary/wuji-core/pkg/driver"
)

func (c *Client) stage(ctx context.Context, paths []string) (map[string]string, error) {
	if c.broker == nil {
		return map[string]string{}, nil
	}
	return c.broker.Upload(ctx, paths)
}

func (c *Client) fetch(ctx context.Context, paths []string) (map[string]string, error) {
	if c.broker == nil {
		return map[string]string{}, nil
	}
	return c.broker.Download(ctx, paths)
}

func (c *Client) fetchTo(ctx context.Context, targets []remotemedia.FetchTarget) (map[string]string, error) {
	if c.broker == nil {
		return map[string]string{}, nil
	}
	return c.broker.DownloadTo(ctx, targets)
}

func (c *Client) prepareText(ctx context.Context, req *driver.TextRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectTextRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteTextRequest(req, m)
	return nil
}

func (c *Client) finalizeImage(ctx context.Context, resp *driver.ImageResponse) error {
	m, err := c.fetch(ctx, remotemedia.CollectImageResponse(*resp))
	if err != nil {
		return err
	}
	remotemedia.RewriteImageResponse(resp, m)
	return nil
}

func (c *Client) prepareImage(ctx context.Context, req *driver.ImageRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectImageRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteImageRequest(req, m)
	return nil
}

func (c *Client) prepareVideo(ctx context.Context, req *driver.VideoRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectVideoRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteVideoRequest(req, m)
	return nil
}

func (c *Client) finalizeVideo(ctx context.Context, resp *driver.VideoResponse) error {
	m, err := c.fetch(ctx, remotemedia.CollectVideoResponse(*resp))
	if err != nil {
		return err
	}
	remotemedia.RewriteVideoResponse(resp, m)
	return nil
}

func (c *Client) prepareAudio(ctx context.Context, req *driver.AudioRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectAudioRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteAudioRequest(req, m)
	return nil
}

func (c *Client) finalizeAudio(ctx context.Context, resp *driver.AudioResponse) error {
	m, err := c.fetch(ctx, remotemedia.CollectAudioResponse(*resp))
	if err != nil {
		return err
	}
	remotemedia.RewriteAudioResponse(resp, m)
	return nil
}

func (c *Client) prepareMesh(ctx context.Context, req *driver.MeshRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectMeshRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteMeshRequest(req, m)
	return nil
}

func (c *Client) finalizeMesh(ctx context.Context, resp *driver.MeshResponse) error {
	m, err := c.fetch(ctx, remotemedia.CollectMeshResponse(*resp))
	if err != nil {
		return err
	}
	remotemedia.RewriteMeshResponse(resp, m)
	return nil
}

func (c *Client) prepareData(ctx context.Context, req *driver.DataRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectDataRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteDataRequest(req, m)
	return nil
}

func (c *Client) prepareRAG(ctx context.Context, req *driver.RAGRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectRAGRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteRAGRequest(req, m)
	return nil
}

func (c *Client) finalizeRAG(ctx context.Context, req driver.RAGRequest) error {
	targets := remotemedia.RAGExportDownload(req)
	if len(targets) == 0 {
		return nil
	}
	_, err := c.fetchTo(ctx, targets)
	return err
}

func (c *Client) prepareDataset(ctx context.Context, req *driver.DatasetRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectDatasetRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteDatasetRequest(req, m)
	return nil
}

func (c *Client) prepareTextTrain(ctx context.Context, req *driver.TextTrainRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectTextTrainRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteTextTrainRequest(req, m)
	return nil
}

func (c *Client) prepareImageTrain(ctx context.Context, req *driver.ImageTrainRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectImageTrainRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteImageTrainRequest(req, m)
	return nil
}

func (c *Client) prepareVideoTrain(ctx context.Context, req *driver.VideoTrainRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectVideoTrainRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteVideoTrainRequest(req, m)
	return nil
}

func (c *Client) prepareAudioTrain(ctx context.Context, req *driver.AudioTrainRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectAudioTrainRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteAudioTrainRequest(req, m)
	return nil
}

func (c *Client) prepareMeshTrain(ctx context.Context, req *driver.MeshTrainRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectMeshTrainRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteMeshTrainRequest(req, m)
	return nil
}

func (c *Client) prepareVoiceTrain(ctx context.Context, req *driver.VoiceTrainRequest) error {
	m, err := c.stage(ctx, remotemedia.CollectVoiceTrainRequest(*req))
	if err != nil {
		return err
	}
	remotemedia.RewriteVoiceTrainRequest(req, m)
	return nil
}

func (c *Client) finalizeTrain(ctx context.Context, resp *driver.TrainResponse) error {
	m, err := c.fetch(ctx, remotemedia.CollectTrainResponse(*resp))
	if err != nil {
		return err
	}
	remotemedia.RewriteTrainResponse(resp, m)
	return nil
}
