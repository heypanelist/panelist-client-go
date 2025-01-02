package panelist

import "github.com/heypanelist/panelist-client-go/common"

type Action interface {
	Serialize(context common.Context) interface{}
	ActionName() string
}
