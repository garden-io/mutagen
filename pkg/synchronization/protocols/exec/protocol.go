package exec

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/mutagen-io/mutagen/pkg/agent"
	"github.com/mutagen-io/mutagen/pkg/agent/transport/exec"
	"github.com/mutagen-io/mutagen/pkg/logging"
	"github.com/mutagen-io/mutagen/pkg/synchronization"
	"github.com/mutagen-io/mutagen/pkg/synchronization/endpoint/remote"
	urlpkg "github.com/mutagen-io/mutagen/pkg/url"
)

// protocolHandler implements the synchronization.ProtocolHandler interface for
// connecting to remote endpoints by running a shell command. It uses the agent
// infrastructure over a custom command transport.
type protocolHandler struct{}

// dialResult provides asynchronous agent dialing results.
type dialResult struct {
	// stream is the stream returned by agent dialing.
	stream io.ReadWriteCloser
	// error is the error returned by agent dialing.
	error error
}

// Connect connects to an exec endpoint.
func (h *protocolHandler) Connect(
	ctx context.Context,
	logger *logging.Logger,
	url *urlpkg.URL,
	prompter string,
	session string,
	version synchronization.Version,
	configuration *synchronization.Configuration,
	alpha bool,
) (synchronization.Endpoint, error) {
	// Verify that the URL is of the correct kind and protocol.
	if url.Kind != urlpkg.Kind_Synchronization {
		panic("non-synchronization URL dispatched to synchronization protocol handler")
	} else if url.Protocol != urlpkg.Protocol_Exec {
		panic("non-Exec URL dispatched to Exec protocol handler")
	}

	// Ensure that no environment variables or parameters are specified. These
	// are neither expected nor supported for SSH URLs.
	if len(url.Environment) > 0 {
		return nil, errors.New("exec URL contains environment variables")
	} else if len(url.Parameters) > 0 {
		return nil, errors.New("exec URL contains internal parameters")
	}

	// Create an exec-based agent transport.
	transport, err := exec.NewTransport(url.Host, prompter)
	if err != nil {
		return nil, fmt.Errorf("unable to create exec transport: %w", err)
	}

	// Create a channel to deliver the dialing result.
	results := make(chan dialResult)

	// Perform dialing in a background Goroutine so that we can monitor for
	// cancellation.
	go func() {
		// Perform the dialing operation.
		stream, err := agent.Dial(logger, transport, agent.CommandSynchronizer, prompter)

		// Transmit the result or, if cancelled, close the stream.
		select {
		case results <- dialResult{stream, err}:
		case <-ctx.Done():
			if stream != nil {
				stream.Close()
			}
		}
	}()

	// Wait for dialing results or cancellation.
	var stream io.ReadWriteCloser
	select {
	case result := <-results:
		if result.error != nil {
			return nil, fmt.Errorf("unable to dial agent endpoint: %w", result.error)
		}
		stream = result.stream
	case <-ctx.Done():
		return nil, errors.New("connect operation cancelled")
	}

	// Create the endpoint client.
	return remote.NewEndpoint(logger, stream, url.Path, session, version, configuration, alpha)
}

func init() {
	// Register the exec protocol handler with the synchronization package.
	synchronization.ProtocolHandlers[urlpkg.Protocol_Exec] = &protocolHandler{}
}
