package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
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

func registerAttachmentAdd(parent *cobra.Command) {
	var (
		url      string
		title    string
		subtitle string
	)

	cmd := &cobra.Command{
		Use:   "add <issue-id>",
		Short: "Add a URL attachment to an issue",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			input := linear.AttachmentCreateInput{
				IssueId: args[0],
				Url:     url,
				Title:   title,
			}
			if subtitle != "" {
				input.Subtitle = &subtitle
			}

			resp, err := linear.AttachmentCreate(ctx, client, input)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			a := resp.AttachmentCreate.Attachment
			output.PrintJSON(map[string]any{
				"id":      a.Id,
				"title":   a.Title,
				"url":     a.Url,
				"created": resp.AttachmentCreate.Success,
			})
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "URL to attach")
	_ = cmd.MarkFlagRequired("url")
	cmd.Flags().StringVar(&title, "title", "", "Attachment title")
	_ = cmd.MarkFlagRequired("title")
	cmd.Flags().StringVar(&subtitle, "subtitle", "", "Attachment subtitle")
	parent.AddCommand(cmd)
}

func registerAttachmentRemove(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "remove <attachment-id>",
		Short: "Remove an attachment",
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
