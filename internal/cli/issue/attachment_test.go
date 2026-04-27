package issue

import "testing"

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
