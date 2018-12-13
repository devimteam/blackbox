package blackbox

type Error struct {
	HttpStatus  int
	RawResponse string
}

func (e Error) Error() string {
	return e.RawResponse
}
