/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/arcalot/arcaflow-plugin-image-builder/internal/carpenter"
	"github.com/spf13/cobra"
	"os"
)

var Push bool
var Build bool

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.PersistentFlags().BoolVarP(&Push, "push", "p", false, "push images to registry")
	buildCmd.PersistentFlags().BoolVarP(&Build, "build", "b", false, "validate requirements and build image")
}

var buildCmd = &cobra.Command{
	Use:   "build an image",
	Short: "build image",
	Run: func(cmd *cobra.Command, args []string) {
		err := carpenter.CliCarpentry(Build, Push, rootLogger, "docker")
		if err != nil {
			rootLogger.Errorf("build command failed (%w)", err)
			os.Exit(1)
		}
	},
}
