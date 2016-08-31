package di

import (
	"github.com/docker/engine-api/client"
	"github.com/icecave/honeycomb/src/docker"
)

// DockerClient returns the docker client used to access the swarm.
func (con *Container) DockerClient() client.APIClient {
	return con.get(
		"docker.client",
		func() (interface{}, error) {
			return client.NewEnvClient()
		},
		nil,
	).(client.APIClient)
}

// ServiceLoader returns the service loader used to load Docker services.
func (con *Container) ServiceLoader() *docker.ServiceLoader {
	return con.get(
		"docker.service-loader",
		func() (interface{}, error) {
			return &docker.ServiceLoader{
				Client:    con.DockerClient(),
				Inspector: con.ServiceInspector(),
				Logger:    con.Logger(),
			}, nil
		},
		nil,
	).(*docker.ServiceLoader)
}

// ServiceInspector returns the service inspector used to create endpoints from
// Docker services.
func (con *Container) ServiceInspector() *docker.ServiceInspector {
	return con.get(
		"docker.service-inspector",
		func() (interface{}, error) {
			return &docker.ServiceInspector{
				Client: con.DockerClient(),
			}, nil
		},
		nil,
	).(*docker.ServiceInspector)
}
