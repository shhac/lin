package upload

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/linear"
)

type UploadedFile struct {
	Filename    string `json:"filename"`
	AssetURL    string `json:"assetUrl"`
	ContentType string `json:"contentType"`
}

func UploadFiles(client graphql.Client, paths []string) ([]UploadedFile, error) {
	results := make([]UploadedFile, 0, len(paths))
	for _, filePath := range paths {
		uploaded, err := uploadOne(client, filePath)
		if err != nil {
			return nil, err
		}
		results = append(results, uploaded)
	}
	return results, nil
}

// uploadOne handles a single file: stat, request an upload URL, PUT, and
// return the asset record. Errors are wrapped with the filename for context.
func uploadOne(client graphql.Client, filePath string) (UploadedFile, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return UploadedFile{}, fmt.Errorf("file not found: %s", filepath.Base(filePath))
	}

	filename := filepath.Base(filePath)
	contentType := detectMIME(filename)
	size := int(info.Size())

	resp, err := linear.FileUpload(context.Background(), client, contentType, filename, size)
	if err != nil {
		return UploadedFile{}, fmt.Errorf("upload failed for %s: %v", filename, err)
	}
	if resp.FileUpload.UploadFile == nil {
		return UploadedFile{}, fmt.Errorf("upload failed for %s: no upload URL returned", filename)
	}
	uf := resp.FileUpload.UploadFile

	headers := make(map[string]string, len(uf.Headers))
	for _, h := range uf.Headers {
		headers[h.Key] = h.Value
	}

	f, err := os.Open(filePath)
	if err != nil {
		return UploadedFile{}, err
	}
	defer func() { _ = f.Close() }()

	if err := httpPutWithHeaders(uf.UploadUrl, f, int64(size), headers); err != nil {
		return UploadedFile{}, fmt.Errorf("upload failed for %s: %v", filename, err)
	}

	return UploadedFile{
		Filename:    filename,
		AssetURL:    uf.AssetUrl,
		ContentType: contentType,
	}, nil
}

// httpPutWithHeaders streams body to url via HTTP PUT with the given headers
// and content length. Returns nil on 2xx, otherwise an error describing the
// response. Body is fully drained.
func httpPutWithHeaders(url string, body io.Reader, contentLength int64, headers map[string]string) error {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.ContentLength = contentLength

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 300 {
		return fmt.Errorf("%d %s", resp.StatusCode, resp.Status)
	}
	return nil
}

func FormatFileMarkdown(files []UploadedFile) string {
	var parts []string
	for _, f := range files {
		if strings.HasPrefix(f.ContentType, "image/") {
			parts = append(parts, fmt.Sprintf("![%s](%s)", f.Filename, f.AssetURL))
		} else {
			parts = append(parts, fmt.Sprintf("[%s](%s)", f.Filename, f.AssetURL))
		}
	}
	return strings.Join(parts, "\n")
}

func detectMIME(filename string) string {
	ext := filepath.Ext(filename)
	if ext != "" {
		if t := mime.TypeByExtension(ext); t != "" {
			return t
		}
	}
	return "application/octet-stream"
}
