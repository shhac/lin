package download

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// eslint-disable-next-line no-control-regex
var unsafeFilenameRE = regexp.MustCompile(`[<>:"|?*\x00-\x1f]`)

var mimeToExt = map[string]string{
	"image/png":        ".png",
	"image/jpeg":       ".jpg",
	"image/gif":        ".gif",
	"image/webp":       ".webp",
	"image/svg+xml":    ".svg",
	"application/pdf":  ".pdf",
	"text/plain":       ".txt",
	"text/csv":         ".csv",
	"application/json": ".json",
	"application/zip":  ".zip",
	"video/mp4":        ".mp4",
	"audio/mpeg":       ".mp3",
}

// inferFilename derives a filename from the response: Content-Disposition
// first, falling back to the URL path with a MIME-derived extension.
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

// SanitizeFilename strips any leading path and replaces filesystem-unsafe
// characters with underscores. Caps at 255 chars.
func SanitizeFilename(name string) string {
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
