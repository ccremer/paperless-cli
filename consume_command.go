package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/ccremer/clustercode/pkg/consumer"
	"github.com/ccremer/clustercode/pkg/paperless"
	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

type ConsumeCommand struct {
	cli.Command

	PaperlessURL   string
	PaperlessToken string
	PaperlessUser  string

	ConsumeDirName string
	ConsumeDelay   time.Duration
}

func newConsumeCommand() *ConsumeCommand {
	c := &ConsumeCommand{}
	c.Command = cli.Command{
		Name:   "consume",
		Usage:  "Consumes a local directory and uploads each file to Paperless instance. The files will be deleted once uploaded.",
		Before: loadConfigFileFn,
		Action: actions(LogMetadata, c.Action),

		Flags: []cli.Flag{
			newURLFlag(&c.PaperlessURL),
			newUsernameFlag(&c.PaperlessUser),
			newTokenFlag(&c.PaperlessToken),
			newConsumeDirFlag(&c.ConsumeDirName),
			newConsumeDelayFlag(&c.ConsumeDelay),
		},
	}
	return c
}

func (c *ConsumeCommand) Action(ctx *cli.Context) error {
	log := logr.FromContextOrDiscard(ctx.Context)
	log.Info("Start consuming directory", "dir", c.ConsumeDirName)

	clt := paperless.NewClient(c.PaperlessURL, c.PaperlessUser, c.PaperlessToken)
	q := consumer.NewQueue[string]()
	q.Subscribe(ctx.Context, func(fileName string) {
		log.V(1).Info("Uploading file...", "file", fileName)
		err := clt.Upload(ctx.Context, fileName, paperless.UploadParams{})
		if err != nil {
			log.Error(err, "Could not upload file")
			return
		}
		if deleteErr := os.Remove(fileName); deleteErr != nil {
			log.Error(err, "Could not delete file, this might be re-uploaded later again", "file", fileName)
		}
		log.Info("File uploaded", "file", fileName)
	})

	walkErr := filepath.WalkDir(c.ConsumeDirName, func(path string, entry fs.DirEntry, err error) error {
		if path == c.ConsumeDirName {
			return nil // same directory, not interesting
		}
		if entry.IsDir() {
			return fs.SkipDir
		}
		if err != nil {
			return fs.SkipDir
		}
		q.Put(path)
		return nil
	})
	if walkErr != nil {
		return fmt.Errorf("cannot walk consumption dir: %w", walkErr)
	}

	watchErr := consumer.StartWatchingDir(ctx.Context, c.ConsumeDirName, c.ConsumeDelay, func(filePath string) {
		q.Put(filePath)
	})
	if watchErr != nil {
		return fmt.Errorf("cannot watch consumption dir: %w", watchErr)
	}
	<-make(chan struct{})
	return nil
}
