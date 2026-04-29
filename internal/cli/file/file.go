package file

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	"github.com/shhac/lin/internal/credential"
	dl "github.com/shhac/lin/internal/download"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/upload"
)

func Register(parent *cobra.Command) {
	file := &cobra.Command{
		Use:   "file",
		Short: "File operations",
	}
	output.HandleUnknownCommand(file, "Upload files: lin file upload <paths...>")

	registerUpload(file)
	registerDownload(file)
	shared.RegisterUsage(file, "file", usageText)

	parent.AddCommand(file)
}

func registerUpload(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "upload <paths...>",
		Short: "Upload files to Linear",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			results, err := upload.UploadFiles(client, args)
			if err != nil {
				output.PrintError(err.Error())
			}
			output.PrintJSON(results)
		},
	}
	parent.AddCommand(cmd)
}

func registerDownload(parent *cobra.Command) {
	var (
		flagOutput    string
		flagOutputDir string
		flagStdout    bool
		flagForce     bool
	)

	cmd := &cobra.Command{
		Use:   "download <url-or-path>",
		Short: "Download a file from Linear",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			apiKey := credential.Resolve()
			if apiKey == "" {
				output.PrintError("Not authenticated. Run: lin auth login")
			}

			orgID, err := dl.GetOrgID(client)
			if err != nil {
				output.PrintError(err.Error())
			}

			parsed, err := dl.ParseFileURL(args[0], orgID)
			if err != nil {
				output.PrintError(err.Error())
			}

			result, err := dl.DownloadFile(parsed.URL, dl.DownloadOpts{
				Output:    flagOutput,
				OutputDir: flagOutputDir,
				Stdout:    flagStdout,
				Force:     flagForce,
				APIKey:    apiKey,
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			if flagStdout {
				// Metadata to stderr when content goes to stdout
				enc := json.NewEncoder(os.Stderr)
				enc.SetEscapeHTML(false)
				_ = enc.Encode(result)
			} else {
				output.PrintJSON(result)
			}
		},
	}
	cmd.Flags().StringVar(&flagOutput, "output", "", "Save to specific file path")
	cmd.Flags().StringVar(&flagOutputDir, "output-dir", "", "Save to directory (default: current directory)")
	cmd.Flags().BoolVar(&flagStdout, "stdout", false, "Write file content to stdout")
	cmd.Flags().BoolVar(&flagForce, "force", false, "Overwrite existing files")
	cmd.MarkFlagsMutuallyExclusive("output", "output-dir", "stdout")

	parent.AddCommand(cmd)
}

const usageText = `lin file — File operations (upload, download)

UPLOAD:
  file upload <paths...>                  Upload one or more files to Linear

DOWNLOAD:
  file download <url-or-path>             Download a file from Linear
    --output <path>                       Save to specific file path
    --output-dir <dir>                    Save to directory (default: CWD)
    --stdout                              Write file content to stdout
    --force                               Overwrite existing files

URL FORMATS:
  Full URL      https://uploads.linear.app/<org>/<uuid>/<uuid>
  Host-relative uploads.linear.app/<org>/<uuid>/<uuid>
  Path only     <org>/<uuid>/<uuid>
  Short path    <uuid>/<uuid>   (org inferred from auth)
  Single UUID   <uuid>          (org inferred from auth)

OUTPUT:
  upload  → [{ filename, assetUrl, contentType }]
  download → { filename, path, size, contentType }
  download --stdout → binary to stdout, metadata JSON to stderr

NOTES:
  --output, --output-dir, and --stdout are mutually exclusive.
  Without --force, download refuses to overwrite existing files.`
