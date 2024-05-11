package utils

import (
	"errors"
	"image"
	"strconv"

	"github.com/ReanSn0w/tk4go/pkg/tools"
	"github.com/adrg/sysfont"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func NewList(log tools.Logger, ffmpeg *FFMPEG) *List {
	return &List{
		log:    log,
		ffmpeg: ffmpeg,
	}
}

type List struct {
	log    tools.Logger
	ffmpeg *FFMPEG
}

func (l *List) Make(spec *Spec, grid, width int, images []image.Image) (image.Image, error) {
	specImage, err := l.specImage(spec)
	if err != nil {
		return nil, err
	}

	resolution := [2]int{}
	if len(images) > 0 {
		resolution[0], resolution[1] = images[0].Bounds().Dx(), images[0].Bounds().Dy()
	}

	if len(images) == 0 {
		return nil, errors.New("no images to make list from")
	}

	list, err := l.grid(resolution, grid, width, images)
	if err != nil {
		return nil, err
	}

	clipped := l.clipImages(specImage, list)
	return clipped, nil
}

func (l *List) grid(resolution [2]int, grid, width int, images []image.Image) (image.Image, error) {
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

func (l *List) clipImages(images ...image.Image) image.Image {
	if images[0] == nil {
		return images[1]
	}

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

func (l *List) specImage(spec *Spec) (image.Image, error) {
	if spec == nil {
		return nil, nil
	}

	font := l.loadFont()
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

func (l *List) loadFont() string {
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
