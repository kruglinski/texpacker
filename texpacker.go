// texpacker project texpacker.go

package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func image_info(name string) (image_width int, image_height int, image_type string, err error) {
	fd, err := os.Open(name)
	if err != nil {
		return 0, 0, "", err
	}

	defer fd.Close()

	img, image_type, err := image.Decode(fd)
	if err != nil {
		return 0, 0, "", err
	}

	img.ColorModel()

	size := img.Bounds().Size()
	return size.X, size.Y, image_type, nil
}

func image_load(name string) (img image.Image, err error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	defer fd.Close()
	img, _, err = image.Decode(fd)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func main() {
	var (
		in_folder string
		out_file  string
	)

	flag.StringVar(&in_folder, "i", "", "input folder")
	flag.StringVar(&out_file, "o", "", "output file")
	flag.Parse()

	_, err := os.Stat(in_folder)
	if err != nil {
		fmt.Println()
		flag.Usage()
		fmt.Println()
		return
	}

	out_fd, err := os.Create(out_file)
	if err != nil {
		fmt.Println("[-]", err)
		return
	}

	defer out_fd.Close()

	var paths []string
	err = filepath.Walk(in_folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		paths = append(paths, path)
		return nil
	})

	if err != nil {
		fmt.Println("[-]", err)
		return
	}

	image_width, image_height, _, err := image_info(paths[0])
	if err != nil {
		fmt.Println("[-]", err)
		return
	}

	rows := int(math.Sqrt(float64(len(paths))))
	cols := rows
	if rows*cols < len(paths) {
		rows += 1
	}

	sheet_width := cols * image_width
	sheet_height := rows * image_height

	rect := image.Rect(0, 0, sheet_width, sheet_height)
	img := image.NewRGBA(rect)
	img_count := 0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Println("adding...", paths[img_count])
			cell_image, err := image_load(paths[img_count])
			if err != nil {
				fmt.Println("[-]", err)
				defer os.Remove(out_file)
				return
			}

			grid_rect := image.Rect(j*image_width, i*image_height, j*image_width+image_width, i*image_height+image_height)
			draw.Draw(img, grid_rect, cell_image, image.ZP, draw.Src)

			img_count++
		}
	}

	err = png.Encode(out_fd, img)
	if err != nil {
		fmt.Println("[-]", err)
		defer os.Remove(out_file)
	} else {
		fmt.Println(out_file, "generated, total", img_count, "images packed!")
	}
}
