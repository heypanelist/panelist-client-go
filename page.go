package panelist

import (
	"errors"

	"github.com/heypanelist/panelist-client-go/common"
	"github.com/heypanelist/panelist-client-go/components/output"
)

// Page represents a page in the Panelist interface.
type Page struct {
	// Name is the unique identifier for the page.  It is used to determine the URL of the page.
	//
	// required: true
	Name string

	// ComponentHandler is a function that returns the components that should be displayed on the page.
	//
	// required: true
	Handler func(context common.Context) []Component

	// Title is a human-readable name of the page.  If not provided, the Name will be used.
	//
	// required: false
	Title string

	// Unlist is a boolean that determines whether or not the page should be displayed in the navigation.
	//
	// required: false, default: false
	Unlist bool

	// UnlistFunc is a function that determines whether or not the page should be displayed in the navigation.
	//
	// required: false
	UnlistFunc func(context common.Context) bool

	// Hidden is a boolean that determines whether or not the page should be accessible at all.
	// Setting hidden will automatically set unlist to true.
	//
	// required: false, default: false
	Hidden bool

	// HiddenFunc is a function that determines whether or not the page should be accessible at all.
	// Setting hidden will automatically set unlist to true.
	//
	// required: false
	HiddenFunc func(context common.Context) bool

	// Icon is the icon that will be displayed in the navigation for the page.
	//
	// required: false
	Icon *output.Icon
}

func (p *Page) Serialize(context common.Context) interface{} {

	components := p.Handler(context)
	if components == nil {
		components = []Component{}
	}

	if p.UnlistFunc != nil {
		if p.UnlistFunc(context) {
			return nil
		}
	}
	if p.Unlist {
		return nil
	}

	return map[string]interface{}{
		"name": p.Name,
		"title": (func() string {
			if p.Title != "" {
				return p.Title
			}
			return p.Name
		})(),
		"icon": (func() interface{} {
			if p.Icon != nil {
				return p.Icon.Serialize(context)
			}
			return nil
		})(),
		"components": (func() []interface{} {
			serializedComponents := make([]interface{}, len(components))
			for i, component := range components {
				serializedComponents[i] = component.Serialize(context)
			}
			return serializedComponents
		})(),
	}
}

func (p *Page) GetPageListItem(context common.Context) map[string]interface{} {
	if err := p.Validate(); err != nil || p.IsUnlisted(context) {
		return nil
	}

	return map[string]interface{}{
		"name": p.Name,
		"title": (func() string {
			if p.Title != "" {
				return p.Title
			}
			return p.Name
		})(),
		"icon": (func() interface{} {
			if p.Icon != nil {
				return p.Icon.Serialize(context)
			}
			return nil
		})(),
	}
}

func (p *Page) IsHidden(context common.Context) bool {
	if p.HiddenFunc != nil && p.HiddenFunc(context) {
		return true
	}
	if p.Hidden {
		return true
	}
	return false
}

func (p *Page) IsUnlisted(context common.Context) bool {
	if p.IsHidden(context) {
		return true
	}
	if p.UnlistFunc != nil && p.UnlistFunc(context) {
		return true
	}
	if p.Unlist {
		return true
	}
	return false
}

func (p *Page) Validate() error {
	if p.Name == "" {
		return errors.New("missing page name")
	}
	if p.Handler == nil {
		return errors.New("missing page handler")
	}
	return nil
}
