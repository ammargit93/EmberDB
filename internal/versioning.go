package internal

import (
	"mime/multipart"
	"time"
)

type FileVersion struct {
	file      multipart.File
	version   int
	timestamp time.Time
}
