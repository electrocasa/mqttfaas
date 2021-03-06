package containerrunner

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http/httputil"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/sks/mqttfaas/pkg/config"
	"github.com/sks/mqttfaas/pkg/types"
)

var attachOptions = dockertypes.ContainerAttachOptions{
	Stream: true,
	Stdin:  true,
	Stdout: true,
	Stderr: true,
}

//ContainerRunner ...
type ContainerRunner struct {
	dockerCLI     DockerCLI
	configuration *config.Config
	defaultEnv    []string
}

//New ...
func New(dockerCLI DockerCLI, configuration *config.Config) *ContainerRunner {
	return &ContainerRunner{
		dockerCLI,
		configuration,
		GetDefaultEnvs(),
	}
}

//Run ...
func (c *ContainerRunner) Run(ctx context.Context, input *types.ImageRunnerInput) ([]byte, error) {
	containerID, err := c.createContainer(ctx, input)
	if err != nil {
		return nil, err
	}
	defer func() {
		if c.configuration.DontUseHotContainers || input.FunctionMetadata.DeleteAfterUse {
			c.dockerCLI.ContainerRemove(ctx, containerID, dockertypes.ContainerRemoveOptions{})
		}
	}()

	hijackResponse, err := c.dockerCLI.ContainerAttach(ctx, containerID, attachOptions)
	if err != nil && err != httputil.ErrPersistEOF {
		return nil, err
	}
	defer hijackResponse.Close()

	err = c.dockerCLI.ContainerStart(ctx, containerID, dockertypes.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	errCh := make(chan error)

	functionOutput := &bytes.Buffer{}
	go func() {
		errCh <- func() error {
			streamer := hijackedIOStreamer{
				inputStream:  ioutil.NopCloser(bytes.NewReader(input.Message)),
				outputStream: functionOutput,
				errorStream:  newErrorStream(input.FunctionMetadata.Name),
				resp:         hijackResponse,
			}

			return streamer.stream(ctx)
		}()
	}()
	select {
	case <-ctx.Done():
		return functionOutput.Bytes(), errors.New("Aborted")

	case err = <-errCh:
		return functionOutput.Bytes(), err
	}
}
