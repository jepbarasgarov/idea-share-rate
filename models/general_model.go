package models

import "mime/multipart"

type ParsedFile struct {
	File        multipart.File
	FileHeader  *multipart.FileHeader
	ContentType string
}

type LimitOffset struct {
	Limit  int
	Offset int
}
