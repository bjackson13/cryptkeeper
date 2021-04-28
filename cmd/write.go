/*
Copyright © 2021 Brennan Jackson btj9560@rit.edu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

/**
* Default values
 */
const DefaultEditor = "vim"
const FileExtension = ".ck"

var DefaultFileName string = time.Now().Format("2006-01-02_1504")

/**
* Flag variables
 */
var fileName string
var passPhrase string
var editor string

// EditorResolver returns the editor to use based on user preferrence or default
type EditorResolver func() string

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write a new encrypted log",
	Long: `Generate a new journalentry. 
Entry will be encrypted when written. 
You may specify a custom file name and storage location.
	
The default name for a file is <DATE>_<TIME>.ck`,
	Run: func(cmd *cobra.Command, args []string) {
		journalEntryAsBytes, err := GetEditorInput()
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(fileName+FileExtension, journalEntryAsBytes, 0755)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(writeCmd)

	// Flags:

	//Persistent Flags:

	/**
	*	filename - changes the file name
	*	short hand: f
	*	default: <DATE>_<TIME>.ck
	 */
	writeCmd.PersistentFlags().StringVarP(&fileName, "filename", "f", DefaultFileName, "Sets the output file name. Default is <DATE>_<TIME>.ck")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	writeCmd.Flags().StringVarP(&passPhrase, "passphrase", "p", "", "Pass phrase used during enryption of log")
	writeCmd.MarkFlagRequired("passphrase")

	writeCmd.Flags().StringVarP(&editor, "editor", "e", DefaultEditor, "Allows you to change the ditor used to write files. Default is vim")
}

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
func resolveEditorArgs(command string, tempFilename string) []string {
	args := []string{tempFilename}

	if strings.Contains(command, "code") || strings.Contains(command, "vsc") {
		args = append([]string{"--wait"}, args...)
	}

	return args
}

/**
* Open users preffered text editor to capture log
 */
func OpenEditor(filename string, resolver EditorResolver) error {
	executable, err := exec.LookPath(resolver())
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(executable, resolveEditorArgs(executable, filename)...)
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

	if OpenEditor(tempFile, getPreferredEditor); err != nil {
		return []byte{}, err
	}

	bytes, err := os.ReadFile(tempFile)
	if err != nil {
		log.Fatal(err)
	}

	return bytes, nil
}
