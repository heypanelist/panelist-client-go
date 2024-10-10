package input

import (
	"github.com/heypanelist/panelist-client-go"
	"github.com/heypanelist/panelist-client-go/components/output"
)

type ButtonVariant string

const (
	ButtonVariantDefault ButtonVariant = "default"
	ButtonVariantOutline ButtonVariant = "outline"
	ButtonVariantGhost   ButtonVariant = "ghost"
)

type ButtonSize string

const (
	ButtonSizeFull ButtonSize = "full"
	ButtonSizeIcon ButtonSize = "icon"
)

type Button struct {
	Label     string
	LeftIcon  *output.Icon
	RightIcon *output.Icon
	Variant   *ButtonVariant
	Size      *ButtonSize
	Action    panelist.Action
}

func (b Button) ComponentName() string {
	return "button"
}

func (b *Button) Serialize(context panelist.Context) interface{} {
	return map[string]interface{}{
		"name": b.ComponentName(),
		"props": map[string]interface{}{
			"label": b.Label,
			"left_icon": func() interface{} {
				if b.LeftIcon == nil {
					return nil
				}
				return b.LeftIcon.Serialize(context)
			}(),
			"right_icon": func() interface{} {
				if b.RightIcon == nil {
					return nil
				}
				return b.RightIcon.Serialize(context)
			},
			"variant": b.Variant,
			"size":    b.Size,
			"action":  b.Action.Serialize(context),
		},
	}
}
