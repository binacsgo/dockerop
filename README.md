# dockerop
Docker client based on sdk, written in Go.

## interface

```go
// DockerOp interface of dockerop
type DockerOp interface {
	// client basic function:
	// NewOpClient create a new client of docker engine
	Close() error // close
	Ping(ctx context.Context) (types.Ping, error)
	Info(ctx context.Context) (types.Info, error)
	RegistryLogin(ctx context.Context, auth types.AuthConfig) error
	// client image function:
	ImageList(ctx context.Context) ([]types.ImageSummary, error)
	ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (string, error)
	ImageRemove(ctx context.Context, imageID string) error
	ImagesPrune(ctx context.Context) error
	// client container function:
	ContainerList(ctx context.Context) ([]types.Container, error)
	ContainerCreate(ctx context.Context, def *ContainerDef) (string, error)
	ContainerStart(ctx context.Context, containerID string) error
	ContainerPause(ctx context.Context, containerID string) error
	ContainerStop(ctx context.Context, containerID string, timeout *time.Duration) error
	ContainerRemove(ctx context.Context, containerID string) error
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
	ContainerLogs(ctx context.Context, containerID string) (string, error)
	//ContainersPrune(ctx context.Context) error
}
```