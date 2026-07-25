package podman

import (
	"context"
	"log"
	"os/exec"
	"strings"

	"go.podman.io/podman/v6/pkg/bindings"
)

// getIdentity retrieves the SSH identity path for the Podman machine.
// It executes 'podman machine inspect' to get the identity path.
func getIdentity() string {
	identityOut, err := exec.Command("podman", "machine", "inspect", "--format", "{{.SSHConfig.IdentityPath}}").Output()
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(string(identityOut))
}

// getHostName finds the default Podman connection URI.
// It parses the list of system connections to find the one marked as default.
func getHostName() string {
	hostName, err := exec.Command("podman", "system", "connection", "list", "--format", "{{.URI}}").Output()
	if err != nil {
		log.Fatal(err)
	}

	isDefault, err := exec.Command("podman", "system", "connection", "list", "--format", "{{.Default}}").Output()
	if err != nil {
		log.Fatal(err)
	}

	isDefaultSlice := strings.Split(string(isDefault), "\n")
	for i := range isDefaultSlice {
		if isDefaultSlice[i] == "true" {
			return strings.TrimSpace(strings.Split(string(hostName), "\n")[i])
		}
	}
	return strings.TrimSpace(strings.Split(string(hostName), "\n")[0])
}

// ConnectPodman establishes a connection to the Podman service.
// It returns a context that carries the connection information for use with other Podman functions.
func ConnectPodman() context.Context {
	identity := getIdentity()
	uri := getHostName()

	conn, err := bindings.NewConnectionWithIdentity(
		context.Background(),
		uri,
		identity,
		false,
	)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
