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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		service := "default"
		if len(args) == 1 {
			service = args[0]
		}
		url := viper.GetString(service + ".url")
		if cmd.Flag("url").Changed {
			url, err = cmd.Flags().GetString("url")
			if err != nil {
				return
			}
		}
		w := &getmap.Service{}
		err = w.SetURL(url)
		if err != nil {
			return
		}
		version := viper.GetString(service + ".version")
		if cmd.Flag("version").Changed {
			version, err = cmd.Flags().GetString("version")
			if err != nil {
				return
			}
		}
		err = w.AddVersion(version)
		if err != nil {
			return
		}
		f, err := cmd.Flags().GetBool("formats")
		if err != nil {
			return
		}
		l, err := cmd.Flags().GetBool("layers")
		if err != nil {
			return
		}
		e, err := cmd.Flags().GetBool("epsg")
		if err != nil {
			return
		}
		c, err := w.GetCapabilities()
		if err != nil {
			return
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
		return
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
