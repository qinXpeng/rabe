package rabes

type Status int

const (
	STATUS_NOT_FOUND Status = 404
	STATUS_PAGE_OK   Status = 200
)

var statusText = map[Status]string{
	STATUS_NOT_FOUND: "404 NOT FOUND",
	STATUS_PAGE_OK:   "200 OK",
}

func (s Status) String() string {
	return statusText[s]
}

func (s Status) Error() string {
	return s.String()
}
func (s Status) Code() int {
	return int(s)
}
