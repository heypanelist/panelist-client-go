package panelist

type Component interface {
	Serialize(context Context) interface{}
	ComponentName() string
}
