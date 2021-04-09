package main

import (
	"io/fs"
	"os"
	"path/filepath"

	genericloader "github.com/fzerorubigd/persianbgbot/internal/generic"

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

		return loadFile(path)
	})

	return errors.Wrapf(err, "walk path %q failed", path)
}
