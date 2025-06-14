//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// versionRE is used to extract Go version from version output
var versionRE = regexp.MustCompile(`go([0-9]+\.[0-9]+(?:\.[0-9]+)?)`)

// runCommand executes a command and streams its output.
func runCommand(name string, args ...string) error {
	fmt.Printf("Running: %s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runCommandInDir executes a command in the given directory.
func runCommandInDir(dir, name string, args ...string) error {
	fmt.Printf("Running in %s: %s %s\n", dir, name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Setup is the primary target that implements the new installation steps.
func Setup() error {
	// Get GOPATH using the helper function - use this once for the entire function
	gopath, err := getGOPATH()
	if err != nil {
		return err
	}
	gobin := os.Getenv("GOBIN")
	if gobin == "" {
		gobin = filepath.Join(gopath, "bin")
	}
	fmt.Printf("Using default GOPATH=%s, GOBIN=%s\n", gopath, gobin)

	// 1. Install gotip-embedded
	if err := runCommand("go", "install", "github.com/clktmr/dl/gotip-embedded@latest"); err != nil {
		return fmt.Errorf("failed to install gotip-embedded: %w", err)
	}
	if _, err := exec.LookPath("gotip-embedded"); err != nil {
		return fmt.Errorf("gotip-embedded not found in PATH after installation: %w", err)
	}

	// 2. Run gotip-embedded download
	if err := runCommand("gotip-embedded", "download"); err != nil {
		return fmt.Errorf("failed to download embedded Go toolchain: %w", err)
	}

	// 3. Create a shortcut for the new Go binary (using absolute path for symlink)
	// First, get the version string using native Go
	verCmd := exec.Command("gotip-embedded", "version")
	verOut, err := verCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to determine Go version: %w", err)
	}

	versionStr := strings.TrimSpace(string(verOut))

	// Extract version using the package-level regex
	match := versionRE.FindStringSubmatch(versionStr)
	if len(match) < 2 {
		return fmt.Errorf("unexpected version output: %q", versionStr)
	}
	goVersion := match[1] // e.g. "1.23" or "1.23.1"
	goVerName := "go" + goVersion

	// Ensure the bin directory exists before creating symlink
	binDir := filepath.Join(gopath, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create GOPATH bin directory: %w", err)
	}

	linkPath := filepath.Join(binDir, goVerName)
	if runtime.GOOS != "windows" {
		// Get absolute path to gotip-embedded
		binPath, err := exec.LookPath("gotip-embedded")
		if err != nil {
			return fmt.Errorf("failed to find gotip-embedded in PATH: %w", err)
		}

		// Remove existing symlink if it exists
		if fi, err := os.Lstat(linkPath); err == nil {
			if fi.Mode()&os.ModeSymlink != 0 {
				if err := os.Remove(linkPath); err != nil {
					fmt.Printf("Warning: Could not remove existing symlink at %s: %v\n", linkPath, err)
				} else {
					fmt.Printf("Removed existing symlink at %s\n", linkPath)
				}
			}
		}

		if err := os.Symlink(binPath, linkPath); err != nil && !os.IsExist(err) {
			return fmt.Errorf("failed to create symlink: %w", err)
		}
		fmt.Printf("Created symlink from %s to %s\n", linkPath, binPath)
	} else {
		// Note: Could use the following for Windows if running with elevated privileges:
		// exec.Command("cmd", "/c", "mklink", linkPath, binPath).Run()
		fmt.Println("Skipping symlink creation on Windows; please create a shortcut named", goVerName, "to gotip-embedded manually.")
	}

	// 4. Install mkrom
	if err := runCommand("go", "install", "github.com/clktmr/n64/tools/mkrom"); err != nil {
		return fmt.Errorf("failed to install mkrom: %w", err)
	}
	if _, err := exec.LookPath("mkrom"); err != nil {
		return fmt.Errorf("mkrom not found in PATH: %w", err)
	}

	// 5. Create .envrc file (correct name for direnv compatibility)
	// Get current working directory to ensure .envrc is created in the proper location
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to determine current directory: %w", err)
	}
	// Figure out the exact toolchain name
	toolchain, err := detectGoToolchain()
	if err != nil {
		return fmt.Errorf("cannot detect embedded Go version for .envrc: %w", err)
	}

	envrcPath := filepath.Join(projectDir, ".envrc")
	content := fmt.Sprintf(`export GOOS="noos"
export GOARCH="mips64"
export GOFLAGS="-tags=n64 '-ldflags=-M=0x00000000:8M -F=0x00000400:8M -stripfn=1'"
export GOTOOLCHAIN="%s"
`, toolchain)

	if err := os.WriteFile(envrcPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write .envrc file: %w", err)
	}
	fmt.Printf("Created .envrc file at %s (GOTOOLCHAIN=%s)\n", envrcPath, toolchain)

	// 6. Ensure direnv is installed (legacy installation)
	if err := InstallDirenv(); err != nil {
		return fmt.Errorf("failed to install direnv: %w", err)
	}

	// 7. Run direnv allow (skip on Windows)
	if runtime.GOOS != "windows" {
		if err := runCommandInDir(projectDir, "direnv", "allow"); err != nil {
			fmt.Println("Warning: Failed to run 'direnv allow'. You may need to run it manually.")
		} else {
			fmt.Println("Successfully ran 'direnv allow' for the .envrc file.")
		}
	}

	// Final instructions
	fmt.Println("\n========== Setup Complete ==========\n")
	fmt.Println("To build:      go build -o test.elf .")
	fmt.Println("To create rom: mkrom test.elf\n\n")
	return nil
}

// InstallDirenv installs direnv on supported OSes and configures the hook.
func InstallDirenv() error {
	var installErr error

	switch runtime.GOOS {
	case "windows":
		installErr = installDirenvWindows()
	case "linux":
		installErr = installDirenvLinux()
	case "darwin":
		installErr = installDirenvDarwin()
	default:
		fmt.Printf("Unsupported OS for direnv installation: %s. Please install manually.\n", runtime.GOOS)
	}

	if installErr != nil {
		return installErr
	}

	// Always call ConfigureDirenvHook after installation (it handles Windows-specific logic internally)
	return ConfigureDirenvHook()
}

func installDirenvWindows() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to determine user home directory: %w", err)
	}

	dest := filepath.Join(home, "direnv.exe")
	if _, err := os.Stat(dest); err == nil {
		fmt.Println("direnv.exe already exists at", dest)
		return nil
	}

	url := "https://github.com/direnv/direnv/releases/latest/download/direnv.windows-amd64.exe"
	fmt.Println("Downloading direnv from", url)
	rsp, err := exec.Command("powershell", "-Command", fmt.Sprintf("Invoke-WebRequest -Uri %s -OutFile %s", url, dest)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("download failed: %s: %w", string(rsp), err)
	}
	fmt.Println("direnv installed at", dest)
	return nil
}

func installDirenvLinux() error {
	if _, err := exec.LookPath("apt-get"); err == nil {
		fmt.Println("Installing direnv via apt-get")
		if err := runCommand("sudo", "apt-get", "update"); err != nil {
			return err
		}
		return runCommand("sudo", "apt-get", "install", "-y", "direnv")
	}
	if _, err := exec.LookPath("dnf"); err == nil {
		fmt.Println("Installing direnv via dnf")
		return runCommand("sudo", "dnf", "install", "-y", "direnv")
	}
	fmt.Println("No supported package manager found. Please install direnv manually.")
	return nil
}

func installDirenvDarwin() error {
	if _, err := exec.LookPath("brew"); err == nil {
		fmt.Println("Installing direnv via Homebrew")
		return runCommand("brew", "install", "direnv")
	}
	fmt.Println("Homebrew not found. Please install Homebrew or direnv manually.")
	return nil
}

// ConfigureDirenvHook appends the direnv hook to the user's shell config.
func ConfigureDirenvHook() error {
	if runtime.GOOS == "windows" {
		fmt.Println("Skipping direnv hook configuration on Windows.")
		fmt.Println("Note: To enable direnv in PowerShell, add the following to your profile:")
		fmt.Println("  - For bash: Add 'eval \"$(direnv hook bash)\"' to your bash profile")
		fmt.Println("  - For PowerShell: Add 'Invoke-Expression \"$(direnv hook powershell)\"' to your PowerShell profile")
		return nil
	}

	shell := os.Getenv("SHELL")
	var config, hook string

	// Determine shell type and corresponding config file and hook
	if strings.Contains(shell, "bash") {
		config = filepath.Join(os.Getenv("HOME"), ".bashrc")
		hook = `eval "$(direnv hook bash)"`
	} else if strings.Contains(shell, "zsh") {
		config = filepath.Join(os.Getenv("HOME"), ".zshrc")
		hook = `eval "$(direnv hook zsh)"`
	} else if strings.Contains(shell, "fish") || filepath.Base(shell) == "fish" {
		// Support both when SHELL contains "fish" and when it's exactly "/usr/bin/fish"
		config = filepath.Join(os.Getenv("HOME"), ".config", "fish", "config.fish")
		hook = `eval (direnv hook fish)`
	} else {
		fmt.Printf("Shell not supported for auto-hook: %s. Please add hook manually.\n", shell)
		fmt.Println("For bash: Add 'eval \"$(direnv hook bash)\"' to your .bashrc")
		fmt.Println("For zsh: Add 'eval \"$(direnv hook zsh)\"' to your .zshrc")
		fmt.Println("For fish: Add 'eval (direnv hook fish)' to your config.fish")
		return nil
	}

	// Check if hook already exists to avoid duplicates
	if _, err := os.Stat(config); err == nil {
		content, err := os.ReadFile(config)
		if err == nil && strings.Contains(string(content), hook) {
			fmt.Printf("direnv hook already present in %s\n", config)
			return nil
		}
	}

	// Ensure the parent directory for the config file exists (important for fish shell)
	configDir := filepath.Dir(config)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	f, err := os.OpenFile(config, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open %s: %w", config, err)
	}
	defer f.Close()

	if _, err := f.WriteString("\n# direnv hook\n" + hook + "\n"); err != nil {
		return fmt.Errorf("failed to write hook: %w", err)
	}

	fmt.Println("Added direnv hook to", config)
	return nil
}

// Test builds the "clearscreen" example and produces a .z64 ROM.
// detectGoToolchain returns the Go toolchain version as a string (e.g., "go1.22" or "go1.22.3")
func detectGoToolchain() (string, error) {
	verCmd := exec.Command("gotip-embedded", "version")
	out, err := verCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to determine Go version: %w", err)
	}
	s := strings.TrimSpace(string(out))
	m := versionRE.FindStringSubmatch(s)
	if len(m) < 2 {
		return "", fmt.Errorf("unexpected version output: %q", s)
	}
	return "go" + m[1], nil
}

// getGOPATH returns the GOPATH with fallback logic
func getGOPATH() (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		return gopath, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine user home directory: %w", err)
	}

	// Allow GOPATH override via environment variable
	gopath = os.Getenv("MAGE_GOPATH")
	if gopath == "" {
		gopath = filepath.Join(home, "gocode")
	}

	return gopath, nil
}
