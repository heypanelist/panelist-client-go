package panelist

import (
	"context"
	"encoding/json"
	"fmt"

	wamp_client "github.com/gammazero/nexus/v3/client"
	"github.com/gammazero/nexus/v3/wamp"
	"github.com/heypanelist/panelist-client-go/internal"
)

func send[I any, O any](
	ctx context.Context,
	wampClient *wamp_client.Client,
	uri internal.URI,
	data I,
) (*Response[O], error) {

	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("panelist: failed to process request: %w", err)
	}
	response, err := wampClient.Call(
		ctx,
		uri.String(),
		nil,
		nil,
		wamp.Dict{"data": string(bytes)},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("panelist: failed to send message: %w", err)
	}

	output := Response[O]{}
	str, ok := response.ArgumentsKw["data"].(string)
	if !ok {
		output.Body = nil
	} else {
		if err := json.Unmarshal([]byte(str), &output.Body); err != nil {
			return nil, fmt.Errorf("panelist: failed to process response: %w", err)
		}
	}

	return &output, nil
}

type Response[T any] struct {
	Success bool          `json:"success"`
	Error   *internal.Err `json:"error"`
	Body    *T            `json:"body"`
}
