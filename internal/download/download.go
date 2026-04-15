package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/linear"
)

const uploadHost = "uploads.linear.app"

var uuidRE = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// eslint-disable-next-line no-control-regex
var unsafeFilenameRE = regexp.MustCompile(`[<>:"|?*\x00-\x1f]`)

var mimeToExt = map[string]string{
	"image/png":       ".png",
	"image/jpeg":      ".jpg",
	"image/gif":       ".gif",
	"image/webp":      ".webp",
	"image/svg+xml":   ".svg",
	"application/pdf": ".pdf",
	"text/plain":      ".txt",
	"text/csv":        ".csv",
	"application/json": ".json",
	"application/zip": ".zip",
	"video/mp4":       ".mp4",
	"audio/mpeg":      ".mp3",
}

type ParsedFileURL struct {
	URL      string
	OrgID    string
	Segments []string
}

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

func ParseFileURL(input string, defaultOrgID string) (ParsedFileURL, error) {
	var pathname string

	if strings.HasPrefix(input, "http://") {
		return ParsedFileURL{}, fmt.Errorf("Refusing http:// URL — only https:// is allowed for file downloads.")
	}

	if strings.HasPrefix(input, "https://") {
		u, err := url.Parse(input)
		if err != nil {
			return ParsedFileURL{}, fmt.Errorf("Invalid URL: %v", err)
		}
		if u.Hostname() != uploadHost {
			return ParsedFileURL{}, fmt.Errorf("Invalid host: %q. Only %s URLs are supported.", u.Hostname(), uploadHost)
		}
		pathname = u.Path
	} else if strings.HasPrefix(input, uploadHost+"/") {
		pathname = input[len(uploadHost):]
	} else {
		pathname = "/" + input
	}

	var segments []string
	for _, s := range strings.Split(pathname, "/") {
		if s != "" {
			segments = append(segments, s)
		}
	}

	if len(segments) == 0 || len(segments) > 3 {
		return ParsedFileURL{}, fmt.Errorf("Cannot parse file URL: %q. Expected 1-3 UUID path segments.", input)
	}

	for _, seg := range segments {
		if !uuidRE.MatchString(seg) {
			return ParsedFileURL{}, fmt.Errorf("Invalid UUID segment: %q. Expected format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", seg)
		}
	}

	var orgID string
	var fileSegments []string

	if len(segments) == 3 {
		orgID = segments[0]
		fileSegments = segments[1:]
	} else {
		if defaultOrgID == "" {
			return ParsedFileURL{}, fmt.Errorf("Cannot infer organization ID. Provide a full URL with org segment, or authenticate first.")
		}
		orgID = defaultOrgID
		fileSegments = segments
	}

	allSegments := append([]string{orgID}, fileSegments...)
	fullURL := fmt.Sprintf("https://%s/%s", uploadHost, strings.Join(allSegments, "/"))

	return ParsedFileURL{URL: fullURL, OrgID: orgID, Segments: allSegments}, nil
}

func GetOrgID(client graphql.Client) (string, error) {
	resp, err := linear.Organization(context.Background(), client)
	if err != nil {
		return "", err
	}
	return resp.Organization.Id, nil
}

func DownloadFile(fileURL string, opts DownloadOpts) (DownloadResult, error) {
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return DownloadResult{}, err
	}
	req.Header.Set("Authorization", opts.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return DownloadResult{}, fmt.Errorf("Download failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return DownloadResult{}, fmt.Errorf("Download failed: %d %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return DownloadResult{}, fmt.Errorf("Download failed: %v", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	filename := inferFilename(resp.Header, fileURL, contentType)

	if opts.Stdout {
		os.Stdout.Write(data)
		return DownloadResult{
			Filename:    filename,
			Path:        "<stdout>",
			Size:        len(data),
			ContentType: contentType,
		}, nil
	}

	destPath, err := resolveDestPath(filename, opts, contentType)
	if err != nil {
		return DownloadResult{}, err
	}

	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return DownloadResult{}, err
	}

	absPath, _ := filepath.Abs(destPath)
	return DownloadResult{
		Filename:    filepath.Base(destPath),
		Path:        absPath,
		Size:        len(data),
		ContentType: contentType,
	}, nil
}

func inferFilename(headers http.Header, rawURL, contentType string) string {
	if disposition := headers.Get("Content-Disposition"); disposition != "" {
		if parsed := parseContentDispositionFilename(disposition); parsed != "" {
			return SanitizeFilename(parsed)
		}
	}

	if u, err := url.Parse(rawURL); err == nil {
		parts := strings.Split(u.Path, "/")
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] != "" {
				ext := mimeToExtension(contentType)
				if ext != "" {
					return parts[i] + ext
				}
				return parts[i]
			}
		}
	}

	return "download"
}

func parseContentDispositionFilename(header string) string {
	// RFC 5987: filename*=UTF-8''encoded
	if idx := strings.Index(strings.ToLower(header), "filename*"); idx >= 0 {
		rest := header[idx:]
		if eqIdx := strings.Index(rest, "UTF-8''"); eqIdx >= 0 {
			encoded := rest[eqIdx+len("UTF-8''"):]
			encoded = strings.SplitN(encoded, ";", 2)[0]
			encoded = strings.TrimSpace(encoded)
			if decoded, err := url.PathUnescape(encoded); err == nil && decoded != "" {
				return decoded
			}
		}
	}

	// filename="quoted"
	if idx := strings.Index(strings.ToLower(header), "filename="); idx >= 0 {
		rest := header[idx+len("filename="):]
		rest = strings.TrimSpace(rest)
		if strings.HasPrefix(rest, "\"") {
			end := strings.Index(rest[1:], "\"")
			if end >= 0 {
				return rest[1 : end+1]
			}
		}
		// unquoted
		return strings.SplitN(strings.TrimSpace(rest), ";", 2)[0]
	}

	return ""
}

func SanitizeFilename(name string) string {
	// Strip leading path
	if idx := strings.LastIndexAny(name, "/\\"); idx >= 0 {
		name = name[idx+1:]
	}
	name = unsafeFilenameRE.ReplaceAllString(name, "_")
	if len(name) > 255 {
		name = name[:255]
	}
	if name == "" {
		return "download"
	}
	return name
}

func mimeToExtension(mimeType string) string {
	base := strings.SplitN(mimeType, ";", 2)[0]
	base = strings.TrimSpace(strings.ToLower(base))
	if ext, ok := mimeToExt[base]; ok {
		return ext
	}
	return ""
}

func resolveDestPath(filename string, opts DownloadOpts, contentType string) (string, error) {
	if opts.Output != "" {
		destPath, _ := filepath.Abs(opts.Output)
		outputExt := strings.ToLower(filepath.Ext(destPath))
		expectedExt := mimeToExtension(contentType)
		if expectedExt != "" && outputExt != "" && outputExt != expectedExt {
			fmt.Fprintf(os.Stderr, "Warning: output extension %q does not match Content-Type %q (expected %q)\n", outputExt, contentType, expectedExt)
		}
		if err := checkOverwrite(destPath, opts.Force); err != nil {
			return "", err
		}
		return destPath, nil
	}

	if opts.OutputDir != "" {
		if _, err := os.Stat(opts.OutputDir); os.IsNotExist(err) {
			return "", fmt.Errorf("Output directory does not exist: %q", opts.OutputDir)
		}
		destPath := filepath.Join(opts.OutputDir, filename)
		if err := checkOverwrite(destPath, opts.Force); err != nil {
			return "", err
		}
		return destPath, nil
	}

	cwd, _ := os.Getwd()
	destPath := filepath.Join(cwd, filename)
	if err := checkOverwrite(destPath, opts.Force); err != nil {
		return "", err
	}
	return destPath, nil
}

func checkOverwrite(path string, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("File already exists: %q. Use --force to overwrite.", path)
		}
	}
	return nil
}
