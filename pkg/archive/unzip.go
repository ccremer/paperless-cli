package archive

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"
)

// Unzip reads and copies every file in the archive to the destination dir.
func Unzip(ctx context.Context, source, dest string) error {
	log := logr.FromContextOrDiscard(ctx)
	log.V(1).Info("Unzipping file", "source", source, "dest", dest)
	archive, openErr := zip.OpenReader(source)
	if openErr != nil {
		return fmt.Errorf("cannot open source file: %w", openErr)
	}
	defer archive.Close()

	for _, f := range archive.File {
		destFilePath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(destFilePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", destFilePath)
		}
		if f.FileInfo().IsDir() {
			log.V(2).Info("Creating directory", "dir", f.FileInfo().Name())
			if mkdirErr := os.MkdirAll(destFilePath, os.ModePerm); mkdirErr != nil {
				return fmt.Errorf("cannot create directory: %w", mkdirErr)
			}
			continue
		}
		log.V(2).Info("Extracting file", "source", f.Name, "dest", destFilePath)

		err := unzipFile(f, destFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func unzipFile(f *zip.File, destFilePath string) error {
	// ensure directory exists where file should be written.
	if mkdirErr := os.MkdirAll(filepath.Dir(destFilePath), os.ModePerm); mkdirErr != nil {
		return fmt.Errorf("cannot create directory: %w", mkdirErr)
	}

	dstFile, dstFileErr := os.OpenFile(destFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if dstFileErr != nil {
		return fmt.Errorf("cannot open destination file: %w", dstFileErr)
	}
	defer dstFile.Close()

	fileInArchive, srcFileErr := f.Open()
	if srcFileErr != nil {
		return fmt.Errorf("cannot open source file: %w", srcFileErr)
	}
	fileInArchive.Close()

	if _, copyErr := io.Copy(dstFile, fileInArchive); copyErr != nil {
		return fmt.Errorf("cannot copy %q to %q: %w", f.Name, dstFile.Name(), copyErr)
	}
	return nil
}
