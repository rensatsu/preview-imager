package main

import "flag"
import "fmt"
import "github.com/fogleman/gg"
import "image"
import "image/png"
import "github.com/disintegration/imaging"
import "os"
import "errors"

// Inspiration:
// https://github.com/montanaflynn/meme-generator/blob/master/main.go

func initImage(width uint, height uint, imagePath string, colorBg string, colorBgOverlay string) *gg.Context {
	ctx := gg.NewContext(int(width), int(height))

	if len(imagePath) == 0 {
		ctx.SetHexColor(colorBg)
		ctx.DrawRectangle(0, 0, float64(ctx.Width()), float64(ctx.Height()))
		ctx.Fill()
	} else {
		existingImageFile, err := os.Open(imagePath)
		if err != nil {
			panic(err)
		}
		defer existingImageFile.Close()

		imageData, _, err := image.Decode(existingImageFile)
		if err != nil {
			panic(err)
		}

		bgImageRes := imaging.Blur(imaging.Fill(imageData, int(width), int(height), imaging.Center, imaging.Lanczos), 3)
		ctx.DrawImage(bgImageRes, 0, 0)
		ctx.DrawRectangle(0, 0, float64(width), float64(height))
		ctx.SetHexColor(colorBgOverlay)
		ctx.Fill()
	}

	return ctx
}

func drawLogo(img *gg.Context, logoPath string, x int, y int) {
	// reading logo file
	if len(logoPath) > 0 {
		logoImage, err := os.Open(logoPath)
		if err != nil {
			panic(err)
		}
		defer logoImage.Close()

		imageData, _, err := image.Decode(logoImage)
		if err != nil {
			panic(err)
		}

		img.DrawImage(imageData, img.Width()-x, img.Height()-y)
	}

}

func main() {
	fmt.Fprintln(os.Stderr, "Preview Imager by Rensatsu")

	colorBg := flag.String("color-bg", "#000000", "Background color")
	colorBgOverlay := flag.String("color-bg-overlay", "#000000", "Background color overlay (overlay on the image)")
	colorFg := flag.String("color-fg", "#ffffff", "Text (foreground) color")
	siteName := flag.String("site-name", "", "Site name")
	title := flag.String("title", "", "Article title")
	imagePath := flag.String("image", "", "Image path")
	logoPath := flag.String("logo-image", "", "Logo image path")
	logoX := flag.Int("logo-x", 0, "Logo image X")
	logoY := flag.Int("logo-y", 0, "Logo image Y")
	targetPath := flag.String("target", "out.png", "Image path, use - for stdout")
	width := flag.Uint("width", 600, "Output image width")
	height := flag.Uint("height", 300, "Output image height")
	paddingX := flag.Float64("padding-x", 0, "Padding X")
	paddingY := flag.Float64("padding-y", 0, "Padding Y")
	lineSpacing := flag.Float64("line-spacing", 1, "Title text line spacing")
	fontSize := flag.Uint("font-size", 36, "Font size for title")
	fontSizeSite := flag.Uint("font-size-site", 20, "Font size for site name")
	paddingYSite := flag.Uint("padding-y-site", 0, "Padding Y for site name")

	flag.Parse()

	if len(*title) == 0 {
		panic(errors.New("title not defined"))
	}

	fmt.Fprintf(os.Stderr, "Color BG: %s\n", *colorBg)
	fmt.Fprintf(os.Stderr, "Color Text: %s\n", *colorFg)
	fmt.Fprintf(os.Stderr, "Title: %s\n", *title)
	fmt.Fprintf(os.Stderr, "Site Name: %s\n", *siteName)

	maxWidth := float64(*width) - 2*(*paddingX)

	img := initImage(*width, *height, *imagePath, *colorBg, *colorBgOverlay)
	drawLogo(img, *logoPath, *logoX, *logoY)

	err := img.LoadFontFace("fonts/IBMPlexSans-Bold.ttf", float64(*fontSize))
	if err != nil {
		panic(err)
	}

	img.SetHexColor(*colorFg)
	img.DrawStringWrapped(
		*title,
		float64(*paddingX),
		float64(*paddingY),
		0,
		0,
		maxWidth,
		float64(*lineSpacing),
		gg.AlignLeft)

	if len(*siteName) > 0 {
		err := img.LoadFontFace("fonts/IBMPlexSans-Regular.ttf", float64(*fontSizeSite))
		if err != nil {
			panic(err)
		}

		img.DrawString(*siteName, float64(*logoX), float64(*height)-float64(*paddingYSite))
	}

	if *targetPath == "-" {
		png.Encode(os.Stdout, img.Image())
	} else {
		img.SavePNG(*targetPath)
	}
}
