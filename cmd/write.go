/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var fileName string
var passPhrase string
var editor string

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write a new encrypted log",
	Long: `Generate a new journalentry. 
Entry will be encrypted when written. 
You may specify a custom file name and storage location.
	
The default name for a file is <DATE>_<TIME>.ck`,
	Run: func(cmd *cobra.Command, args []string) {
		journalBytes, err := GetEditorInput()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(string(journalBytes[:]))
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
	writeCmd.PersistentFlags().StringVarP(&fileName, "filename", "f", time.Now().Format("2006-01-02_1504"), "Sets the output file name. Default is <DATE>_<TIME>.ck")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	writeCmd.Flags().StringVarP(&passPhrase, "passphrase", "p", "", "Pass phrase used during enryption of log")
	writeCmd.MarkFlagRequired("passphrase")

	writeCmd.Flags().StringVarP(&editor, "editor", "e", "vim", "Allows you to change the ditor used to write files. Default is vim")
}

/**
Open users preffered text editor to capture log
*/
func OpenEditor(filename string) error {
	executable, err := exec.LookPath(editor)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(executable, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

/*
	Create a temporary file in a temp directory and open it in the users preffered editor.
	Delete the temp file once we read in the contents of it
*/
func GetEditorInput() ([]byte, error) {
	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		log.Fatal(err)
	}

	tempFile := file.Name()
	defer os.Remove(tempFile)

	if err = file.Close(); err != nil {
		return []byte{}, err
	}

	if OpenEditor(tempFile); err != nil {
		return []byte{}, err
	}

	bytes, err := ioutil.ReadFile(tempFile)
	if err != nil {
		log.Fatal(err)
	}

	return bytes, nil
}
