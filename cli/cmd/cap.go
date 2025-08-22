// Package cmd implements the CLI commands for the WMS tool.
package cmd

import (
	"fmt"

	"github.com/sgaunet/wms/getmap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getcapCommand = &cobra.Command{
	Use:     "cap",
	Aliases: []string{"getcap"},
	Args:    cobra.RangeArgs(0, 1),
	Short:   "Get the capabilities of a WMS",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		service := "default"
		if len(args) == 1 {
			service = args[0]
		}
		url := viper.GetString(service + ".url")
		if cmd.Flag("url").Changed {
			url, err = cmd.Flags().GetString("url")
			if err != nil {
				return fmt.Errorf("getting url: %w", err)
			}
		}
		if url == "" {
			return ErrEmptyURL
		}
		w := &getmap.Service{}
		err = w.SetURL(url)
		if err != nil {
			return fmt.Errorf("setting URL: %w", err)
		}
		version := viper.GetString(service + ".version")
		if cmd.Flag("version").Changed {
			version, err = cmd.Flags().GetString("version")
			if err != nil {
				return fmt.Errorf("getting version: %w", err)
			}
		}
		err = w.AddVersion(version)
		if err != nil {
			return fmt.Errorf("adding version: %w", err)
		}
		f, err := cmd.Flags().GetBool("formats")
		if err != nil {
			return fmt.Errorf("getting formats flag: %w", err)
		}
		l, err := cmd.Flags().GetBool("layers")
		if err != nil {
			return fmt.Errorf("getting layers flag: %w", err)
		}
		e, err := cmd.Flags().GetBool("epsg")
		if err != nil {
			return fmt.Errorf("getting epsg flag: %w", err)
		}
		c, err := w.GetCapabilities()
		if err != nil {
			return fmt.Errorf("getting capabilities: %w", err)
		}
		if !f && !l && !e {
			fmt.Println(c)
		}
		if f {
			fmt.Println(c.Formats)
		}
		if l {
			layers := c.GetLayerNames()
			for _, l := range layers {
				fmt.Println(l)
			}
		}
		if e {
			fmt.Println(c.GetBBoxes().GetEPSG())
		}
		return nil
	},
}

func init() {
	getcapCommand.Flags().StringP("url", "u", "", "Set url")
	getcapCommand.Flags().StringP("version", "v", "", "Set version")

	getcapCommand.Flags().BoolP("formats", "f", false, "Get available formats")
	getcapCommand.Flags().BoolP("layers", "l", false, "Get available layers")
	getcapCommand.Flags().BoolP("epsg", "e", false, "Get available epsg-codes")

	root.AddCommand(getcapCommand)
}
