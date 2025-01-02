package panelist

import (
	"errors"

	"github.com/heypanelist/panelist-client-go/common"
	"github.com/heypanelist/panelist-client-go/components/output"
)

type PageGroup struct {
	// Name is the unique identifier for the page group.  It is used to determine the URL of the page group and its childrens.
	//
	// required: true
	Name string

	// Handler returns a slice of pages that belong to the page group.
	//
	// required: true
	Handler func(context common.Context) []Page

	// Title is the human-readable name of the page group.  If not provided, the Name will be used.
	//
	// required: false
	Title string
	// Icon is the icon that will be displayed in the navigation for the page group.
	//
	// required: false
	Icon *output.Icon

	// Unlist is a boolean that determines whether or not the page group should be displayed in the navigation.
	//
	// required: false, default: false
	Unlist bool

	// UnlistFunc is a function that determines whether or not the page group should be displayed in the navigation.
	//
	// required: false
	UnlistFunc func(context common.Context) bool

	// Hidden is a boolean that determines whether or not the page group and its pages should be accessible by URL.
	// Setting hidden will automatically set unlist to true.
	//
	// required: false, default: false
	Hidden bool

	// HiddenFunc is a function that determines whether or not the page group and its pages should be accessible by URL.
	// Setting hidden will automatically set unlist to true.
	//
	// required: false
	HiddenFunc func(context common.Context) bool
}

func (pg *PageGroup) Serialize(context common.Context) interface{} {

	pages := pg.Handler(context)
	if pages == nil || len(pages) == 0 {
		return nil
	}

	return map[string]interface{}{
		"name": pg.Name,
		"title": (func() string {
			if pg.Title != "" {
				return pg.Title
			}
			return pg.Name
		})(),
		"icon": (func() interface{} {
			if pg.Icon != nil {
				return pg.Icon.Serialize(context)
			}
			return nil
		})(),
		"pages": (func() []interface{} {
			serializedPages := make([]interface{}, len(pages))
			for i, page := range pages {
				serializedPages[i] = page.Serialize(context)
			}
			return serializedPages
		})(),
	}
}

func (pg *PageGroup) GetPageGroupListItem(context common.Context) map[string]interface{} {
	if err := pg.Validate(); err != nil || pg.IsUnlisted(context) {
		return nil
	}

	pageListItems := []map[string]interface{}{}
	for _, page := range pg.Handler(context) {
		pageListItem := page.GetPageListItem(context)
		if pageListItem != nil {
			pageListItems = append(pageListItems, pageListItem)
		}
	}
	if len(pageListItems) == 0 {
		return nil
	}

	return map[string]interface{}{
		"name": pg.Name,
		"title": (func() string {
			if pg.Title != "" {
				return pg.Title
			}
			return pg.Name
		})(),
		"pages": pageListItems,
		"icon": (func() interface{} {
			if pg.Icon != nil {
				return pg.Icon.Serialize(context)
			}
			return nil
		})(),
	}
}

func (pg *PageGroup) IsHidden(context common.Context) bool {
	if pg.HiddenFunc != nil && pg.HiddenFunc(context) {
		return true
	}
	if pg.Hidden {
		return true
	}
	return false
}

func (pg *PageGroup) Validate() error {
	if pg.Name == "" {
		return errors.New("missing page name")
	}
	if pg.Handler == nil {
		return errors.New("missing page handler")
	}
	return nil
}

func (pg *PageGroup) IsUnlisted(context common.Context) bool {
	if pg.IsHidden(context) {
		return true
	}
	if pg.UnlistFunc != nil && pg.UnlistFunc(context) {
		return true
	}
	if pg.Unlist {
		return true
	}
	return false
}
