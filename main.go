package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"net/http"
)

type RestartResponse struct {
	Ok      bool
	Message string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res := RestartResponse{
			Ok:      true,
			Message: "Healthy",
		}

		_ = json.NewEncoder(w).Encode(res)
	})

	http.HandleFunc("/restart", restart)

	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		return
	}
}

func restart(w http.ResponseWriter, r *http.Request) {
	containerName := r.URL.Query().Get("container")
	if containerName == "" {
		_ = json.NewEncoder(w).Encode(RestartResponse{
			Ok:      false,
			Message: "Please provide a `container` as a URL query parameter",
		})
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	var res RestartResponse

	err = cli.ContainerRestart(context.Background(), containerName, container.StopOptions{})
	if err != nil {
		res = RestartResponse{
			Ok:      false,
			Message: err.Error(),
		}
	} else {
		res = RestartResponse{
			Ok:      true,
			Message: fmt.Sprintf("Restarted %s successfully", containerName),
		}
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println(err)
	}
}
