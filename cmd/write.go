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
	"time"

	"github.com/spf13/cobra"
)

var DefaultFileName string = time.Now().Format("2006-01-02_1504")

/**
* Flag variables
 */
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
		// get the input from editor
		journalEntryAsBytes, err := GetEditorInput()
		if err != nil {
			log.Fatal(err)
		}

		//encrypt using passphrase
		encryptedJournal, err := EncryptJournal(journalEntryAsBytes, passPhrase)
		if err != nil {
			log.Fatal(err)
		}

		//write encrypted input to file
		err = os.WriteFile(fileName+FileExtension, encryptedJournal, 0700)
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
	writeCmd.Flags().StringVarP(&passPhrase, "passphrase", "p", "", "Pass phrase used during encryption of journal entry")
	writeCmd.MarkFlagRequired("passphrase")

	writeCmd.Flags().StringVarP(&editor, "editor", "e", DefaultEditor, "Allows you to change the editor used to write files. Default is vim")
}
