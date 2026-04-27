package issue

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// parseGitLabMRURL extracts (projectPathWithNamespace, MR number) from a GitLab MR URL.
// Accepts forms like https://gitlab.com/group/sub/project/-/merge_requests/42[/...].
func parseGitLabMRURL(raw string) (string, int, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Path == "" {
		return "", 0, fmt.Errorf("invalid GitLab URL: %q", raw)
	}
	path := strings.Trim(u.Path, "/")
	const sep = "/-/merge_requests/"
	idx := strings.Index(path, sep)
	if idx < 0 {
		return "", 0, fmt.Errorf("not a GitLab MR URL (missing /-/merge_requests/): %q", raw)
	}
	projectPath := path[:idx]
	rest := strings.TrimPrefix(path[idx:], "/-/merge_requests/")
	numStr := rest
	if slash := strings.Index(rest, "/"); slash >= 0 {
		numStr = rest[:slash]
	}
	number, err := strconv.Atoi(numStr)
	if err != nil || number <= 0 {
		return "", 0, fmt.Errorf("invalid MR number in URL: %q", raw)
	}
	if projectPath == "" {
		return "", 0, fmt.Errorf("missing project path in GitLab MR URL: %q", raw)
	}
	return projectPath, number, nil
}

// parseDiscordMessageURL extracts (channelId, messageId) from a Discord message URL.
// Accepts forms like https://discord.com/channels/<guild>/<channel>/<message>.
func parseDiscordMessageURL(raw string) (string, string, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Path == "" {
		return "", "", fmt.Errorf("invalid Discord URL: %q", raw)
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 4 || parts[0] != "channels" {
		return "", "", fmt.Errorf("not a Discord message URL: %q", raw)
	}
	channelID := parts[2]
	messageID := parts[3]
	if channelID == "" || messageID == "" {
		return "", "", fmt.Errorf("missing channel or message ID in Discord URL: %q", raw)
	}
	return channelID, messageID, nil
}
