package persianbgbot

import (
	"bytes"
	"embed"
	genericloader "github.com/fzerorubigd/persianbgbot/pkg/generic"
	"log"
	"path/filepath"
)

//go:embed contrib
var contrib embed.FS

func loadDir(path string) {
	files, err := contrib.ReadDir(path)
	if err != nil {
		log.Println(err)
		return
	}

	for i := range files {
		if files[i].IsDir() {
			loadDir(filepath.Join(path, files[i].Name()))
			continue
		}

		b, err := contrib.ReadFile(filepath.Join(path, files[i].Name()))
		if err != nil {
			log.Println(err)
			continue
		}

		if err := genericloader.RegisterCard(bytes.NewBuffer(b)); err != nil {
			log.Println(err)
		}
	}

}

func init() {
	loadDir("contrib")
}
