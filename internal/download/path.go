package download

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// resolveDestPath returns the absolute path to write the downloaded file,
// honoring --output, --output-dir, or defaulting to cwd. It also warns when
// --output's extension doesn't match the Content-Type.
func resolveDestPath(filename string, opts DownloadOpts, contentType string) (string, error) {
	if opts.Output != "" {
		destPath, _ := filepath.Abs(opts.Output)
		outputExt := strings.ToLower(filepath.Ext(destPath))
		expectedExt := mimeToExtension(contentType)
		if expectedExt != "" && outputExt != "" && outputExt != expectedExt {
			fmt.Fprintf(os.Stderr, "Warning: output extension %q does not match Content-Type %q (expected %q)\n", outputExt, contentType, expectedExt)
		}
		if err := checkOverwrite(destPath, opts.Force); err != nil {
			return "", err
		}
		return destPath, nil
	}

	if opts.OutputDir != "" {
		if _, err := os.Stat(opts.OutputDir); os.IsNotExist(err) {
			return "", fmt.Errorf("output directory does not exist: %q", opts.OutputDir)
		}
		destPath := filepath.Join(opts.OutputDir, filename)
		if err := checkOverwrite(destPath, opts.Force); err != nil {
			return "", err
		}
		return destPath, nil
	}

	cwd, _ := os.Getwd()
	destPath := filepath.Join(cwd, filename)
	if err := checkOverwrite(destPath, opts.Force); err != nil {
		return "", err
	}
	return destPath, nil
}

func checkOverwrite(path string, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("file already exists: %q, use --force to overwrite", path)
		}
	}
	return nil
}
