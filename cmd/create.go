/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [file_name]",
	Short: "Creates a migration file in current directory with provided name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := fmt.Sprintf("%d_%s.sql", time.Now().Unix(), args[0])

		file, err := os.Create(fileName)
		if err != nil {
			slog.Error("cannot create migration file",
				err)
			return
		}
		defer file.Close()

		content := "-- +seeker Up\n --write your up-migration here--\n-- +seeker Down\n--write your down-migration here--"

		if _, err = file.WriteString(content); err != nil {
			slog.Error("cannot write content to migration file",
				err)
			return
		}

		slog.Info("migration file have been successfully created",
			"filename",
			fileName)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
