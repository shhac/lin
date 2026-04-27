package issue

import (
	"github.com/spf13/cobra"

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
