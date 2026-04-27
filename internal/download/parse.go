package download

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/shhac/lin/internal/filters"
)

type ParsedFileURL struct {
	URL      string
	OrgID    string
	Segments []string
}

// ParseFileURL accepts an uploads.linear.app URL (or path-only / host-relative
// shortcuts) and returns a normalized form. defaultOrgID supplies the org
// segment when the input has only the file UUIDs.
func ParseFileURL(input string, defaultOrgID string) (ParsedFileURL, error) {
	var pathname string

	if strings.HasPrefix(input, "http://") {
		return ParsedFileURL{}, fmt.Errorf("refusing http:// URL — only https:// is allowed for file downloads")
	}

	if strings.HasPrefix(input, "https://") {
		u, err := url.Parse(input)
		if err != nil {
			return ParsedFileURL{}, fmt.Errorf("invalid URL: %v", err)
		}
		if u.Hostname() != uploadHost {
			return ParsedFileURL{}, fmt.Errorf("invalid host: %q, only %s URLs are supported", u.Hostname(), uploadHost)
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
		return ParsedFileURL{}, fmt.Errorf("cannot parse file URL: %q, expected 1-3 UUID path segments", input)
	}

	for _, seg := range segments {
		if !filters.IsUUID(seg) {
			return ParsedFileURL{}, fmt.Errorf("invalid UUID segment: %q, expected format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", seg)
		}
	}

	var orgID string
	var fileSegments []string

	if len(segments) == 3 {
		orgID = segments[0]
		fileSegments = segments[1:]
	} else {
		if defaultOrgID == "" {
			return ParsedFileURL{}, fmt.Errorf("cannot infer organization ID, provide a full URL with org segment, or authenticate first")
		}
		orgID = defaultOrgID
		fileSegments = segments
	}

	allSegments := append([]string{orgID}, fileSegments...)
	fullURL := fmt.Sprintf("https://%s/%s", uploadHost, strings.Join(allSegments, "/"))

	return ParsedFileURL{URL: fullURL, OrgID: orgID, Segments: allSegments}, nil
}
