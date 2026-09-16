package driver

import (
	"context"

	"github.com/coditary/wuji-core/pkg/data"
)

// TextGenerator generates text from prompts.
type TextGenerator interface {
	GenerateText(ctx context.Context, req TextRequest) (*TextResponse, error)
}

// VideoToAudioConverter extracts an audio track from a video container.
type VideoToAudioConverter interface {
	VideoToAudio(ctx context.Context, req Video2AudioRequest) (*Video2AudioResponse, error)
}

// ImageToTextGenerator describes an image and returns plain text.
type ImageToTextGenerator interface {
	ImageToText(ctx context.Context, req TextRequest) (*TextResponse, error)
}

// VideoToTextGenerator understands video natively and returns plain text.
type VideoToTextGenerator interface {
	VideoToText(ctx context.Context, req TextRequest) (*TextResponse, error)
}

// DocumentToTextGenerator extracts plain text from a document file.
type DocumentToTextGenerator interface {
	DocumentToText(ctx context.Context, req TextRequest) (*TextResponse, error)
}

// AudioToTextGenerator transcribes audio to plain text.
type AudioToTextGenerator interface {
	AudioToText(ctx context.Context, req TextRequest) (*TextResponse, error)
}

// ImageGenerator generates images from text prompts (text-to-image).
type ImageGenerator interface {
	GenerateImage(ctx context.Context, req ImageGenerateRequest) (*ImageResponse, error)
}

// ImageToImageGenerator transforms an image guided by a prompt.
type ImageToImageGenerator interface {
	ImageToImage(ctx context.Context, req ImageToImageRequest) (*ImageResponse, error)
}

// ImageInpaintGenerator fills masked regions in an image.
type ImageInpaintGenerator interface {
	InpaintImage(ctx context.Context, req ImageInpaintRequest) (*ImageResponse, error)
}

// ImageUpscaleGenerator increases image resolution with AI super-resolution.
type ImageUpscaleGenerator interface {
	UpscaleImage(ctx context.Context, req ImageUpscaleRequest) (*ImageResponse, error)
}

// ImageDownscaleGenerator reduces image resolution (e.g. SSIM-aware, Lanczos, …).
type ImageDownscaleGenerator interface {
	DownscaleImage(ctx context.Context, req ImageDownscaleRequest) (*ImageResponse, error)
}

// ImageScaleGenerator handles arbitrary scale factors (e.g. 3×, 2.5×) in one step.
type ImageScaleGenerator interface {
	ScaleImage(ctx context.Context, req ImageScaleRequest) (*ImageResponse, error)
}

// ImageEditGenerator edits an image from an instruction prompt.
type ImageEditGenerator interface {
	EditImage(ctx context.Context, req ImageEditRequest) (*ImageResponse, error)
}

// ImageControlNetGenerator generates images with structural control.
type ImageControlNetGenerator interface {
	ControlNetImage(ctx context.Context, req ImageControlNetRequest) (*ImageResponse, error)
}

// ImageDepthToImageGenerator generates images from a depth map.
type ImageDepthToImageGenerator interface {
	DepthToImage(ctx context.Context, req ImageDepthToImageRequest) (*ImageResponse, error)
}

// ImageVariationGenerator creates variations of an input image.
type ImageVariationGenerator interface {
	ImageVariation(ctx context.Context, req ImageVariationRequest) (*ImageResponse, error)
}

// ImageStyleTransferGenerator applies a reference style to generation.
type ImageStyleTransferGenerator interface {
	StyleTransferImage(ctx context.Context, req ImageStyleTransferRequest) (*ImageResponse, error)
}

// ImageSpriteGenerator generates sprite sheet atlases.
type ImageSpriteGenerator interface {
	SpriteImage(ctx context.Context, req ImageSpriteRequest) (*ImageResponse, error)
}

// VideoGenerator generates videos from text prompts (text-to-video).
type VideoGenerator interface {
	GenerateVideo(ctx context.Context, req VideoGenerateRequest) (*VideoResponse, error)
}

// VideoImageToVideoGenerator animates an image into video.
type VideoImageToVideoGenerator interface {
	ImageToVideo(ctx context.Context, req VideoImageToVideoRequest) (*VideoResponse, error)
}

// VideoInterpolateGenerator increases frame rate via interpolation.
type VideoInterpolateGenerator interface {
	InterpolateVideo(ctx context.Context, req VideoInterpolateRequest) (*VideoResponse, error)
}

// VideoUpscaleGenerator increases video resolution.
type VideoUpscaleGenerator interface {
	UpscaleVideo(ctx context.Context, req VideoUpscaleRequest) (*VideoResponse, error)
}

// MusicGenerator generates music from a text prompt.
type MusicGenerator interface {
	GenerateMusic(ctx context.Context, req AudioMusicRequest) (*AudioResponse, error)
}

// MelodyToMusicGenerator generates music guided by a reference melody.
type MelodyToMusicGenerator interface {
	MelodyToMusic(ctx context.Context, req AudioMelodyToMusicRequest) (*AudioResponse, error)
}

// SFXGenerator generates sound effects from a text prompt.
type SFXGenerator interface {
	GenerateSFX(ctx context.Context, req AudioSFXRequest) (*AudioResponse, error)
}

// AudioSpeechGenerator generates spoken audio via an audio model.
type AudioSpeechGenerator interface {
	GenerateAudioSpeech(ctx context.Context, req AudioSpeechRequest) (*AudioResponse, error)
}

// MeshGenerator generates or transforms mesh assets.
type MeshGenerator interface {
	GenerateMesh(ctx context.Context, req MeshRequest) (*MeshResponse, error)
}

// VoiceCreator creates a voice profile from a sample.
type VoiceCreator interface {
	CreateVoice(ctx context.Context, req VoiceCreateRequest) (*VoiceResponse, error)
}

// VoiceConverter converts source audio to a cloned voice.
type VoiceConverter interface {
	ConvertVoice(ctx context.Context, req VoiceConvertRequest) (*VoiceResponse, error)
}

// VoiceCatalogProvider lists speakable voice profiles and their supported speech parameters.
type VoiceCatalogProvider interface {
	ListVoiceProfiles(ctx context.Context) ([]VoiceProfileInfo, error)
}

// TextTrainer fine-tunes text generation models.
type TextTrainer interface {
	TrainText(ctx context.Context, req TextTrainRequest) (*TrainResponse, error)
}

// ImageTrainer fine-tunes image generation models.
type ImageTrainer interface {
	TrainImage(ctx context.Context, req ImageTrainRequest) (*TrainResponse, error)
}

// VideoTrainer fine-tunes video generation models.
type VideoTrainer interface {
	TrainVideo(ctx context.Context, req VideoTrainRequest) (*TrainResponse, error)
}

// AudioTrainer fine-tunes audio generation models.
type AudioTrainer interface {
	TrainAudio(ctx context.Context, req AudioTrainRequest) (*TrainResponse, error)
}

// MeshTrainer fine-tunes mesh generation models.
type MeshTrainer interface {
	TrainMesh(ctx context.Context, req MeshTrainRequest) (*TrainResponse, error)
}

// VoiceTrainer fine-tunes voice cloning models.
type VoiceTrainer interface {
	TrainVoice(ctx context.Context, req VoiceTrainRequest) (*TrainResponse, error)
}

// DataProducer runs structured data tasks (embed, detect, forecast, graph, …).
type DataProducer interface {
	ProduceData(ctx context.Context, req DataRequest) (data.DataShape, error)
}

// RAGProducer runs RAG tasks and returns WDD output (used by remote drivers).
type RAGProducer interface {
	ProduceRAG(ctx context.Context, req RAGRequest) (data.DataShape, error)
}

// RAGIndexer ingests documents into a collection.
type RAGIndexer interface {
	Index(ctx context.Context, req RAGRequest) (*RAGIndexResponse, error)
}

// RAGQuerier retrieves relevant chunks for a query.
type RAGQuerier interface {
	Query(ctx context.Context, req RAGRequest) (*RAGQueryResponse, error)
}

// RAGAnswerer generates an answer with retrieval (optional driver-native path).
type RAGAnswerer interface {
	Answer(ctx context.Context, req RAGRequest) (*RAGAnswerResponse, error)
}

// RAGCollectionManager lists and maintains collections.
type RAGCollectionManager interface {
	ListCollections(ctx context.Context, req RAGRequest) ([]RAGCollectionInfo, error)
	CollectionInfo(ctx context.Context, req RAGRequest) (RAGCollectionInfo, error)
	DeleteCollection(ctx context.Context, req RAGRequest) error
	PurgeSource(ctx context.Context, req RAGRequest) error
	ExportCollection(ctx context.Context, req RAGRequest) error
	ImportCollection(ctx context.Context, req RAGRequest) error
	RenameCollection(ctx context.Context, req RAGRequest) error
}

// DatasetLister lists datasets.
type DatasetLister interface {
	ListDatasets(ctx context.Context, req DatasetListRequest) (*DatasetResponse, error)
}

// DatasetCreator creates datasets.
type DatasetCreator interface {
	CreateDataset(ctx context.Context, req DatasetCreateRequest) (*DatasetResponse, error)
}

// DatasetDeleter deletes datasets.
type DatasetDeleter interface {
	DeleteDataset(ctx context.Context, req DatasetDeleteRequest) (*DatasetResponse, error)
}

// DatasetIngester imports files into a dataset.
type DatasetIngester interface {
	IngestDataset(ctx context.Context, req DatasetIngestRequest) (*DatasetIngestResult, error)
}

// DatasetVersioner creates dataset snapshots.
type DatasetVersioner interface {
	CreateDatasetVersion(ctx context.Context, req DatasetVersionRequest) (*DatasetVersionResult, error)
}

// DatasetVersionLister lists dataset versions.
type DatasetVersionLister interface {
	ListDatasetVersions(ctx context.Context, req DatasetListVersionsRequest) ([]DatasetVersionEntry, error)
}
