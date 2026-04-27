package issue

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
)

func registerAttachment(parent *cobra.Command) {
	attachment := &cobra.Command{
		Use:   "attachment",
		Short: "Attachment operations",
	}
	parent.AddCommand(attachment)

	registerAttachmentList(attachment)
	registerAttachmentAdd(attachment)
	registerAttachmentRemove(attachment)

	output.HandleUnknownCommand(attachment, "Run 'lin issue usage' for available attachment subcommands")
}

func registerAttachmentList(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "list <issue-id>",
		Short: "List attachments on an issue",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueAttachments(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]any, len(resp.Issue.Attachments.Nodes))
			for i, a := range resp.Issue.Attachments.Nodes {
				items[i] = map[string]any{
					"id":         a.Id,
					"title":      a.Title,
					"url":        a.Url,
					"subtitle":   a.Subtitle,
					"sourceType": a.SourceType,
				}
			}

			output.PrintJSON(items)
		},
	})
}

// attachmentFlags collects the integration-flag state for `lin issue attachment add`.
type attachmentFlags struct {
	githubPR    bool
	githubIssue bool
	gitlabMR    bool
	slack       bool
	discord     bool
	syncThread  bool
}

// linkResult is the uniform shape returned by every attachment-link mutation.
type linkResult struct {
	success    bool
	id         string
	title      string
	url        string
	sourceType *string
}

func registerAttachmentAdd(parent *cobra.Command) {
	var (
		title string
		flags attachmentFlags
	)

	cmd := &cobra.Command{
		Use:   "add <issue-id> <url>",
		Short: "Link a URL to an issue (rich attachment when integration matches)",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			issueID := args[0]
			rawURL := args[1]

			if err := validateAttachmentFlags(flags); err != nil {
				output.WriteError(err)
			}

			var titlePtr *string
			if title != "" {
				titlePtr = ptr.To(title)
			}

			client := linear.GetClient()
			ctx := context.Background()

			op, opErr := selectLinkOp(ctx, client, flags, issueID, rawURL, titlePtr)
			if opErr != nil {
				output.WriteError(opErr)
			}

			result, err := op()
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{
				"created":    result.success,
				"id":         result.id,
				"title":      result.title,
				"url":        result.url,
				"sourceType": ptr.Deref(result.sourceType),
			})
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Override the attachment title")
	cmd.Flags().BoolVar(&flags.githubPR, "github-pr", false, "Link as a GitHub pull request (rich attachment with PR sync)")
	cmd.Flags().BoolVar(&flags.githubIssue, "github-issue", false, "Link as a GitHub issue (rich attachment with issue sync)")
	cmd.Flags().BoolVar(&flags.gitlabMR, "gitlab-mr", false, "Link as a GitLab merge request (rich attachment with MR sync)")
	cmd.Flags().BoolVar(&flags.slack, "slack", false, "Link as a Slack message")
	cmd.Flags().BoolVar(&flags.syncThread, "sync-thread", false, "Sync the Slack thread with the issue's comment thread (requires --slack)")
	cmd.Flags().BoolVar(&flags.discord, "discord", false, "Link as a Discord message")
	parent.AddCommand(cmd)
}

// validateAttachmentFlags returns nil if the flag combination is valid, or a
// structured error describing the conflict.
func validateAttachmentFlags(f attachmentFlags) *apierrors.APIError {
	var selected []string
	if f.githubPR {
		selected = append(selected, "--github-pr")
	}
	if f.githubIssue {
		selected = append(selected, "--github-issue")
	}
	if f.gitlabMR {
		selected = append(selected, "--gitlab-mr")
	}
	if f.slack {
		selected = append(selected, "--slack")
	}
	if f.discord {
		selected = append(selected, "--discord")
	}
	if len(selected) > 1 {
		return apierrors.Newf(apierrors.FixableByAgent,
			"conflicting flags: %s", strings.Join(selected, ", ")).
			WithHint("specify at most one integration flag")
	}
	if f.syncThread && !f.slack {
		return apierrors.New("--sync-thread requires --slack", apierrors.FixableByAgent).
			WithHint("add --slack to sync the thread to a comment thread")
	}
	return nil
}

// selectLinkOp returns a closure that performs the appropriate Linear mutation
// based on the chosen flags. It validates URL-derived args eagerly so callers
// see a clean error rather than a confusing GraphQL response.
func selectLinkOp(
	ctx context.Context,
	client graphql.Client,
	f attachmentFlags,
	issueID, rawURL string,
	titlePtr *string,
) (func() (linkResult, error), *apierrors.APIError) {
	switch {
	case f.githubPR:
		return func() (linkResult, error) {
			r, err := linear.AttachmentLinkGitHubPR(ctx, client, issueID, rawURL, titlePtr)
			if err != nil {
				return linkResult{}, err
			}
			a := r.AttachmentLinkGitHubPR
			return linkResult{a.Success, a.Attachment.Id, a.Attachment.Title, a.Attachment.Url, a.Attachment.SourceType}, nil
		}, nil
	case f.githubIssue:
		return func() (linkResult, error) {
			r, err := linear.AttachmentLinkGitHubIssue(ctx, client, issueID, rawURL, titlePtr)
			if err != nil {
				return linkResult{}, err
			}
			a := r.AttachmentLinkGitHubIssue
			return linkResult{a.Success, a.Attachment.Id, a.Attachment.Title, a.Attachment.Url, a.Attachment.SourceType}, nil
		}, nil
	case f.gitlabMR:
		projectPath, number, err := parseGitLabMRURL(rawURL)
		if err != nil {
			return nil, apierrors.New(err.Error(), apierrors.FixableByAgent).
				WithHint("GitLab MR URLs look like https://gitlab.com/<group>/<project>/-/merge_requests/<n>")
		}
		return func() (linkResult, error) {
			r, err := linear.AttachmentLinkGitLabMR(ctx, client, issueID, rawURL, projectPath, float64(number), titlePtr)
			if err != nil {
				return linkResult{}, err
			}
			a := r.AttachmentLinkGitLabMR
			return linkResult{a.Success, a.Attachment.Id, a.Attachment.Title, a.Attachment.Url, a.Attachment.SourceType}, nil
		}, nil
	case f.slack:
		var syncPtr *bool
		if f.syncThread {
			syncPtr = ptr.To(true)
		}
		return func() (linkResult, error) {
			r, err := linear.AttachmentLinkSlack(ctx, client, issueID, rawURL, titlePtr, syncPtr)
			if err != nil {
				return linkResult{}, err
			}
			a := r.AttachmentLinkSlack
			return linkResult{a.Success, a.Attachment.Id, a.Attachment.Title, a.Attachment.Url, a.Attachment.SourceType}, nil
		}, nil
	case f.discord:
		channelID, messageID, err := parseDiscordMessageURL(rawURL)
		if err != nil {
			return nil, apierrors.New(err.Error(), apierrors.FixableByAgent).
				WithHint("Discord URLs look like https://discord.com/channels/<guild>/<channel>/<message>")
		}
		return func() (linkResult, error) {
			r, err := linear.AttachmentLinkDiscord(ctx, client, issueID, rawURL, channelID, messageID, titlePtr)
			if err != nil {
				return linkResult{}, err
			}
			a := r.AttachmentLinkDiscord
			return linkResult{a.Success, a.Attachment.Id, a.Attachment.Title, a.Attachment.Url, a.Attachment.SourceType}, nil
		}, nil
	default:
		return func() (linkResult, error) {
			r, err := linear.AttachmentLinkURL(ctx, client, issueID, rawURL, titlePtr)
			if err != nil {
				return linkResult{}, err
			}
			a := r.AttachmentLinkURL
			return linkResult{a.Success, a.Attachment.Id, a.Attachment.Title, a.Attachment.Url, a.Attachment.SourceType}, nil
		}, nil
	}
}

func registerAttachmentRemove(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "remove <attachment-id>",
		Short: "Remove an attachment (works for any source type)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.AttachmentDelete(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"deleted": resp.AttachmentDelete.Success})
		},
	})
}

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
