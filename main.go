package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/urfave/cli"
	"gopkg.in/gographics/imagick.v2/imagick"
)

func main() {
	app := cli.NewApp()
	app.Name = "depicture"
	app.Usage = "generate colorschemes from wallpapers"
	app.Version = "0.1.0"

	app.Action = func(c *cli.Context) error {
		filePath := c.Args().Get(0)
		if len(filePath) < 1 {
			fmt.Printf("Error: Please provide a file.\n")
			os.Exit(1)
		}

		imagick.Initialize()
		defer imagick.Terminate()

		mw := imagick.NewMagickWand()

		if err := mw.ReadImage(filePath); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if err := mw.QuantizeImage(16, imagick.COLORSPACE_SRGB, 0, false, false); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		colorCount := mw.GetImageColors()
		if colorCount < 16 {
			fmt.Printf("Error: Couldn't find enough colors.\n")
			os.Exit(1)
		}

		var colors []string

		for i := uint(0); i < 16; i++ {
			pw, err := mw.GetImageColormapColor(i)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			defer pw.Destroy()

			r := round(pw.GetRed() * 255.0)
			g := round(pw.GetGreen() * 255.0)
			b := round(pw.GetBlue() * 255.0)

			hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)

			colors = append(colors, hex)
		}
		mw.Destroy()

		commentColor := colors[0]
		n := string([]rune(commentColor)[1])

		switch n {
		case "0", "1":
			colors = append(colors, "#666666")
		case "2":
			colors = append(colors, "#757575")
		case "3", "4":
			colors = append(colors, "#999999")
		case "5":
			colors = append(colors, "#8a8a8a")
		case "6-9":
			colors = append(colors, "#a1a1a1")
		default:
			colors = append(colors, colors[15])
		}

		type templateData struct {
			Colors []string
		}

		data := &templateData{Colors: colors}

		xTemplate, err := template.ParseFiles("./templates/Xresources")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if err := xTemplate.Execute(os.Stdout, data); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		return nil
	}

	app.Run(os.Args)
}
