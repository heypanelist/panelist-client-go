package block

import (
	"github.com/heypanelist/panelist-client-go/common"
	"github.com/heypanelist/panelist-client-go/components/input"
)

type PageHeader struct {
	Title         string
	Subtitle      *string
	ActionButtons []input.Button
}

func (ph PageHeader) ComponentName() string {
	return "page-header"
}

func (ph PageHeader) Serialize(context common.Context) interface{} {
	serializedButtons := make([]interface{}, len(ph.ActionButtons))
	for i, button := range ph.ActionButtons {
		serializedButtons[i] = button.Serialize(context)
	}

	return map[string]interface{}{
		"name": ph.ComponentName(),
		"props": map[string]interface{}{
			"title":          ph.Title,
			"subtitle":       ph.Subtitle,
			"action_buttons": serializedButtons,
		},
	}
}
