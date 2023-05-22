package executor

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"server/internal/entity"
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

type BuildResult struct {
	Stream string `json:"stream"`
}

func getImageId(rd io.Reader) (string, error, int) {
	var lastLine string
	buildRes := &BuildResult{}

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		var line string = scanner.Text()
		if line[2:8] == "stream" {
			lastLine = line
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err, entity.PROCESS_SERVER_ERROR
	}

	json.Unmarshal([]byte(lastLine), buildRes)
	if len(buildRes.Stream) != 32 || len(buildRes.Stream) < 32 {

		return "", fmt.Errorf("%s\n", buildRes.Stream), entity.PROCESS_COMPILE_ERROR
	}

	return buildRes.Stream[19:31], nil, entity.PROCESS_OK
}

func (e *ExecutorCPP) ExecuteFromSource(source string) (output string, err error, errcode int) {
	goRoot := os.Getenv("GO_ROOT")
	if goRoot == "" {
		return "", fmt.Errorf("Could not find GO_ROOT variable in .env"), entity.PROCESS_SERVER_ERROR
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/source.cpp", goRoot), []byte(source), 0644)
	if err != nil {
		return "", err, entity.PROCESS_SERVER_ERROR
	}
	defer os.Remove(fmt.Sprintf("%s/source.cpp", goRoot))

	cli, err := client.NewEnvClient()
	if err != nil {
		return "", err, entity.PROCESS_SERVER_ERROR
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	tar, err := archive.TarWithOptions("../", &archive.TarOptions{})
	if err != nil {
		return "", err, entity.PROCESS_SERVER_ERROR
	}

	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "DockerfileCPP",
		Remove:     true,
	}

	res, err := cli.ImageBuild(ctx, tar, buildOptions)
	if err != nil {
		return "", err, entity.PROCESS_SERVER_ERROR
	}
	defer res.Body.Close()

	imageid, err, errcode := getImageId(res.Body)
	if errcode != 0 {
		return "", err, errcode
	}

	resp_create, err := cli.ContainerCreate(ctx, &container.Config{Image: imageid}, &container.HostConfig{}, nil, nil, "cpp_executor")
	if err != nil {
		return "", err, entity.PROCESS_SERVER_ERROR
	}

	if err := cli.ContainerStart(ctx, resp_create.ID, types.ContainerStartOptions{}); err != nil {
		return "", err, entity.PROCESS_SERVER_ERROR
	}
	defer cli.ContainerRemove(ctx, resp_create.ID, types.ContainerRemoveOptions{})

	reader, err := cli.ContainerLogs(ctx, resp_create.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", err, entity.PROCESS_SERVER_ERROR
	}

	logsAsBytes, err := io.ReadAll(reader)
	if err != nil {
		return "", nil, entity.PROCESS_SERVER_ERROR
	}

	result := string(logsAsBytes)
	result = strings.ReplaceAll(result, "\u0000", "")
	result = strings.ReplaceAll(result, "\u0001", "")
	result = strings.ReplaceAll(result, "\u0002", "")
	result = strings.Replace(result, "\r", "", 1)

	return result, nil, entity.PROCESS_OK
}
