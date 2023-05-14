package executor

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type ExecutorCPP struct{}

func NewExecutorCPP() *ExecutorCPP {
	return &ExecutorCPP{}
}

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

type BuildResult struct {
	Stream string `json:"stream"`
}

func getImageId(rd io.Reader) (string, error) {
	var lastLine string
	buildRes := &BuildResult{}

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	json.Unmarshal([]byte(lastLine), buildRes)
	if len(buildRes.Stream) != 32 || len(buildRes.Stream) < 32 {
		return "", fmt.Errorf("Build failed")
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return "", fmt.Errorf("Error occured")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return buildRes.Stream[19:31], nil
}

func (e *ExecutorCPP) ExecuteFromSource(source string) (output string, err error) {
	goRoot := os.Getenv("GO_ROOT")
	if goRoot == "" {
		return "", fmt.Errorf("Could not find GO_ROOT variable in .env")
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/source.cpp", goRoot), []byte(source), 0644)
	if err != nil {
		return "", err
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	tar, err := archive.TarWithOptions("../", &archive.TarOptions{})
	if err != nil {
		return "", err
	}

	buildOptions := types.ImageBuildOptions{
		Dockerfile:  "DockerfileCPP",
		Remove:      true,
		ForceRemove: true,
		NoCache:     true,
	}

	res, err := cli.ImageBuild(ctx, tar, buildOptions)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	imageid, err := getImageId(res.Body)
	if err != nil {
		return "", err
	}

	resp_create, err := cli.ContainerCreate(ctx, &container.Config{Image: imageid}, &container.HostConfig{}, nil, nil, "cpp_executor")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp_create.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	reader, err := cli.ContainerLogs(ctx, resp_create.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", err
	}

	logsAsBytes, err := io.ReadAll(reader)
	if err != nil {
		return "", nil
	}

	err = cli.ContainerRemove(ctx, resp_create.ID, types.ContainerRemoveOptions{})
	if err != nil {
		return "", err
	}

	result := string(logsAsBytes)
	result = strings.ReplaceAll(result, "\u0000", "")
	result = strings.ReplaceAll(result, "\u0001", "")
	result = strings.ReplaceAll(result, "\u0002", "")
	result = strings.Replace(result, "\r", "", 1)

	return result, nil
}
