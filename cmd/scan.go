/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/vucchaid/scanner/pkg/clamav"
	"github.com/vucchaid/scanner/pkg/docker"
	"github.com/vucchaid/scanner/pkg/log"

	"os"

	"github.com/spf13/cobra"
)

const (
	imageName = "ajilaag/clamav-rest"
)

var rootCmd = &cobra.Command{
	Use: "scanner",
	Long: `Scanner is an utility program to scan files for viruses and malwares.
This tool uses clamav to do the findings. 
For more information on clamav, visit https://www.clamav.net

Please note: This utility is not stable, meaning - it is written in free time, for fun.
It is just an easy way to setup, scan files through clamav. It may or may not be updated in future.

Requirement : Docker must be running on your system.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var scanCmd = &cobra.Command{
	Use: "scan",
	Long: `
scanner scan command is used to trigger a scan. Usage: scanner scan -f "<fileName>"`,
	Short: `scan command is used to trigger a scan. Usage: scanner scan -f "<fileName>"`,
	Run: func(cmd *cobra.Command, args []string) {

		logger, err := log.GetLogger()
		if err != nil {
			return
		}

		fileName, _ := cmd.Flags().GetString("file")

		if fileName == "" {
			logger.Warn("must provide file to proceed. exiting...")
			return
		}

		if err = docker.CheckAndRunDockerFuncs(imageName, logger); err != nil {
			logger.Warn("caught error: " + err.Error())
			return
		}

		clamav.Interact(fileName, logger)

	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP("file", "f", "", "specify file to be scanned")
}
