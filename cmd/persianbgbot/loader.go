package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	genericloader "github.com/fzerorubigd/persianbgbot/pkg/generic"

	"github.com/pkg/errors"
)

func loadFile(file string) error {
	fl, err := os.Open(file)
	if err != nil {
		return errors.Wrapf(err, "open file %q failed", file)
	}
	defer func() {
		_ = fl.Close()
	}()

	return genericloader.RegisterCard(fl)
}

func loadPath(path string) error {
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if err := loadFile(path); err != nil {
			log.Printf("Load file %q failed, ignoring", path)
		}

		return nil
	})

	return errors.Wrapf(err, "walk path %q failed", path)
}
