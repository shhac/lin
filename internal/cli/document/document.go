package document

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/output"
)

func Register(parent *cobra.Command) {
	document := &cobra.Command{
		Use:   "document",
		Short: "Document operations",
	}
	parent.AddCommand(document)

	registerSearch(document)
	registerList(document)
	registerGet(document)
	registerNew(document)
	registerUpdate(document)
	registerHistory(document)
	registerUsage(document)
	output.HandleUnknownCommand(document, "To view a document: lin document get <id>")
}

