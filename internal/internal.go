package internal

type URI string

func (u URI) String() string {
	return string(u)
}

type Err string

func (e Err) Error() string {
	return string(e)
}
