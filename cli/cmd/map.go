package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/sgaunet/wms/getmap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	bboxExpectedParts      = 4
	layerStyleSeparatorLen = 2
	percentToRatio         = 100
	defaultDirPerm         = 0o750
)

var getmapCommand = &cobra.Command{
	Use:     "map",
	Aliases: []string{"getmap"},
	Args:    cobra.RangeArgs(0, 1),
	Short:   "Download a WMS-Tile",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		s := &getmap.Service{}
		service := "default"
		if len(args) == 1 {
			service = args[0]
		}
		url := viper.GetString(service + ".url")
		if cmd.Flag("url").Changed {
			url, err = cmd.Flags().GetString("url")
			if err != nil {
				return fmt.Errorf("getting url flag: %w", err)
			}
		}
		if url == "" {
			return ErrEmptyURL
		}
		err = s.SetURL(url)
		if err != nil {
			return fmt.Errorf("setting URL: %w", err)
		}
		version := viper.GetString(service + ".version")
		if cmd.Flag("version").Changed {
			version, err = cmd.Flags().GetString("version")
			if err != nil {
				return fmt.Errorf("getting version flag: %w", err)
			}
		}
		err = s.AddVersion(version)
		if err != nil {
			return fmt.Errorf("adding version: %w", err)
		}
		format := viper.GetString(service + ".format")
		if cmd.Flag("format").Changed {
			format, err = cmd.Flags().GetString("format")
			if err != nil {
				return fmt.Errorf("getting format flag: %w", err)
			}
		}
		if format != "" {
			err = s.AddFormat(format)
			if err != nil {
				return fmt.Errorf("adding format: %w", err)
			}
		}
		layers := viper.GetStringSlice(service + ".layers")
		if cmd.Flag("layers").Changed {
			layers, err = cmd.Flags().GetStringSlice("layers")
			if err != nil {
				return fmt.Errorf("getting layers flag: %w", err)
			}
		}
		layerReq := []string{}
		stylesReq := make(map[string]string)
		for _, l := range layers {
			split := strings.Split(l, "/")
			layerReq = append(layerReq, split[0])
			if len(split) == layerStyleSeparatorLen {
				stylesReq[split[0]] = split[1]
			}
		}
		if len(layerReq) != 0 {
			err = s.AddLayers(layerReq...)
			if err != nil {
				return fmt.Errorf("adding layers: %w", err)
			}
		}
		for l, sty := range stylesReq {
			err = s.AddStyle(l, sty)
			if err != nil {
				return fmt.Errorf("adding style: %w", err)
			}
		}
		epsg := viper.GetInt(service + ".epsg")
		if cmd.Flag("epsg").Changed {
			epsg, err = cmd.Flags().GetInt("epsg")
			if err != nil {
				return fmt.Errorf("getting epsg flag: %w", err)
			}
		}
		if epsg != 0 {
			err = s.AddEPSG(epsg)
			if err != nil {
				return fmt.Errorf("adding EPSG: %w", err)
			}
		}
		name := viper.GetString(service + ".file-name")
		if cmd.Flag("file-name").Changed {
			name, err = cmd.Flags().GetString("file-name")
			if err != nil {
				return fmt.Errorf("getting file-name flag: %w", err)
			}
		}
		if name == "" {
			name = "example"
		}
		if cmd.Flag("save").Changed {
			save, err := cmd.Flags().GetString("save")
			if err != nil {
				return fmt.Errorf("getting save flag: %w", err)
			}
			viper.Set(save+".url", url)
			viper.Set(save+".version", version)
			viper.Set(save+".format", format)
			viper.Set(save+".layers", layers)
			viper.Set(save+".epsg", epsg)
			viper.Set(save+".file-name", name)
			fmt.Println("Saving service:", save)
			err = viper.WriteConfig()
			if err != nil {
				return fmt.Errorf("writing config: %w", err)
			}
		}
		if cmd.Flag("dry-run").Changed {
			fmt.Println(s)
			fmt.Println("File name:", name)
		} else {
			var bboxes [][]float64
			if cmd.Flag("bbox").Changed {
				bboxStr, err := cmd.Flags().GetStringSlice("bbox")
				if err != nil {
					return fmt.Errorf("getting bbox flag: %w", err)
				}
				if len(bboxStr) != bboxExpectedParts {
					return ErrInvalidBBox
				}
				minx, err := strconv.ParseFloat(bboxStr[0], 64)
				if err != nil {
					return fmt.Errorf("parsing minx: %w", err)
				}
				miny, err := strconv.ParseFloat(bboxStr[1], 64)
				if err != nil {
					return fmt.Errorf("parsing miny: %w", err)
				}
				maxx, err := strconv.ParseFloat(bboxStr[2], 64)
				if err != nil {
					return fmt.Errorf("parsing maxx: %w", err)
				}
				maxy, err := strconv.ParseFloat(bboxStr[3], 64)
				if err != nil {
					return fmt.Errorf("parsing maxy: %w", err)
				}
				bboxes = append(bboxes, []float64{minx, miny, maxx, maxy})
			} else if cmd.Flag("bbox-file").Changed {
				bf, err := cmd.Flags().GetString("bbox-file")
				if err != nil {
					return fmt.Errorf("getting bbox-file flag: %w", err)
				}
				file, err := os.Open(bf) // #nosec G304
				if err != nil {
					return fmt.Errorf("opening bbox file: %w", err)
				}
				defer func() {
					if cerr := file.Close(); cerr != nil {
						err = fmt.Errorf("closing file: %w", cerr)
					}
				}()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					split := strings.Split(line, ",")
					if len(split) != bboxExpectedParts {
						return fmt.Errorf("invalid: BBox-File %s: %w", bf, ErrInvalidBBox)
					}
					var bboxL []float64
					for _, b := range split {
						i, err := strconv.ParseFloat(b, 64)
						if err != nil {
							return fmt.Errorf("parsing bbox coordinate: %w", err)
						}
						bboxL = append(bboxL, i)
					}
					bboxes = append(bboxes, bboxL)
				}
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("scanning file: %w", err)
				}
			} else {
				return ErrInvalidBBoxReq
			}
			width, err := cmd.Flags().GetInt("width")
			if err != nil {
				return fmt.Errorf("getting width flag: %w", err)
			}
			height, err := cmd.Flags().GetInt("height")
			if err != nil {
				return fmt.Errorf("getting height flag: %w", err)
			}
			scale, err := cmd.Flags().GetInt("scale")
			if err != nil {
				return fmt.Errorf("getting scale flag: %w", err)
			}
			dpi, err := cmd.Flags().GetInt("dpi")
			if err != nil {
				return fmt.Errorf("getting dpi flag: %w", err)
			}
			expand, err := cmd.Flags().GetInt("expand")
			if err != nil {
				return fmt.Errorf("getting expand flag: %w", err)
			}
			c, err := cmd.Flags().GetBool("cut")
			if err != nil {
				return fmt.Errorf("getting cut flag: %w", err)
			}
			errs := make(chan error, len(bboxes))
			for i, b := range bboxes {
				go getImageC(s, bboxes, i, b, expand, width, height, scale, dpi, name, c, errs)
			}
			for range bboxes {
				x := <-errs
				if x != nil {
					return x
				}
			}
			pwd, _ := os.Getwd()
			fmt.Println("Done. Your requested file is here: " + filepath.Join(pwd, "output"))
		}
		return nil
	},
}

func init() {
	getmapCommand.Flags().StringP("url", "u", "", "Set url")
	getmapCommand.Flags().StringP("version", "v", "", "Set version")
	getmapCommand.Flags().StringP("format", "f", "", "Set format")
	getmapCommand.Flags().StringSliceP("layers", "l", nil, "Set layers")
	getmapCommand.Flags().IntP("epsg", "e", 0, "Set epsg-code")

	getmapCommand.Flags().IntP("width", "w", 0, "Set width of output image in px")
	getmapCommand.Flags().IntP("height", "h", 0, "Set height of output image in px")
	getmapCommand.Flags().IntP("scale", "s", 0, "Set scale of output image (dpi required!)")
	getmapCommand.Flags().IntP("dpi", "i", 0, "Set dpi of output image (scale required!)")

	getmapCommand.Flags().StringSliceP("bbox", "b", nil, "Set bbox in meters (minx,miny,maxx,maxy)")
	getmapCommand.Flags().StringP("bbox-file", "B", "", "Set bbox file")

	getmapCommand.Flags().IntP("expand", "E", 0, "Expands bbox in %")
	getmapCommand.Flags().BoolP("cut", "C", false,
		"Cuts image to unexpanded bbox (interesting for dynamically generated texts and symbols)")

	getmapCommand.Flags().StringP("file-name", "n", "", "Set file name")

	getmapCommand.Flags().Bool("dry-run", false, "Validate your request without executing")
	getmapCommand.Flags().String("save", "", "Save your request settings")

	getmapCommand.Flags().StringP("user", "", "", "Set user for Basic Authentication")
	getmapCommand.Flags().StringP("password", "", "", "Set password for Basic Authentication")

	root.AddCommand(getmapCommand)
}

func getImageC(s *getmap.Service, bboxes [][]float64, i int, b []float64, expand int,
	width, height, scale, dpi int, name string, cut bool, errs chan error) {
	var r *bytes.Reader
	var err error
	expandX := (b[2] - b[0]) * float64(expand) / percentToRatio
	expandY := (b[3] - b[1]) * float64(expand) / percentToRatio
	if width != 0 || height != 0 {
		width = int(math.Round(float64(width) * (1 + float64(expand)/100)))
		height = int(math.Round(float64(height) * (1 + float64(expand)/100)))
		r, width, height, err = s.GetMap(
			b[0]-expandX/2, b[1]-expandY/2, b[2]+expandX/2, b[3]+expandY/2,
			getmap.WidthHeightOption(width, height))
	} else {
		r, width, height, err = s.GetMap(
			b[0]-expandX/2, b[1]-expandY/2, b[2]+expandX/2, b[3]+expandY/2,
			getmap.ScaleDPIOption(scale, dpi))
	}
	if err != nil {
		errs <- err
		return
	}
	ext := s.GetFileExt()
	if ext == "" {
		ext = "png"
	}
	_ = os.MkdirAll("output", defaultDirPerm)
	var filePath string
	if len(bboxes) > 1 {
		filePath = filepath.Join("output", fmt.Sprintf("%02d_%v.%v", i+1, name, ext))
	} else {
		filePath = filepath.Join("output", fmt.Sprintf("%v.%v", name, ext))
	}
	if err != nil {
		errs <- err
		return
	}
	img, err := imaging.Decode(r)
	if err != nil {
		errs <- err
		return
	}
	if cut {
		width = int(math.Round(float64(width) / (1 + float64(expand)/100)))
		height = int(math.Round(float64(height) / (1 + float64(expand)/100)))
		img = imaging.CropCenter(img, width, height)
	}
	err = imaging.Save(img, filePath)
	if err != nil {
		errs <- err
		return
	}
	errs <- err
}
