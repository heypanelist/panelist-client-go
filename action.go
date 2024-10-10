package panelist

type Action interface {
	Serialize(context Context) interface{}
}
