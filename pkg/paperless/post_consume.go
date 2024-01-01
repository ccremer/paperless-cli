package paperless

import (
	"time"
)

type PostConsumeDocumentOptions struct {
	Id               string
	FileName         string
	Created          time.Time
	Modified         time.Time
	Added            time.Time
	SourcePath       string
	ArchivePath      string
	ThumbnailPath    string
	DownloadUrl      string
	ThumbnailUrl     string
	Correspondent    string
	Tags             []string
	OriginalFileName string
}
