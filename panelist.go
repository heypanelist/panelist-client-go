// Package panelist is the main package in the Panelist client library.  It provides the Panelist
// struct which is the entry point for the library and houses other main types and
// functions used to interact with the Panelist server.
package panelist

import (
	"context"
	"fmt"

	wamp_client "github.com/gammazero/nexus/v3/client"
	"github.com/heypanelist/panelist-client-go/internal/client"
)

const clientVersion = "go-0.0.1"

// Panelist is the client that connects to the Panelist server
type Panelist struct {
	wampClient *wamp_client.Client
	config     Config

	pages []Page
}

type Config struct {
	// Name of this client.  Should be unique and in [kebab case](https://developer.mozilla.org/en-US/docs/Glossary/Kebab_case).
	Name string
	// WorkspaceSlug is the workspace on the server the client wants to join.
	WorkspaceSlug string
	// ServerPort is the port of your Panelist server.
	ServerPort int
	// Hostname is the host of your Panelist server.
	ServerHost string
}

// New creates a new instance of the Panelist client
func New(config Config) *Panelist {
	p := &Panelist{
		config: config,
	}
	return p
}

// AddPages adds the provided pages and makes them visible to
// the Panelist server.  This function should be called before Listen.
func (p *Panelist) AddPages(pages ...Page) error {
	p.pages = append(p.pages, pages...)
	return nil
}

// Listen listens for incoming messages and connects the client to the Panelist server
// using the configuration provided. It should be called last, after all pages have been registered.
// This function blocks the current goroutine.
func (p *Panelist) Listen() (err error) {
	ctx := context.Background()

	if p.wampClient, err = wamp_client.ConnectNet(ctx, fmt.Sprintf("tcp://%s:%d/", p.config.ServerHost, p.config.ServerPort), wamp_client.Config{
		Realm: "panelist",
	}); err != nil {
		return fmt.Errorf("panelist: failed to connect to server: %w", err)
	}

	pages := []string{}
	for _, page := range p.pages {
		pages = append(pages, page.Name)
	}

	response, err := send[client.RegisterRequest, client.RegisterResponse](ctx,
		p.wampClient,
		client.UriClientRegister,
		client.RegisterRequest{
			Name:          p.config.Name,
			ClientVersion: clientVersion,
			WorkspaceSlug: p.config.WorkspaceSlug,
			Pages:         pages,
		})
	if err != nil {
		return fmt.Errorf("panelist: failed to register client: %w", err)
	}
	if response.Error != nil {
		return fmt.Errorf("panelist: failed to register client: %w", *response.Error)
	}

	select {}
}
