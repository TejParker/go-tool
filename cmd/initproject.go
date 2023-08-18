package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	new2 "go-tool/module"
)

var new = &cobra.Command{
	Use:   "new",
	Short: "Simon's customize service, automatically generate Service directory!",
	Long:  `This Command will create project base directory fastly and generate project basic structure for convenience~`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Tips:\n\t Please  input project name！")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		new2.NewServiceProject(cmd, args)
	},
}
