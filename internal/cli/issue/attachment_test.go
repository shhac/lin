package issue

import (
	"strings"
	"testing"
)

func TestValidateAttachmentFlags(t *testing.T) {
	cases := []struct {
		name      string
		flags     attachmentFlags
		wantErr   bool
		errSubstr string
	}{
		{name: "no flags ok", flags: attachmentFlags{}},
		{name: "github-pr alone ok", flags: attachmentFlags{githubPR: true}},
		{name: "slack alone ok", flags: attachmentFlags{slack: true}},
		{name: "slack with sync-thread ok", flags: attachmentFlags{slack: true, syncThread: true}},
		{
			name:      "two integration flags conflict",
			flags:     attachmentFlags{githubPR: true, slack: true},
			wantErr:   true,
			errSubstr: "conflicting flags",
		},
		{
			name:      "three integration flags conflict",
			flags:     attachmentFlags{githubPR: true, gitlabMR: true, discord: true},
			wantErr:   true,
			errSubstr: "conflicting flags",
		},
		{
			name:      "sync-thread without slack",
			flags:     attachmentFlags{syncThread: true},
			wantErr:   true,
			errSubstr: "--sync-thread requires --slack",
		},
		{
			name:      "sync-thread with non-slack integration flag",
			flags:     attachmentFlags{githubPR: true, syncThread: true},
			wantErr:   true,
			errSubstr: "--sync-thread requires --slack",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateAttachmentFlags(tc.flags)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tc.errSubstr) {
					t.Errorf("error = %q, want substring %q", err.Error(), tc.errSubstr)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestSelectLinkOp_DerivedArgsValidation exercises the URL-derived arg paths
// (gitlab MR, discord) returning a structured error before any network call.
func TestSelectLinkOp_DerivedArgsValidation(t *testing.T) {
	cases := []struct {
		name      string
		flags     attachmentFlags
		url       string
		wantOK    bool
		errSubstr string
	}{
		{
			name:      "gitlab-mr with bad URL",
			flags:     attachmentFlags{gitlabMR: true},
			url:       "https://gitlab.com/g/p/-/issues/1",
			errSubstr: "GitLab MR",
		},
		{
			name:   "gitlab-mr with valid URL",
			flags:  attachmentFlags{gitlabMR: true},
			url:    "https://gitlab.com/group/project/-/merge_requests/42",
			wantOK: true,
		},
		{
			name:      "discord with bad URL",
			flags:     attachmentFlags{discord: true},
			url:       "https://discord.com/invite/abc/123",
			errSubstr: "Discord",
		},
		{
			name:   "discord with valid URL",
			flags:  attachmentFlags{discord: true},
			url:    "https://discord.com/channels/1/2/3",
			wantOK: true,
		},
		{
			name:   "default URL link does not validate URL shape",
			flags:  attachmentFlags{},
			url:    "anything",
			wantOK: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			op, err := selectLinkOp(nil, nil, tc.flags, "issue", tc.url, nil)
			if tc.wantOK {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if op == nil {
					t.Fatal("expected non-nil op")
				}
				return
			}
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tc.errSubstr) {
				t.Errorf("error = %q, want substring %q", err.Error(), tc.errSubstr)
			}
		})
	}
}

func TestParseGitLabMRURL(t *testing.T) {
	cases := []struct {
		name        string
		url         string
		wantPath    string
		wantNum     int
		wantErr     bool
	}{
		{
			name:     "simple project",
			url:      "https://gitlab.com/group/project/-/merge_requests/42",
			wantPath: "group/project",
			wantNum:  42,
		},
		{
			name:     "nested subgroups",
			url:      "https://gitlab.com/group/sub1/sub2/project/-/merge_requests/7",
			wantPath: "group/sub1/sub2/project",
			wantNum:  7,
		},
		{
			name:     "trailing path segments",
			url:      "https://gitlab.com/group/project/-/merge_requests/3/diffs",
			wantPath: "group/project",
			wantNum:  3,
		},
		{
			name:    "missing merge_requests segment",
			url:     "https://gitlab.com/group/project/-/issues/3",
			wantErr: true,
		},
		{
			name:    "non-numeric MR number",
			url:     "https://gitlab.com/group/project/-/merge_requests/foo",
			wantErr: true,
		},
		{
			name:    "empty",
			url:     "",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path, num, err := parseGitLabMRURL(tc.url)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got path=%q num=%d", path, num)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if path != tc.wantPath {
				t.Errorf("path = %q, want %q", path, tc.wantPath)
			}
			if num != tc.wantNum {
				t.Errorf("num = %d, want %d", num, tc.wantNum)
			}
		})
	}
}

func TestParseDiscordMessageURL(t *testing.T) {
	cases := []struct {
		name        string
		url         string
		wantChannel string
		wantMessage string
		wantErr     bool
	}{
		{
			name:        "standard",
			url:         "https://discord.com/channels/111111111111111111/222222222222222222/333333333333333333",
			wantChannel: "222222222222222222",
			wantMessage: "333333333333333333",
		},
		{
			name:        "ptb subdomain",
			url:         "https://ptb.discord.com/channels/111/222/333",
			wantChannel: "222",
			wantMessage: "333",
		},
		{
			name:    "too few segments",
			url:     "https://discord.com/channels/111/222",
			wantErr: true,
		},
		{
			name:    "wrong root",
			url:     "https://discord.com/invite/abc/123/456",
			wantErr: true,
		},
		{
			name:    "empty",
			url:     "",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			channel, message, err := parseDiscordMessageURL(tc.url)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got channel=%q message=%q", channel, message)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if channel != tc.wantChannel {
				t.Errorf("channel = %q, want %q", channel, tc.wantChannel)
			}
			if message != tc.wantMessage {
				t.Errorf("message = %q, want %q", message, tc.wantMessage)
			}
		})
	}
}
