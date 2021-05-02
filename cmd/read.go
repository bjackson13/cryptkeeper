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
	"log"
	"os"

	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read an encrypted journal entry",
	Long: `Read allows you to select a cryptkeeper file and read the decrypted contents.
	Requires first argument to be the name of the file you wish to decrypt and read.
	pass phrase flag "-p" is required as well. This must match the passphrase used to encryt the journal.
	Example:
		cryptkeeper read my_journal_entry.ck -p mysecretpassphrase`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName = args[0]

		// Read encrypted journal into memory
		//TODO: Improve file not found error handling
		encryptedJournal, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatal(err)
		}

		// decrypt journal
		decryptedJournal, err := DecryptJournal(encryptedJournal, passPhrase)
		if err != nil {
			log.Fatal(err)
		}

		//display in editor
		DisplayOutputInEditor(decryptedJournal)
	},
}

func init() {
	rootCmd.AddCommand(readCmd)

	readCmd.Flags().StringVarP(&passPhrase, "passphrase", "p", "", "Pass phrase used during decryption of journal entry")
	readCmd.MarkFlagRequired("passphrase")

	readCmd.Flags().StringVarP(&editor, "editor", "e", DefaultEditor, "Allows you to change the editor used to read files. Default is vim")
}
