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
	var results []UploadedFile
	for _, filePath := range paths {
		info, err := os.Stat(filePath)
		if err != nil {
			return nil, fmt.Errorf("File not found: %s", filepath.Base(filePath))
		}

		filename := filepath.Base(filePath)
		contentType := detectMIME(filename)
		size := int(info.Size())

		resp, err := linear.FileUpload(context.Background(), client, contentType, filename, size)
		if err != nil {
			return nil, fmt.Errorf("Upload failed for %s: %v", filename, err)
		}
		if resp.FileUpload.UploadFile == nil {
			return nil, fmt.Errorf("Upload failed for %s: no upload URL returned", filename)
		}

		uf := resp.FileUpload.UploadFile

		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("PUT", uf.UploadUrl, f)
		if err != nil {
			f.Close()
			return nil, err
		}
		for _, h := range uf.Headers {
			req.Header.Set(h.Key, h.Value)
		}
		req.ContentLength = int64(size)

		putResp, err := http.DefaultClient.Do(req)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("Upload failed for %s: %v", filename, err)
		}
		io.Copy(io.Discard, putResp.Body)
		putResp.Body.Close()
		if putResp.StatusCode >= 300 {
			return nil, fmt.Errorf("Upload failed for %s: %d %s", filename, putResp.StatusCode, putResp.Status)
		}

		results = append(results, UploadedFile{
			Filename:    filename,
			AssetURL:    uf.AssetUrl,
			ContentType: contentType,
		})
	}
	return results, nil
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
