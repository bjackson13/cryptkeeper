package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// EditorResolver returns the editor to use based on user preferrence or default
type EditorResolver func() string

/**
* Get the users prefered editor either from passed flag or default
* TODO: Add ability to set editor in config file
* Implements EditorResolver type
 */
func getPreferredEditor() string {
	// Supported Editors - I have confirmed all of these
	supportedEditors := []string{"vim", "vi", "code", "vsc", "nano"}

	if editor == "" {
		return DefaultEditor
	}

	// determine if requested editor is supported or not
	isSupported := func() bool {
		for _, val := range supportedEditors {
			if editor == val {
				return true
			}
		}
		return false
	}()

	if !isSupported {
		return DefaultEditor
	}

	return editor
}

/**
* Appends needed arguments to certain editors that are required for proper functionality
 */
func appendEditorArgs(command string, tempFilename string) []string {
	args := []string{tempFilename}

	if strings.Contains(command, "code") || strings.Contains(command, "vsc") {
		args = append([]string{"--wait"}, args...)
	}

	return args
}

/**
* Open users preffered text editor to capture log
 */
func openEditor(filename string, resolver EditorResolver) error {
	executable, err := exec.LookPath(resolver())
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(executable, appendEditorArgs(executable, filename)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

/**
* Create a temporary file in a temp directory and open it in the users preffered editor.
* Delete the temp file once we read in the contents of it
 */
func GetEditorInput() ([]byte, error) {
	file, err := os.CreateTemp(os.TempDir(), "*")
	if err != nil {
		log.Fatal(err)
	}

	tempFile := file.Name()
	defer os.Remove(tempFile)

	if err = file.Close(); err != nil {
		return []byte{}, err
	}

	if openEditor(tempFile, getPreferredEditor); err != nil {
		return []byte{}, err
	}

	bytes, err := os.ReadFile(tempFile)
	if err != nil {
		log.Fatal(err)
	}

	return bytes, nil
}

/**
* Displays decrypted plain text journal in user prefered editor
 */
func DisplayOutputInEditor(plaintextJournal []byte) {
	file, err := os.CreateTemp(os.TempDir(), "*")
	if err != nil {
		log.Fatal(err)
	}

	tempFile := file.Name()
	defer os.Remove(tempFile)

	if err = file.Close(); err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(tempFile, plaintextJournal, 0700)
	if err != nil {
		log.Fatal(err)
	}

	if openEditor(tempFile, getPreferredEditor); err != nil {
		log.Fatal(err)
	}
}
