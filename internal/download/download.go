package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/linear"
)

const uploadHost = "uploads.linear.app"

type DownloadOpts struct {
	Output    string
	OutputDir string
	Stdout    bool
	Force     bool
	APIKey    string
}

type DownloadResult struct {
	Filename    string `json:"filename"`
	Path        string `json:"path"`
	Size        int    `json:"size"`
	ContentType string `json:"contentType"`
}

func GetOrgID(client graphql.Client) (string, error) {
	resp, err := linear.Organization(context.Background(), client)
	if err != nil {
		return "", err
	}
	return resp.Organization.Id, nil
}

// fetchedFile is the network-side payload of a download, before any
// disk-vs-stdout routing decisions.
type fetchedFile struct {
	Data        []byte
	ContentType string
	Filename    string
}

// DownloadFile fetches fileURL and routes the bytes to stdout or disk per opts.
func DownloadFile(fileURL string, opts DownloadOpts) (DownloadResult, error) {
	f, err := fetchFile(fileURL, opts.APIKey)
	if err != nil {
		return DownloadResult{}, err
	}
	return writeFetched(f, opts)
}

// fetchFile performs the authenticated HTTP GET and reads the body. It does
// not touch the local filesystem.
func fetchFile(fileURL, apiKey string) (fetchedFile, error) {
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return fetchedFile{}, err
	}
	req.Header.Set("Authorization", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fetchedFile{}, fmt.Errorf("download failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return fetchedFile{}, fmt.Errorf("download failed: %d %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fetchedFile{}, fmt.Errorf("download failed: %v", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return fetchedFile{
		Data:        data,
		ContentType: contentType,
		Filename:    inferFilename(resp.Header, fileURL, contentType),
	}, nil
}

// writeFetched routes f to stdout or a resolved local path per opts and
// returns the user-facing result.
func writeFetched(f fetchedFile, opts DownloadOpts) (DownloadResult, error) {
	if opts.Stdout {
		_, _ = os.Stdout.Write(f.Data)
		return DownloadResult{
			Filename:    f.Filename,
			Path:        "<stdout>",
			Size:        len(f.Data),
			ContentType: f.ContentType,
		}, nil
	}

	destPath, err := resolveDestPath(f.Filename, opts, f.ContentType)
	if err != nil {
		return DownloadResult{}, err
	}
	if err := os.WriteFile(destPath, f.Data, 0o644); err != nil {
		return DownloadResult{}, err
	}

	absPath, _ := filepath.Abs(destPath)
	return DownloadResult{
		Filename:    filepath.Base(destPath),
		Path:        absPath,
		Size:        len(f.Data),
		ContentType: f.ContentType,
	}, nil
}
