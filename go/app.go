package main

import "flag"
import "fmt"
import "github.com/fogleman/gg"
import "image"
import "image/png"
import "github.com/disintegration/imaging"
import "os"
import "errors"

func initImage(width uint, height uint, imagePath string, colorBg string, colorBgOverlay string, blurStrength float64) *gg.Context {
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

		bgImageRes := imaging.Blur(
			imaging.Fill(imageData, int(width), int(height), imaging.Center, imaging.Lanczos),
			blurStrength)
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

	colorBg := flag.String("colorBg", "#000000", "Background color")
	colorBgOverlay := flag.String("colorBgOverlay", "#000000", "Background color overlay (overlay on the image)")
	colorFg := flag.String("colorFg", "#ffffff", "Text (foreground) color")
	siteName := flag.String("siteName", "", "Site name")
	title := flag.String("title", "", "Article title")
	imagePath := flag.String("imagePath", "", "Image path")
	logoPath := flag.String("logoPath", "", "Logo image path")
	logoX := flag.Int("logoX", 0, "Logo image X")
	logoY := flag.Int("logoY", 0, "Logo image Y")
	targetPath := flag.String("targetPath", "out.png", "Image path, use - for stdout")
	width := flag.Uint("width", 600, "Output image width")
	height := flag.Uint("height", 300, "Output image height")
	paddingX := flag.Float64("paddingX", 0, "Padding X")
	paddingY := flag.Float64("paddingY", 0, "Padding Y")
	lineSpacing := flag.Float64("lineSpacing", 1, "Title text line spacing")
	fontSize := flag.Uint("fontSize", 36, "Font size for title")
	fontSizeSite := flag.Uint("fontSizeSite", 20, "Font size for site name")
	paddingYSite := flag.Uint("paddingYSite", 0, "Padding Y for site name")
	blurStrength := flag.Float64("blurStrength", 3, "Blur strength")

	fontTitle := flag.String("fontTitle", "", "Font for title")
	fontSiteName := flag.String("fontSiteName", "", "Font for site name")

	flag.Parse()

	if len(*title) == 0 {
		panic(errors.New("title not defined"))
	}

	fmt.Fprintf(os.Stderr, "Color BG: %s\n", *colorBg)
	fmt.Fprintf(os.Stderr, "Color Text: %s\n", *colorFg)
	fmt.Fprintf(os.Stderr, "Title: %s\n", *title)
	fmt.Fprintf(os.Stderr, "Site Name: %s\n", *siteName)
	fmt.Fprintf(os.Stderr, "Image path: %s\n", *imagePath)
	fmt.Fprintf(os.Stderr, "Logo image path: %s\n", *logoPath)

	maxWidth := float64(*width) - 2*(*paddingX)

	img := initImage(*width, *height, *imagePath, *colorBg, *colorBgOverlay, *blurStrength)
	drawLogo(img, *logoPath, *logoX, *logoY)

	if len(*fontTitle) > 0 {
		err := img.LoadFontFace(*fontTitle, float64(*fontSize))
		if err != nil {
			panic(err)
		}
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
		if len(*fontSiteName) > 0 {
			err := img.LoadFontFace(*fontSiteName, float64(*fontSizeSite))
			if err != nil {
				panic(err)
			}
		}

		img.DrawString(*siteName, float64(*logoX), float64(*height)-float64(*paddingYSite))
	}

	if *targetPath == "-" {
		png.Encode(os.Stdout, img.Image())
	} else {
		img.SavePNG(*targetPath)
	}
}
