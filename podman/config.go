package podman

import (
	"context"
	"os/exec"
	"strings"

	"go.podman.io/podman/v6/pkg/bindings"
)

type Connection struct {
	ctx context.Context
}

// getIdentity retrieves the SSH identity path for the Podman machine.
// It executes 'podman machine inspect' to get the identity path.
func getIdentity() (string, error) {
	identityOut, err := exec.Command("podman", "machine", "inspect", "--format", "{{.SSHConfig.IdentityPath}}").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(identityOut)), nil
}

// getHostName finds the default Podman connection URI.
// It parses the list of system connections to find the one marked as default.
func getHostName() (string, error) {
	hostName, err := exec.Command("podman", "system", "connection", "list", "--format", "{{.URI}}").Output()
	if err != nil {
		return "", err
	}

	isDefault, err := exec.Command("podman", "system", "connection", "list", "--format", "{{.Default}}").Output()
	if err != nil {
		return "", err
	}

	isDefaultSlice := strings.Split(string(isDefault), "\n")
	for i := range isDefaultSlice {
		if isDefaultSlice[i] == "true" {
			return strings.TrimSpace(strings.Split(string(hostName), "\n")[i]), nil
		}
	}
	return strings.TrimSpace(strings.Split(string(hostName), "\n")[0]), nil
}

// ConnectPodman establishes a connection to the Podman service.
// It returns a context that carries the connection information for use with other Podman functions.
func ConnectPodman() (Connection, error) {
	identity, err := getIdentity()
	if err != nil {
		return Connection{}, err
	}
	uri, err := getHostName()
	if err != nil {
		return Connection{}, err
	}
	conn, err := bindings.NewConnectionWithIdentity(
		context.Background(),
		uri,
		identity,
		false,
	)
	if err != nil {
		return Connection{}, err
	}
	return Connection{conn}, nil
}
