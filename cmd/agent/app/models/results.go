package models

type Result struct {
	data []byte
	err  error
}

func NewResult(data []byte, err error) Result {
	return Result{data: data, err: err}
}

func (r Result) Bytes() []byte {
	return r.data
}

func (r Result) HasError() bool {
	if r.err != nil {
		return false
	}
	return true
}

func (r Result) Error() error {
	return r.err
}
