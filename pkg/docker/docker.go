package docker

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"go.uber.org/zap"
)

var (
	logger       *zap.Logger
	dockerClient *client.Client
)

func CheckAndRunDockerFuncs(imageName string, logger *zap.Logger) error {

	svc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	setClientAndLogger(svc, logger)

	if !checkForImage(imageName) {
		pullImage(imageName)
	}

	ok, containerId := checkForContainer(imageName)
	if !ok {
		containerId = createContainer(imageName)
	}

	startContainer(containerId)

	logger.Info("docker steps complete")

	defer dockerClient.Close()

	return nil

}

func checkForImage(imageName string) bool {

	logger.Info("fetching image list")

	imageList, err := dockerClient.ImageList(context.TODO(), types.ImageListOptions{
		All: true,
	})
	if err != nil {
		logger.Panic("failed to fetch image list", zap.Error(err))
	}

	logger.Info("fetch completed, looking for image")

	for i := range imageList {
		if len(imageList[i].RepoDigests) > 0 {
			if strings.Contains(imageList[i].RepoDigests[0], imageName) {
				logger.Info("image found")
				return true
			}
		}
	}

	logger.Info("couldn't find image")
	return false
}

func checkForContainer(imageName string) (bool, string) {

	logger.Info("fetching container list")

	containers, err := dockerClient.ContainerList(context.TODO(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	logger.Info("fetch completed, looking for container built on image")

	for i := range containers {
		if containers[i].Image == imageName {
			logger.Info("container found")
			return true, containers[i].ID
		}
	}

	logger.Info("couldn't find container")
	return false, ""

}

func createContainer(imageName string) string {

	out, err := dockerClient.ContainerCreate(context.TODO(), &container.Config{
		Image: imageName,
		Tty:   false,
		ExposedPorts: nat.PortSet{
			"9000/tcp": struct{}{},
		},
	},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"9000/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "9000",
					},
				},
			},
		}, nil, nil, "clamav-scanner")
	if err != nil {
		panic(err)
	}

	logger.Info("container created")

	return out.ID
}

func pullImage(imageName string) {

	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}

	logger.Info("attempting to pull image")

	out, err := dockerClient.ImagePull(context.TODO(), imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer os.Remove(f.Name())
	defer f.Close()
	defer out.Close()

	io.Copy(f, out)

	logger.Info("image pull completed")
}

func startContainer(containerId string) {

	logger.Info("starting container")

	if err := dockerClient.ContainerStart(context.TODO(), containerId, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	logger.Info("container started")
}

func setClientAndLogger(svc *client.Client, zapLogger *zap.Logger) {
	dockerClient, logger = svc, zapLogger
}
