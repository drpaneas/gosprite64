//go:build mage
// +build mage

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Base directory for our toolchain installation.
var baseDir = filepath.Join(os.Getenv("HOME"), "toolchains", "nintendo64")

// Directories for the custom Go and GOPATH.
var goDir = filepath.Join(baseDir, "go")
var gopathDir = filepath.Join(baseDir, "gopath")

// expectedVersionFingerprint is the substring we expect in our custom Go version.
const expectedVersionFingerprint = "af62b1cff2"

// runCommand executes a command and streams its output.
func runCommand(name string, args ...string) error {
	fmt.Printf("Running: %s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runCommandInDir executes a command in a given directory.
func runCommandInDir(dir, name string, args ...string) error {
	fmt.Printf("Running in %s: %s %s\n", dir, name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SetCustomEnv programmatically sets the environment variables needed for your custom Go.
func SetCustomEnv() {
	os.Setenv("GOROOT", goDir)
	os.Setenv("GOPATH", gopathDir)
	gobin := filepath.Join(gopathDir, "bin")
	os.Setenv("GOBIN", gobin)

	// Prepend our custom directories to PATH.
	path := os.Getenv("PATH")
	newPath := gobin + string(os.PathListSeparator) + filepath.Join(goDir, "bin") + string(os.PathListSeparator) + path
	os.Setenv("PATH", newPath)

	fmt.Println("Custom environment variables set:")
	fmt.Printf("  GOROOT=%s\n", goDir)
	fmt.Printf("  GOPATH=%s\n", gopathDir)
	fmt.Printf("  GOBIN=%s\n", gobin)
	fmt.Printf("  PATH=%s\n", os.Getenv("PATH"))
}

// VerifyCustomGo ensures that "go version" prints the expected custom version fingerprint.
func VerifyCustomGo() error {
	// Set our custom environment.
	SetCustomEnv()

	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return fmt.Errorf("failed to run go version: %w", err)
	}
	versionOutput := string(out)
	fmt.Printf("Detected Go version: %s\n", versionOutput)
	if !strings.Contains(versionOutput, expectedVersionFingerprint) {
		return fmt.Errorf("incorrect Go environment: got '%s', expected version containing '%s'",
			versionOutput, expectedVersionFingerprint)
	}
	return nil
}

// ConfigureDirenvHook detects the current shell and adds the appropriate direnv hook if not present.
func ConfigureDirenvHook() error {
	shell := os.Getenv("SHELL")
	var configFile, hookCmd string

	if strings.Contains(shell, "bash") {
		configFile = filepath.Join(os.Getenv("HOME"), ".bashrc")
		hookCmd = `eval "$(direnv hook bash)"`
	} else if strings.Contains(shell, "zsh") {
		configFile = filepath.Join(os.Getenv("HOME"), ".zshrc")
		hookCmd = `eval "$(direnv hook zsh)"`
	} else if strings.Contains(shell, "fish") {
		configFile = filepath.Join(os.Getenv("HOME"), ".config", "fish", "config.fish")
		hookCmd = `eval (direnv hook fish)`
	} else {
		fmt.Printf("Your shell (%s) is not explicitly supported. Please add the direnv hook manually.\n", shell)
		return nil
	}

	// Check if the hook is already present.
	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", configFile, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == hookCmd {
			fmt.Printf("Direnv hook already present in %s\n", configFile)
			return nil
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Append the hook if not present.
	f, err := os.OpenFile(configFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s for appending: %w", configFile, err)
	}
	defer f.Close()

	fmt.Printf("Appending direnv hook to %s...\n", configFile)
	if _, err := f.WriteString("\n# Load direnv\n" + hookCmd + "\n"); err != nil {
		return fmt.Errorf("failed to append to %s: %w", configFile, err)
	}
	fmt.Printf("Direnv hook added to %s. Please reload your shell.\n", configFile)
	return nil
}

// InstallDirenv installs direnv using the appropriate package manager and then configures the hook.
func InstallDirenv() error {
	fmt.Println("Installing direnv...")
	osType := runtime.GOOS

	if osType == "linux" {
		if _, err := exec.LookPath("apt-get"); err == nil {
			fmt.Println("Debian/Ubuntu detected. Installing direnv via apt-get...")
			if err := runCommand("sudo", "apt-get", "update"); err != nil {
				return err
			}
			if err := runCommand("sudo", "apt-get", "install", "-y", "direnv"); err != nil {
				return err
			}
		} else if _, err := exec.LookPath("dnf"); err == nil {
			fmt.Println("Fedora detected. Installing direnv via dnf...")
			if err := runCommand("sudo", "dnf", "install", "-y", "direnv"); err != nil {
				return err
			}
		} else if _, err := exec.LookPath("pacman"); err == nil {
			fmt.Println("Arch Linux detected. Installing direnv via pacman...")
			if err := runCommand("sudo", "pacman", "-Syu", "direnv"); err != nil {
				return err
			}
		} else {
			fmt.Println("No supported package manager found. Please install direnv manually.")
		}
	} else if osType == "darwin" {
		if _, err := exec.LookPath("brew"); err == nil {
			fmt.Println("macOS detected. Installing direnv via Homebrew...")
			if err := runCommand("brew", "install", "direnv"); err != nil {
				return err
			}
		} else {
			fmt.Println("Homebrew not installed. Please install Homebrew or install direnv manually.")
		}
	} else {
		fmt.Printf("Unsupported OS: %s. Please install direnv manually.\n", osType)
		return nil
	}

	return ConfigureDirenvHook()
}

// AutoDirenvAllow automatically runs "direnv allow" in the base directory.
func AutoDirenvAllow() error {
	fmt.Printf("Automatically allowing direnv in %s...\n", baseDir)
	return runCommandInDir(baseDir, "direnv", "allow")
}

// SetupToolchain ensures that the base directory exists.
func SetupToolchain() error {
	fmt.Printf("Ensuring base directory %s exists...\n", baseDir)
	return os.MkdirAll(baseDir, 0755)
}

// BuildCustomGo clones the custom Go repository (mips branch) and builds it.
// If the custom Go binary already exists, it skips the build.
func BuildCustomGo() error {
	// Check if the custom Go binary exists.
	goBinaryPath := filepath.Join(goDir, "bin", "go")
	if _, err := os.Stat(goBinaryPath); err == nil {
		fmt.Println("Custom Go binary already built, skipping rebuild.")
		return nil
	}

	// Clone repository if it doesn't exist.
	if _, err := os.Stat(goDir); os.IsNotExist(err) {
		fmt.Println("Cloning custom Go repository...")
		if err := runCommand("git", "clone", "https://github.com/clktmr/go", "-b", "mips", goDir); err != nil {
			return err
		}
	} else {
		fmt.Println("Custom Go repository already cloned.")
	}

	// Ensure we're on the mips branch.
	cmd := exec.Command("git", "checkout", "mips")
	cmd.Dir = goDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Build the custom Embedded Go.
	srcDir := filepath.Join(goDir, "src")
	fmt.Println("Building custom Embedded Go...")
	if err := runCommandInDir(srcDir, "bash", "-c", "./make.bash"); err != nil {
		return err
	}
	fmt.Println("Custom Go built successfully.")
	return nil
}

// SetupGopath creates the GOPATH directory structure.
func SetupGopath() error {
	dirs := []string{
		gopathDir,
		filepath.Join(gopathDir, "src"),
		filepath.Join(gopathDir, "pkg"),
		filepath.Join(gopathDir, "bin"),
	}
	for _, d := range dirs {
		fmt.Printf("Ensuring directory %s exists...\n", d)
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	fmt.Println("GOPATH structure created.")
	return nil
}

// ConfigureEnv writes an .envrc file in the base directory to set environment variables.
func ConfigureEnv() error {
	envrcPath := filepath.Join(baseDir, ".envrc")
	content := fmt.Sprintf(`export GOROOT=%s
export GOPATH=%s
export GOBIN=%s
export PATH=$GOBIN:%s/bin:$PATH
`, goDir, gopathDir, filepath.Join(gopathDir, "bin"), goDir)
	fmt.Printf("Creating .envrc at %s\n", envrcPath)
	if err := os.WriteFile(envrcPath, []byte(content), 0644); err != nil {
		return err
	}
	fmt.Println(".envrc file created successfully.")
	return nil
}

// EditGoMod uses Go's own module editing command to update go.mod.
// It replaces "github.com/drpaneas/gosprite64" with the local path.
func EditGoMod() error {
	localModulePath := filepath.Join(gopathDir, "src", "gosprite64")
	fmt.Println("Editing go.mod to replace module path with local module path...")
	return runCommand("go", "mod", "edit", "-replace=github.com/drpaneas/gosprite64="+localModulePath)
}

// InstallEmgo installs the emgo tool—but only after verifying we are using the correct Go environment.
func InstallEmgo() error {
	if err := VerifyCustomGo(); err != nil {
		return fmt.Errorf("cannot install emgo: %w", err)
	}
	fmt.Println("Installing emgo tool...")
	if err := runCommand("go", "install", "-v", "github.com/embeddedgo/tools/emgo@latest"); err != nil {
		return err
	}
	fmt.Println("emgo installed successfully.")
	return nil
}

// CloneGosprite64 clones the gosprite64 repository into GOPATH/src.
func CloneGosprite64() error {
	projectDir := filepath.Join(gopathDir, "src", "gosprite64")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		fmt.Println("Cloning gosprite64 repository...")
		if err := runCommand("git", "clone", "https://github.com/drpaneas/gosprite64", projectDir); err != nil {
			return err
		}
	} else {
		fmt.Println("gosprite64 repository already cloned.")
	}
	return nil
}

// BuildExample builds the clearscreen example using emgo.
// It verifies that the custom Go environment is active.
// When running in GitHub Actions, it also tidies modules and updates module paths.
func BuildExample() error {
	if err := VerifyCustomGo(); err != nil {
		return fmt.Errorf("cannot build example: %w", err)
	}
	exampleDir := filepath.Join(gopathDir, "src", "gosprite64", "examples", "clearscreen")
	// For GitHub Actions, tidy and update module paths using go mod edit instead of sed.
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		fmt.Println("Detected GitHub Actions environment; running emgo mod tidy...")
		if err := runCommandInDir(exampleDir, "emgo", "mod", "tidy"); err != nil {
			return err
		}
		fmt.Println("Editing go.mod in the example to use the local module path...")
		if err := EditGoMod(); err != nil {
			fmt.Println("Warning: Failed to edit go.mod:", err)
		}
	}
	fmt.Println("Building clearscreen example...")
	if err := runCommandInDir(exampleDir, "emgo", "build"); err != nil {
		return err
	}
	fmt.Println("Clearscreen example built successfully.")
	return nil
}

// Setup runs the setup tasks for the toolchain, excluding cloning GoSprite64.
func Setup() error {
	if err := SetupToolchain(); err != nil {
		return err
	}
	if err := BuildCustomGo(); err != nil {
		return err
	}
	if err := SetupGopath(); err != nil {
		return err
	}
	if err := ConfigureEnv(); err != nil {
		return err
	}
	if err := InstallDirenv(); err != nil {
		return err
	}
	if err := AutoDirenvAllow(); err != nil {
		return err
	}
	if err := InstallEmgo(); err != nil {
		return err
	}
	fmt.Println("Toolchain setup complete.")
	return nil
}

// Test is a separate target that clones the GoSprite64 repository, updates go.mod, and builds the example.
func Test() error {
	if err := CloneGosprite64(); err != nil {
		return err
	}
	if err := EditGoMod(); err != nil {
		return err
	}
	if err := BuildExample(); err != nil {
		return err
	}
	fmt.Println("Test target complete: GoSprite64 cloned and example built.")
	return nil
}
