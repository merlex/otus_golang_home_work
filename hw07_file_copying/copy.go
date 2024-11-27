package main

import (
	"errors"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSourcePathEmpty       = errors.New("source path is required")
	ErrDestinationPathEmpty  = errors.New("destination path is required")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" {
		return ErrSourcePathEmpty
	}
	if toPath == "" {
		return ErrDestinationPathEmpty
	}

	sourceFileStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset < 0 {
		offset = 0
	}
	if offset > sourceFileStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	if limit <= 0 {
		limit = sourceFileStat.Size()
	}
	if limit > sourceFileStat.Size()-offset {
		limit = sourceFileStat.Size() - offset
	}

	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = sourceFile.Seek(offset, 0)
	if err != nil {
		return err
	}
	bar := progressbar.Default(limit)
	proxyReader := progressbar.NewReader(sourceFile, bar)
	_, err = io.CopyN(destinationFile, &proxyReader, limit)
	if err != nil {
		return err
	}

	return nil
}
