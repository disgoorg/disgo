package discord

import (
	"fmt"
	"regexp"
)

// MaxFileTypes is the maximum number of file types that can be filtered for at once.
// https://docs.discord.com/developers/reference#file-type-filtering
const MaxFileTypes = 10

var fileTypeExtensionPattern = regexp.MustCompile(`^\.[\w](?:[\w.-]*[\w])?$`)

// NewFileType creates a new FileType for a custom file extension, e.g. NewFileType(".pdf").
// The extension must be dot-prefixed and only contain latin letters, digits, dashes or dots.
// Use the FileTypeImage, FileTypeVideo or FileTypeAudio constants for the preset file groups.
func NewFileType(extension string) (FileType, error) {
	if !fileTypeExtensionPattern.MatchString(extension) {
		return "", fmt.Errorf(`file type extension %q must be dot-prefixed (e.g. ".pdf") and only contain latin letters, digits, dashes or dots`, extension)
	}
	return FileType(extension), nil
}

// FileType represents a type of file to filter for in FileUploadComponents and ApplicationCommandOptionAttachment.
// It can be a preset group like FileTypeImage, FileTypeVideo or FileTypeAudio, or a custom file extension like ".pdf".
// File types only match against the file extension, Discord does not validate the actual content of the file.
// Up to [MaxFileTypes] file types can be specified, each file type is validated when marshaled.
// https://docs.discord.com/developers/reference#file-type-filtering
type FileType string

// Preset FileType groups.
const (
	FileTypeImage FileType = "image"
	FileTypeVideo FileType = "video"
	FileTypeAudio FileType = "audio"
)

func (t FileType) String() string {
	return string(t)
}

func (t FileType) isValid() bool {
	return t == FileTypeImage || t == FileTypeVideo || t == FileTypeAudio || fileTypeExtensionPattern.MatchString(string(t))
}
