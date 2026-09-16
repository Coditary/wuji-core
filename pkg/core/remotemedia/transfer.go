package remotemedia

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/netx"
)

// Enabled reports whether the client should stage media for a core endpoint.
func Enabled(endpoint string) bool {
	if os.Getenv("WUJI_SKIP_MEDIA_TRANSFER") == "1" {
		return false
	}
	return !netx.IsUnix(endpoint)
}

// FetchTarget maps a server path to a preferred local destination.
type FetchTarget struct {
	Server string
	Local  string
}

// Broker moves local media to/from a remote core daemon.
type Broker struct {
	client    wujiv1.WujiCoreClient
	localRoot string
}

func NewBroker(client wujiv1.WujiCoreClient, localRoot string) *Broker {
	if localRoot == "" {
		localRoot = "."
	}
	return &Broker{client: client, localRoot: localRoot}
}

func (b *Broker) sessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Upload reads local files and directories; returns localPath -> serverPath.
func (b *Broker) Upload(ctx context.Context, localPaths []string) (map[string]string, error) {
	entries, dirUploads, err := BuildUploadPlan(localPaths)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return map[string]string{}, nil
	}

	sessionID := b.sessionID()
	blobs := make([]*wujiv1.FileBlob, 0, len(entries))
	for _, e := range entries {
		data, err := os.ReadFile(e.LocalPath)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", e.LocalPath, err)
		}
		payload, zstdFlag, err := BlobData(data)
		if err != nil {
			return nil, err
		}
		blobs = append(blobs, &wujiv1.FileBlob{
			Name: e.BlobName,
			Data: payload,
			Zstd: zstdFlag,
		})
	}

	resp, err := b.client.PutFiles(ctx, &wujiv1.PutFilesRequest{
		SessionId: sessionID,
		Files:     blobs,
	})
	if err != nil {
		return nil, fmt.Errorf("upload media: %w", err)
	}
	if len(resp.GetPaths()) != len(entries) {
		return nil, fmt.Errorf("upload media: expected %d paths, got %d", len(entries), len(resp.GetPaths()))
	}

	mapping := make(map[string]string)
	stagingRoot := resp.GetStagingRoot()
	for i, e := range entries {
		if e.RootPath != "" {
			mapping[e.RootPath] = ServerDirForLocal(stagingRoot, e.RootPath)
			continue
		}
		mapping[e.LocalPath] = resp.GetPaths()[i]
	}
	for _, d := range dirUploads {
		if mapping[d.LocalDir] == "" && stagingRoot != "" {
			mapping[d.LocalDir] = ServerDirForLocal(stagingRoot, d.LocalDir)
		}
	}
	return mapping, nil
}

// Download fetches server files; returns serverPath -> localPath.
func (b *Broker) Download(ctx context.Context, serverPaths []string) (map[string]string, error) {
	targets := make([]FetchTarget, 0, len(serverPaths))
	for _, p := range serverPaths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		targets = append(targets, FetchTarget{Server: p, Local: ""})
	}
	return b.DownloadTo(ctx, targets)
}

// DownloadTo fetches server files into preferred local paths when set.
func (b *Broker) DownloadTo(ctx context.Context, targets []FetchTarget) (map[string]string, error) {
	var serverPaths []string
	for _, t := range targets {
		if strings.TrimSpace(t.Server) == "" {
			continue
		}
		serverPaths = append(serverPaths, t.Server)
	}
	if len(serverPaths) == 0 {
		return map[string]string{}, nil
	}

	resp, err := b.client.GetFiles(ctx, &wujiv1.GetFilesRequest{Paths: serverPaths})
	if err != nil {
		return nil, fmt.Errorf("download media: %w", err)
	}

	sessionID := b.sessionID()
	outBase := filepath.Join(b.localRoot, ".wuji", "remote-out", sanitizeSessionID(sessionID))
	mapping := make(map[string]string)

	// One server path may expand to many blobs (directory).
	if len(resp.GetFiles()) == len(serverPaths) {
		for i, serverPath := range serverPaths {
			blob := resp.GetFiles()[i]
			localPath := preferredLocal(targets, serverPath)
			if localPath == "" {
				localPath = filepath.Join(outBase, filepath.Base(serverPath))
			}
			raw, err := MaybeDecompress(blob.GetData(), blob.GetZstd())
			if err != nil {
				return nil, err
			}
			if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
				return nil, err
			}
			if err := os.WriteFile(localPath, raw, 0o644); err != nil {
				return nil, fmt.Errorf("write %s: %w", localPath, err)
			}
			mapping[serverPath] = localPath
		}
		return mapping, nil
	}

	// Directory tree or batched blobs — write by blob name under outBase.
	prefix := commonBlobPrefix(resp.GetFiles())
	for _, blob := range resp.GetFiles() {
		raw, err := MaybeDecompress(blob.GetData(), blob.GetZstd())
		if err != nil {
			return nil, err
		}
		name := strings.TrimPrefix(blob.GetName(), prefix)
		if name == "" {
			name = blob.GetName()
		}
		localPath := filepath.Join(outBase, name)
		if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(localPath, raw, 0o644); err != nil {
			return nil, fmt.Errorf("write %s: %w", localPath, err)
		}
	}
	for _, t := range targets {
		if t.Local != "" {
			mapping[t.Server] = t.Local
			continue
		}
		base := filepath.Base(t.Server)
		mapping[t.Server] = filepath.Join(outBase, base)
	}
	return mapping, nil
}

func preferredLocal(targets []FetchTarget, serverPath string) string {
	for _, t := range targets {
		if t.Server == serverPath && t.Local != "" {
			return t.Local
		}
	}
	return ""
}

func commonBlobPrefix(blobs []*wujiv1.FileBlob) string {
	if len(blobs) == 0 {
		return ""
	}
	first := blobs[0].GetName()
	if !strings.Contains(first, "/") && !strings.Contains(first, string(os.PathSeparator)) {
		return ""
	}
	parts := strings.Split(strings.ReplaceAll(first, "\\", "/"), "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[0] + "/"
}
