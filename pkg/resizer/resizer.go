package resizer

import (
	"image"
	"log"
	"os"
	"path"
	"sync"

	"github.com/nfnt/resize"
	"image/jpeg"
)

type Resizer struct {
	files []os.DirEntry
	wg    sync.WaitGroup
}

func NewResizer() (*Resizer, error) {
	return &Resizer{}, nil
}

func (r *Resizer) Resize(filePath string, fileName string, out string) error {

	file, err := os.Open(path.Join(filePath, fileName))
	if err != nil {
		r.wg.Done()
		return err
	}

	if stat, _ := file.Stat(); stat.IsDir() {
		log.Printf("%t :: %s", !stat.IsDir(), fileName)
		r.wg.Done()
		return nil
	}
	srcImg, _, err := image.Decode(file)
	if err != nil {
		r.wg.Done()
		return err
	}

	newImage := resize.Resize(28, 28, srcImg, resize.Lanczos3)

	newFile, err := os.Create(path.Join(out, fileName[:len(fileName)-4]+".bmp"))
	if err != nil {
		return err
	}
	defer func() {
		err = newFile.Close()
		if err != nil {
			log.Printf("file error: %s", err)
			r.wg.Done()
			return
		}
	}()

	err = jpeg.Encode(newFile, newImage, nil)
	if err != nil {
		r.wg.Done()
		return err
	}

	r.wg.Done()
	return nil
}

func (r *Resizer) ResizeAll(in, out string) (err error) {

	r.files, err = os.ReadDir(in)
	if err != nil {
		return err
	}

	for _, file := range r.files {
		r.wg.Add(1)
		f := file

		go func() {
			err = r.Resize(in, f.Name(), out)
			if err != nil {
				log.Printf("resize error: %s", err)
			}
		}()
	}
	r.wg.Wait()

	return nil
}
