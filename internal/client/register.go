package client

type RegisterRequest struct {
	Name          string   `json:"name"`
	ClientVersion string   `json:"client_version"`
	WorkspaceSlug string   `json:"workspace_slug"`
	Pages         []string `json:"pages"`
}

type RegisterResponse struct {
	ID string `json:"id"`
}
