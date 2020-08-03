package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/binacsgo/dockerop"
)

const prompt = `
========> Please switch the function you want to test:
	ping:     cli.Ping
	info:     cli.Info
	ilist:    cli.ImageList
	ipull:    cli.ImagePull
	iremove:  cli.ImageRemove
	iprune:   cli.ImagesPrune
	clist:    cli.ContainerList
	ccreate:  cli.ContainerCreate
	cstart:   cli.ContainerStart
	cpause:   cli.ContainerPause
	cstop:    cli.ContainerStop
	cremove:  cli.ContainerRemove
	cinspect: cli.ContainerInspect
	clogs:    cli.ContainerLogs`

var (
	fcode string
	host  string
)

func init() {
	flag.StringVar(&host, "host", client.DefaultDockerHost, "Host location")
}

func handle(resp interface{}, err error) {
	handleErr(err)
	log.Printf("get %+v\n", resp)
}

func handleErr(err error) {
	if err != nil {
		log.Printf("err %+v\n", err)
	}
}

func main() {
	flag.Parse()
	cli, err := dockerop.NewOpClient(host, nil)
	if err != nil {
		log.Printf("new client err %+v\n", err)
	}
	log.Printf("client created success!\n")

	for {
		fmt.Println(prompt)
		fmt.Scanln(&fcode)
		doTest(cli, fcode)
	}

}

func doTest(cli *dockerop.OpClient, fcode string) {
	ctx := context.Background()
	var ref, imageID, image, cname, containerID string
	switch fcode {
	case "ping":
		handle(cli.Ping(ctx))
	case "info":
		handle(cli.Info(ctx))
	case "ilist":
		handle(cli.ImageList(ctx))
	case "ipull":
		fmt.Println("Input ref url (such as docker.io/library/nginx:latest):")
		fmt.Scanln(&ref)
		handle(cli.ImagePull(ctx, ref, types.ImagePullOptions{}))
	case "iremove":
		fmt.Println("Input imageID:")
		fmt.Scanln(&imageID)
		handleErr(cli.ImageRemove(ctx, imageID))
	case "iprune":
		cli.ImagesPrune(ctx)
	case "clist":
		handle(cli.ContainerList(ctx))
	case "ccreate":
		fmt.Println("Input image name:")
		fmt.Scanln(&image)
		fmt.Println("Input container name:")
		fmt.Scanln(&cname)
		handle(cli.ContainerCreate(ctx, &dockerop.ContainerDef{
			Image: image,
			CName: cname,
			PortSet: nat.PortSet{
				"80": struct{}{}, //docker容器对外开放的端口
			},
			PortMap: nat.PortMap{
				"80": []nat.PortBinding{nat.PortBinding{
					HostIP:   "0.0.0.0", //docker 容器映射的宿主机的ip
					HostPort: "80",      //docker 容器映射到宿主机的端口
				}},
			},
			Mounts: []mount.Mount{ //docker 容器目录挂在到宿主机目录
				mount.Mount{
					Type:   mount.TypeBind,
					Source: "/home/",
					Target: "/work/volume",
				},
			},
		}))
	case "cstart":
		fmt.Println("Input containerID:")
		fmt.Scanln(&containerID)
		handleErr(cli.ContainerStart(ctx, containerID))
	case "cpause":
		fmt.Println("Input containerID:")
		fmt.Scanln(&containerID)
		handleErr(cli.ContainerPause(ctx, containerID))
	case "cstop":
		fmt.Println("Input containerID:")
		fmt.Scanln(&containerID)
		handleErr(cli.ContainerStop(ctx, containerID, nil))
	case "cremove":
		fmt.Println("Input containerID:")
		fmt.Scanln(&containerID)
		handleErr(cli.ContainerRemove(ctx, containerID))
	case "cinspect":
		fmt.Println("Input containerID:")
		fmt.Scanln(&containerID)
		handle(cli.ContainerInspect(ctx, containerID))
	case "clogs":
		fmt.Println("Input container or containerID:")
		fmt.Scanln(&cname)
		handle(cli.ContainerLogs(ctx, cname))
	default:
		fmt.Println("not support:", fcode)
	}
}
