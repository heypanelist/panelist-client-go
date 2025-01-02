package client

type RegisterRequest struct {
	Name          string `json:"name"`
	ClientVersion string `json:"client_version"`
}

type RegisterResponse struct {
	ID string `json:"id"`
}
