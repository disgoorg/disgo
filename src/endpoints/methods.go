package endpoints

type Method string

const (
	DELETE Method = "DELETE"
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
)

func (m Method) String() string {
	return string(m)
}
