package panelist

import "github.com/heypanelist/panelist-client-go/common"

type Component interface {
	Serialize(context common.Context) interface{}
	ComponentName() string
}
