package output

import (
	"fmt"

	"github.com/heypanelist/panelist-client-go"
)

type Icon struct {
	Src IconSrc
}

type IconSrcFrom string

const (
	IconSrcFromUrl    IconSrcFrom = "url"
	IconSrcFromLucide IconSrcFrom = "lucide"
)

type IconSrc struct {
	From  IconSrcFrom
	Value string
}

func (i Icon) ComponentName() string {
	return "icon"
}

func (i Icon) Serialize(context panelist.Context) interface{} {
	return map[string]interface{}{
		"name": i.ComponentName(),
		"props": map[string]interface{}{
			"src": fmt.Sprintf("%s-%s", i.Src.From, i.Src.Value),
		},
	}
}
