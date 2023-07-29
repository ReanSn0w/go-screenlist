package list

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"os"
	"strconv"

	"github.com/ReanSn0w/go-screenlist/pkg/video"
	"github.com/adrg/sysfont"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func Make(spec *video.Spec, grid, width int, images []image.Image) (image.Image, error) {
	specImage, err := specImage(spec)
	if err != nil {
		return nil, err
	}

	list, err := makeImagesGrid(spec.Resolution, grid, width, images)
	if err != nil {
		return nil, err
	}

	clipped := clipImages(specImage, list)
	return clipped, nil
}

func Save(name string, spes *video.Spec, grid, width int, images []image.Image) error {
	image, err := Make(spes, grid, width, images)
	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, image, nil)
	if err != nil {
		return err
	}

	err = os.WriteFile(name, buffer.Bytes(), 0644)
	return err
}

func makeImagesGrid(resolution [2]int, grid, width int, images []image.Image) (image.Image, error) {
	if len(images) == 0 {
		return nil, errors.New("no images to make list from")
	}

	rows := len(images) / grid
	if len(images)%grid != 0 {
		rows++
	}

	listWidth := resolution[0]*grid + 10*grid + 1
	listHeight := resolution[1]*rows + 10*rows + 1

	image := image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: listWidth,
			Y: listHeight,
		},
	})

	for row := 0; row < rows; row++ {
		for line := 0; line < grid; line++ {
			index := row*grid + line
			if index >= len(images) {
				break
			}

			x := line*resolution[0] + line*10 + 1
			y := row*resolution[1] + row*10 + 1

			for i := 0; i < resolution[0]; i++ {
				for j := 0; j < resolution[1]; j++ {
					image.Set(x+i, y+j, images[index].At(i, j))
				}
			}
		}
	}

	resized := resize.Resize(uint(width), 0, image, resize.Lanczos3)
	return resized, nil
}

func clipImages(images ...image.Image) image.Image {
	width, height := 0, (len(images)-1)*10
	for _, image := range images {
		if image == nil {
			height -= 10
			continue
		}

		width = image.Bounds().Max.X
		height += image.Bounds().Max.Y
	}

	image := image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: width,
			Y: height,
		},
	})

	ySkip := 0

	for _, item := range images {
		for i := 0; i < item.Bounds().Max.X; i++ {
			for j := 0; j < item.Bounds().Max.Y; j++ {
				image.Set(i, j+ySkip, item.At(i, j))
			}
		}

		ySkip += item.Bounds().Max.Y + 10
	}

	return image
}

func specImage(spec *video.Spec) (image.Image, error) {
	font := loadFont()
	if font == "" {
		return nil, errors.New("no font found")
	}

	values := []string{
		"Title: " + spec.Title,
		"Resolution: " + strconv.Itoa(spec.Resolution[0]) + "x" + strconv.Itoa(spec.Resolution[1]),
		"Fps: " + strconv.Itoa(spec.Fps),
		"Size: " + strconv.Itoa(spec.Size),
		"Duration: " + strconv.Itoa(spec.Duration/60) + "m " + strconv.Itoa(spec.Duration%60) + "s ",
	}

	dc := gg.NewContext(spec.Resolution[0], len(values)*18+10)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	if err := dc.LoadFontFace(font, 14); err != nil {
		return nil, err
	}

	for i, value := range values {
		dc.DrawStringAnchored(value, 10, float64(18*(i+1)), 0, 0)
	}

	return dc.Image(), nil
}

func loadFont() string {
	finder := sysfont.NewFinder(nil)
	// if font := finder.Match("Helvetica"); font != nil {
	// 	return font.Filename
	// }

	if font := finder.Match("Arial"); font != nil {
		return font.Filename
	}

	if font := finder.Match("Times New Roman"); font != nil {
		return font.Filename
	}

	return ""
}
