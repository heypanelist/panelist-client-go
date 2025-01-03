package panelist

import (
	"context"
	"fmt"
	"regexp"

	wamp_client "github.com/gammazero/nexus/v3/client"
	"github.com/heypanelist/panelist-client-go/internal/client"
)

const clientVersion = "go-0.0.1"

// Panelist is the client that connects to the Panelist server
type Panelist struct {
	wampClient *wamp_client.Client
	config     Config

	pages      []Page
	pageGroups []PageGroup
}

type Config struct {
	// Name of this client.  Should be unique and in [kebab case](https://developer.mozilla.org/en-US/docs/Glossary/Kebab_case).
	Name string
	// Workspace is the workspace on the server the client wants to join.
	Workspace string
	// Port number of the Panelist server.
	ServerPort int
	// Hostname of the Panelist server.
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

// AddPageGroups adds the provided page groups and makes them visible to
// the Panelist server.  This function should be called before Listen.
func (p *Panelist) AddPageGroups(pages ...Page) error {
	p.pages = append(p.pages, pages...)
	return nil
}

// Listen listens for incoming messages and connects the client to the Panelist server
// using the configuration provided. It should be called last, after all pages have been registered.
// This function blocks the current goroutine.
func (p *Panelist) Listen() (err error) {
	ctx := context.Background()

	clientNameRegex := regexp.MustCompile(`^[a-z0-9-]{3,20}$`)
	if !clientNameRegex.MatchString(p.config.Name) {
		return fmt.Errorf(
			"panelist: client name must be between 3 and 64 characters long and match the kebab-case format",
		)
	}
	for _, page := range p.pages {
		if err := page.Validate(); err != nil {
			return fmt.Errorf("panelist: one or more of your pages have the error '%w'", err)
		}
	}

	if p.wampClient, err = wamp_client.ConnectNet(ctx, fmt.Sprintf("tcp://%s:%d/", p.config.ServerHost, p.config.ServerPort), wamp_client.Config{
		Realm: "panelist",
	}); err != nil {
		return fmt.Errorf("panelist: failed to connect to server: %w", err)
	}

	response, err := send[client.RegisterRequest, client.RegisterResponse](ctx,
		p.wampClient,
		client.UriClientRegister,
		client.RegisterRequest{
			Name:          p.config.Name,
			ClientVersion: clientVersion,
		})
	if err != nil {
		return fmt.Errorf("panelist: failed to register client: %w", err)
	}
	if response.Error != nil {
		return fmt.Errorf("panelist: failed to register client: %w", *response.Error)
	}

	select {}
}
