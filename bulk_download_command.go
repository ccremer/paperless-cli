package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ccremer/paperless-cli/pkg/archive"
	"github.com/ccremer/paperless-cli/pkg/errors"
	"github.com/ccremer/paperless-cli/pkg/localdb"
	"github.com/ccremer/paperless-cli/pkg/paperless"
	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

type BulkDownloadCommand struct {
	cli.Command

	PaperlessURL   string
	PaperlessToken string
	PaperlessUser  string

	TargetPath              string
	Content                 string
	UnzipEnabled            bool
	OverwriteExistingTarget bool
	Incremental             bool
}

const desc = `Use this command to create a local offline-copy of all documents.
If --%s is given, it will only download documents that don't exist locally, to save bandwidth.`

func newBulkDownloadCommand() *BulkDownloadCommand {
	c := &BulkDownloadCommand{}
	c.Command = cli.Command{
		Name:        "bulk-download",
		Usage:       "Downloads all documents at once",
		Description: fmt.Sprintf(desc, newIncrementalFlag(nil).Name),
		Before:      loadConfigFileFn,
		Action:      actions(LogMetadata, c.Action),
		Flags: []cli.Flag{
			newURLFlag(&c.PaperlessURL),
			newUsernameFlag(&c.PaperlessUser),
			newTokenFlag(&c.PaperlessToken),
			newTargetPathFlag(&c.TargetPath),
			newDownloadContentFlag(&c.Content),
			newUnzipFlag(&c.UnzipEnabled),
			newOverwriteFlag(&c.OverwriteExistingTarget),
			newIncrementalFlag(&c.Incremental),
		},
	}
	return c
}

func (c *BulkDownloadCommand) Action(ctx *cli.Context) error {
	log := logr.FromContextOrDiscard(ctx.Context)
	if c.Incremental {
		c.OverwriteExistingTarget = true
		c.UnzipEnabled = true
	}

	if prepareErr := c.prepareTarget(); prepareErr != nil {
		return prepareErr
	}
	clt := paperless.NewClient(c.PaperlessURL, c.PaperlessUser, c.PaperlessToken)

	log.Info("Getting list of documents")
	documents, queryErr := clt.QueryDocuments(ctx.Context, paperless.QueryParams{
		TruncateContent: true,
		Ordering:        "id",
		PageSize:        100,
	})
	if queryErr != nil {
		return queryErr
	}
	documentIDs := paperless.MapToDocumentIDs(documents)
	var db *localdb.Database

	if c.Incremental {
		log.V(1).Info("Opening DB", "dir", c.getTargetPath())
		newDb, openErr := localdb.Open(c.getTargetPath())
		if openErr != nil {
			return openErr
		}
		db = newDb
		newDocuments := c.filterMissingDocuments(db, documents)
		for _, doc := range newDocuments {
			db.Put(doc)
		}

		deletedDocuments := c.filterDeletedDocuments(db, paperless.MapToDocumentMap(documents))
		for _, deletedDoc := range deletedDocuments {
			db.Remove(deletedDoc)
		}
		if err := c.removeFiles(ctx, deletedDocuments); err != nil {
			return fmt.Errorf("cannot delete local documents: %w", err)
		}
		log.Info("Cleaned up deleted documents", "count", len(deletedDocuments))
		documentIDs = paperless.MapToDocumentIDs(newDocuments)
	}

	if len(documentIDs) == 0 {
		log.Info("Nothing to download")
		if db != nil {
			log.V(1).Info("Saving DB")
			return db.Close()
		}
		return nil
	}

	tmpFile, err := c.downloadDocuments(ctx, clt, documentIDs)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name()) // cleanup if not renamed

	if c.UnzipEnabled {
		unzipErr := c.unzip(ctx, tmpFile)
		if unzipErr != nil {
			return unzipErr
		}
		if db == nil {
			return nil
		}
		log.V(1).Info("Saving DB")
		return db.Close()
	}
	return c.move(ctx, tmpFile)
}

func (c *BulkDownloadCommand) removeFiles(ctx *cli.Context, deletedDocs []paperless.Document) error {
	log := logr.FromContextOrDiscard(ctx.Context)

	dir := c.getTargetPath()

	files := map[string]paperless.Document{}
	for _, doc := range deletedDocs {
		files[doc.ArchivedFileName] = doc
		files[doc.OriginalFileName] = doc
	}

	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		fileName := filepath.Base(path)
		if doc, found := files[fileName]; found {
			log.V(1).Info("Removing deleted document", "id", doc.ID, "path", path)
			_ = os.Remove(path)
		}
		return nil
	})
	return err
}

func (c *BulkDownloadCommand) downloadDocuments(ctx *cli.Context, clt *paperless.Client, documentIDs []int) (*os.File, error) {
	log := logr.FromContextOrDiscard(ctx.Context)

	tmpFile, createTempErr := os.CreateTemp(os.TempDir(), "paperless-bulk-download-")
	if createTempErr != nil {
		return nil, fmt.Errorf("cannot open temporary file: %w", createTempErr)
	}

	log.Info("Downloading documents", "count", len(documentIDs))
	downloadErr := clt.BulkDownload(ctx.Context, tmpFile, paperless.BulkDownloadParams{
		FollowFormatting: true,
		Content:          paperless.BulkDownloadContent(c.Content),
		DocumentIDs:      documentIDs,
	})
	return tmpFile, errors.Wrap(downloadErr, "could not download documents")
}

func (c *BulkDownloadCommand) unzip(ctx *cli.Context, tmpFile *os.File) error {
	log := logr.FromContextOrDiscard(ctx.Context)
	downloadFilePath := c.getTargetPath()
	if c.Content == paperless.BulkDownloadArchives.String() {
		downloadFilePath = filepath.Join(downloadFilePath, paperless.BulkDownloadArchives.String())
	}
	if c.Content == paperless.BulkDownloadOriginal.String() {
		downloadFilePath = filepath.Join(downloadFilePath, paperless.BulkDownloadOriginal.String())
	}
	if unzipErr := archive.Unzip(ctx.Context, tmpFile.Name(), downloadFilePath); unzipErr != nil {
		return fmt.Errorf("cannot unzip file %q to %q: %w", tmpFile.Name(), downloadFilePath, unzipErr)
	}
	log.Info("Unzipped archive to dir", "dir", downloadFilePath)
	return nil
}

func (c *BulkDownloadCommand) move(ctx *cli.Context, tmpFile *os.File) error {
	log := logr.FromContextOrDiscard(ctx.Context)
	downloadFilePath := c.getTargetPath()
	if renameErr := os.Rename(tmpFile.Name(), downloadFilePath); renameErr != nil {
		return fmt.Errorf("cannot move temp file: %w", renameErr)
	}
	log.Info("Downloaded zip archive", "file", downloadFilePath)
	return nil
}

func (c *BulkDownloadCommand) getTargetPath() string {
	if c.TargetPath != "" {
		return c.TargetPath
	}
	if c.UnzipEnabled {
		return "documents"
	}
	return "documents.zip"
}

func (c *BulkDownloadCommand) prepareTarget() error {
	target := c.getTargetPath()
	if c.OverwriteExistingTarget {
		if c.Incremental {
			return nil
		}
		return os.RemoveAll(target)
	}
	_, err := os.Stat(target)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return fmt.Errorf("target %q exists, abort", target)
}

func (c *BulkDownloadCommand) filterMissingDocuments(db *localdb.Database, documentsOnServer []paperless.Document) []paperless.Document {
	missing := make([]paperless.Document, 0)
	for i := 0; i < len(documentsOnServer); i++ {
		serverDoc := documentsOnServer[i]
		localDoc := db.FindByID(serverDoc.ID)
		if localDoc == nil {
			missing = append(missing, serverDoc)
		}
	}
	return missing
}

func (c *BulkDownloadCommand) filterDeletedDocuments(db *localdb.Database, documentsOnServer map[int]paperless.Document) []paperless.Document {
	extra := make([]paperless.Document, 0)
	allLocalDocs := db.GetAll()
	for i := 0; i < len(allLocalDocs); i++ {
		localDoc := allLocalDocs[i]
		if _, found := documentsOnServer[localDoc.ID]; !found {
			extra = append(extra, localDoc)
		}
	}
	return extra
}
