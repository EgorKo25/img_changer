package rotater

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"golang.org/x/image/tiff"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"path"
	"sync"
)

type Rotater struct {
	files []os.DirEntry
	wg    sync.WaitGroup
}

func NewRotater() (*Rotater, error) {
	return &Rotater{}, nil
}

func (r *Rotater) Rotate(in, fileName, out string, angle float64) error {

	imagePath, err := os.Open(path.Join(in, fileName))
	if err != nil {
		r.wg.Done()
		return err
	}
	defer func() {
		err = imagePath.Close()
		if err != nil {
			log.Printf("file error: %s", err)
			return
		}
	}()

	log.Printf("%s%s", in, fileName)
	srcImage, err := tiff.Decode(imagePath)
	if err != nil {
		r.wg.Done()
		return errors.New(fmt.Sprintf("decode error: %s", err))
	}

	rotImage := imaging.Rotate(srcImage, angle, color.White)

	newImage, err := os.Create(path.Join(out, fileName[:len(fileName)-5]+fmt.Sprintf("_%.f", angle)+".bmp"))
	if err != nil {
		r.wg.Done()
		return err
	}
	defer func() {
		err = newImage.Close()
		if err != nil {
			r.wg.Done()
			log.Printf("file error: %s", err)
			return
		}
	}()

	err = jpeg.Encode(newImage, rotImage, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		r.wg.Done()
		return err
	}

	r.wg.Done()
	return nil
}

func (r *Rotater) RotateAll(in, out string) (err error) {

	r.files, err = os.ReadDir(in)
	if err != nil {
		return err
	}

	for _, file := range r.files {
		for i := -25; i < 30; i += 5 {
			r.wg.Add(1)
			go func(angle float64) {
				err = r.Rotate(in, file.Name(), out, angle)
				if err != nil {
					log.Printf("gorutine error: %s", err)
				}
			}(float64(i))

		}
		r.wg.Wait()

	}

	return nil
}
