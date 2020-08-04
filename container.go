package dockerop

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

type ContainerDef struct {
	Image   string
	CName   string
	Cmd     strslice.StrSlice
	PortSet nat.PortSet   // 容器开放端口 List of exposed ports
	PortMap nat.PortMap   // 容器与主机端口映射 Port mapping between the exposed port (container) and the host
	Mounts  []mount.Mount // 挂载目录
}

// ContainerList return all of the containers
func (c *OpClient) ContainerList(ctx context.Context) ([]types.Container, error) {
	return c.client.ContainerList(ctx, types.ContainerListOptions{All: true})
}

// ContainerCreate create a container
func (c *OpClient) ContainerCreate(ctx context.Context, def *ContainerDef) (string, error) {
	resp, err := c.client.ContainerCreate(ctx, &container.Config{
		Image:        def.Image,   //镜像名称
		Tty:          true,        //docker run命令中的-t选项
		OpenStdin:    true,        //docker run命令中的-i选项
		Cmd:          def.Cmd,     //docker 容器中执行的命令
		ExposedPorts: def.PortSet, //端口
		//WorkingDir:   "/work",   //docker容器中的工作目录
	}, &container.HostConfig{
		PortBindings: def.PortMap,
		Mounts:       def.Mounts,
	}, &network.NetworkingConfig{}, nil, def.CName)
	return resp.ID, err
}

// ContainerStart start the container
func (c *OpClient) ContainerStart(ctx context.Context, containerID string) error {
	return c.client.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

// ContainerPause pause the container
func (c *OpClient) ContainerPause(ctx context.Context, containerID string) error {
	return c.client.ContainerPause(ctx, containerID)
}

// ContainerStop stops a container. In case the container fails to stop
// gracefully within a time frame specified by the timeout argument, it
// is forcefully terminated (killed).
// If the timeout is nil, the container's StopTimeout value is used, if
// set, otherwise the engine default. A negative timeout value can be
// specified, meaning no timeout, i.e. no forceful termination is performed.
func (c *OpClient) ContainerStop(ctx context.Context, containerID string, timeout *time.Duration) error {
	return c.client.ContainerStop(ctx, containerID, timeout)
}

// ContainerRemove remove the container
func (c *OpClient) ContainerRemove(ctx context.Context, containerID string) error {
	return c.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: true, Force: true})
}

// ContainersPrune prune the container source
//func (c *OpClient) ContainersPrune(ctx context.Context) error {
//	_, err := c.client.ContainersPrune(ctx, filters.Args{})
//	return err
//}

// ContainerInspect pause the container
func (c *OpClient) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	return c.client.ContainerInspect(ctx, containerID)
}

// ContainerLogs return the logs in container
func (c *OpClient) ContainerLogs(ctx context.Context, containerID string) (string, error) {
	reader, err := c.client.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Details: true})
	if err != nil {
		return "", err
	}
	info := new(bytes.Buffer)
	info.ReadFrom(reader)
	return info.String(), nil
}

// ContainerExec exec cmds in remote docker and return stdout/stderr
// cmd such as : []string{"ls", "-al"},
func (c *OpClient) ContainerExec(ctx context.Context, containerID string, cmd []string) (string, error) {
	idr, err := c.client.ContainerExecCreate(ctx, containerID, types.ExecConfig{AttachStdout: true, AttachStderr: true, Cmd: cmd})
	if err != nil {
		return "", err
	}
	resp, err := c.client.ContainerExecAttach(ctx, idr.ID, types.ExecStartCheck{})
	defer resp.Close()
	if err != nil {
		return "", err
	}
	// 去除头部乱码
	//for i := 0; i < 8; i++ {
	//	resp.Reader.ReadByte()
	//}
	var res string
	for {
		text, err := resp.Reader.ReadString('\n')
		res += text
		if err == io.EOF {
			break
		}
	}
	return res, nil
}
