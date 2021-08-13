package interfaces

type HttpClient interface {
	Url(url string) error
	Post(body string) ([]byte, error)
}
