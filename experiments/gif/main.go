package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/anthonynsimon/bild/effect"
	"github.com/pkg/errors"
)

func main() {
	f, err := os.Open("t-9-0.png")
	if err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}

	fmt.Printf("%s  -  %v\n", format, img.Bounds())

	wf, err := os.Create("wowza.png")
	if err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}
	defer wf.Close()

	//dst := image.NewGray(img.Bounds())
	//Canny(dst, img)
	img = effect.Grayscale(img)
	img = effect.Sobel(img)
	img = effect.Invert(img)

	err = png.Encode(wf, img)
	if err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}
}

func Canny(dst *image.Gray, src image.Image) error {
	if dst == nil {
		return errors.New("edge: dst is nil")
	}
	if src == nil {
		return errors.New("edge: src is nil")
	}

	b := src.Bounds()
	srcGray, ok := src.(*image.Gray)
	if !ok {
		srcGray = image.NewGray(b)
		draw.Draw(srcGray, b, src, b.Min, draw.Src)
	}

	mag, dir := image.NewGray(b), image.NewGray(b)
	//if err := Sobel(mag, dir, srcGray); err != nil {
	//	return err
	//}

	// Non-maximum supression.
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			d := dir.Pix[(y-b.Min.Y)*dir.Stride+(x-b.Min.X)*1]
			var m0, m1 uint8
			switch d {
			case 0: // west and east
				m0 = atOrZero(mag, x-1, y)
				m1 = atOrZero(mag, x+1, y)
			case 45: // north-east and south-west
				m0 = atOrZero(mag, x+1, y-1)
				m1 = atOrZero(mag, x-1, y+1)
			case 90: // north and south
				m0 = atOrZero(mag, x, y-1)
				m1 = atOrZero(mag, x, y+1)
			case 135: // north-west and south-east
				m0 = atOrZero(mag, x-1, y-1)
				m1 = atOrZero(mag, x+1, y+1)
			default:
				return fmt.Errorf("edge: bad direction (%d, %d): %d", x, y, d)
			}

			m := mag.Pix[(y-b.Min.Y)*mag.Stride+(x-b.Min.X)*1]
			if m > m0 && m > m1 {
				m = 0xff
			} else if m > m0 || m > m1 {
				m = 0x80
			} else {
				m = 0x00
			}
			dst.Pix[(y-b.Min.Y)*dst.Stride+(x-b.Min.X)*1] = m
		}
	}

	return nil
}

func atOrZero(m *image.Gray, x, y int) uint8 {
	if !image.Pt(x, y).In(m.Rect) {
		return 0
	}
	return m.Pix[(y-m.Rect.Min.Y)*m.Stride+(x-m.Rect.Min.X)*1]
}
