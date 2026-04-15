package keychain

import (
	"os/exec"
	"runtime"
	"strings"
)

const service = "app.paulie.lin"

// Get retrieves a keychain entry. macOS only; returns "" on other platforms.
func Get(account string) (string, error) {
	if runtime.GOOS != "darwin" {
		return "", nil
	}
	out, err := exec.Command("security", "find-generic-password",
		"-s", service, "-a", account, "-w").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Store saves a keychain entry. macOS only; no-op on other platforms.
func Store(account, password string) error {
	if runtime.GOOS != "darwin" {
		return nil
	}
	return exec.Command("security", "add-generic-password",
		"-s", service, "-a", account, "-w", password, "-U").Run()
}

// Delete removes a keychain entry. macOS only; no-op on other platforms.
func Delete(account string) error {
	if runtime.GOOS != "darwin" {
		return nil
	}
	return exec.Command("security", "delete-generic-password",
		"-s", service, "-a", account).Run()
}
