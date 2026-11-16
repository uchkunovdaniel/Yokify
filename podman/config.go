package podman

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"go.podman.io/podman/v6/pkg/bindings"
)

// Connection represents a connection to the Podman service.
type Connection struct {
	ctx context.Context
}

// SystemConnection is one entry from `podman system connection list`.
type SystemConnection struct {
	Name      string `json:"Name"`
	URI       string `json:"URI"`
	Identity  string `json:"Identity"`
	IsMachine bool   `json:"IsMachine"`
	Default   bool   `json:"Default"`
	ReadWrite bool   `json:"ReadWrite"`
}

// GetConnections returns the Podman system connections configured on this host.
func GetConnections() ([]SystemConnection, error) {
	out, err := exec.Command("podman", "system", "connection", "list", "--format", "json").Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("podman system connection list: %w: %s",
				err, strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, fmt.Errorf("podman system connection list: %w", err)
	}

	var conns []SystemConnection
	if err := json.Unmarshal(out, &conns); err != nil {
		return nil, fmt.Errorf("parse podman connection list: %w", err)
	}
	if len(conns) == 0 {
		return nil, errors.New("no podman connections configured")
	}
	return conns, nil
}

// defaultConnection returns the connection marked as default, or the first one.
func defaultConnection(conns []SystemConnection) SystemConnection {
	for _, c := range conns {
		if c.Default {
			return c
		}
	}
	return conns[0]
}

// ConnectPodman establishes a connection to the Podman service.
// It returns a context that carries the connection information for use with other Podman functions.
func ConnectPodman() (Connection, error) {
	conns, err := GetConnections()
	if err != nil {
		return Connection{}, err
	}
	c := defaultConnection(conns)

	ctx, err := bindings.NewConnectionWithIdentity(context.Background(), c.URI, c.Identity, c.IsMachine)
	if err != nil {
		return Connection{}, fmt.Errorf("connect to podman %q: %w", c.Name, err)
	}
	return Connection{ctx}, nil
}
