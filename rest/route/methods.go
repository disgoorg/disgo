package route

// Method is a HTTP request Method
type Method string

// HTTP Methods used by Discord
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

// QueryValues is used to supply query param value pairs to Route.Compile
type QueryValues map[string]interface{}