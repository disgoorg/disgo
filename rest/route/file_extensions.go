package route

// FileExtension is the type of image on Discord's CDN
type FileExtension string

// The available FileExtension(s)
const (
	PNG   FileExtension = "png"
	JPEG  FileExtension = "jpg"
	WEBP  FileExtension = "webp"
	GIF   FileExtension = "gif"
	BLANK FileExtension = ""
)

func (f FileExtension) String() string {
	return string(f)
}
