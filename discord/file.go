package discord

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"

	"github.com/DisgoOrg/disgo/json"
)

type Payload interface {
	ToBody() (interface{}, error)
}

// MultipartBuffer holds the Body & ContentType of the multipart body
type MultipartBuffer struct {
	Buffer      *bytes.Buffer
	ContentType string
}

// PayloadWithFiles returns the given payload as multipart body with all files in it
//goland:noinspection GoUnusedExportedFunction
func PayloadWithFiles(v interface{}, files ...*File) (buffer *MultipartBuffer, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.FormDataContentType()
	defer func() {
		_ = writer.Close()
	}()

	var payload []byte
	payload, err = json.Marshal(v)
	if err != nil {
		return
	}

	part, err := writer.CreatePart(partHeader(`form-data; name="payload_json"`, "application/json"))
	if err != nil {
		return
	}

	if _, err = part.Write(payload); err != nil {
		return
	}

	for i, file := range files {
		var name string
		if file.Flags.Has(FileFlagSpoiler) {
			name = "SPOILER_" + file.Name
		} else {
			name = file.Name
		}
		part, err = writer.CreatePart(partHeader(fmt.Sprintf(`form-data; name="file%d"; filename="%s"`, i, name), "application/octet-stream"))
		if err != nil {
			return
		}

		if _, err = io.Copy(part, file.Reader); err != nil {
			return
		}
	}

	buffer = &MultipartBuffer{
		Buffer:      body,
		ContentType: writer.FormDataContentType(),
	}
	return
}

func partHeader(contentDisposition string, contentType string) textproto.MIMEHeader {
	return textproto.MIMEHeader{
		"Content-Disposition": []string{contentDisposition},
		"Content-Type":        []string{contentType},
	}
}

// NewFile returns a new File struct with the given name, io.Reader & FileFlags
//goland:noinspection GoUnusedExportedFunction
func NewFile(name string, reader io.Reader, flags ...FileFlags) *File {
	return &File{
		Name:   name,
		Reader: reader,
		Flags:  FileFlagNone.Add(flags...),
	}
}

// File holds all information about a given io.Reader
type File struct {
	Name   string
	Reader io.Reader
	Flags  FileFlags
}

// FileFlags are used to mark Attachments as Spoiler
type FileFlags int

// all FileFlags
const (
	FileFlagSpoiler FileFlags = 1 << iota
	FileFlagNone    FileFlags = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (f FileFlags) Add(bits ...FileFlags) FileFlags {
	for _, bit := range bits {
		f |= bit
	}
	return f
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f FileFlags) Remove(bits ...FileFlags) FileFlags {
	for _, bit := range bits {
		f &^= bit
	}
	return f
}

// Has will ensure that the bit includes all the bits entered
func (f FileFlags) Has(bits ...FileFlags) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return false
		}
	}
	return true
}

// Missing will check whether the bit is missing any one of the bits
func (f FileFlags) Missing(bits ...FileFlags) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return true
		}
	}
	return false
}
