package podman

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

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
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
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
func ConnectPodman(ctx context.Context) (*Connection, error) {
	if host := os.Getenv("CONTAINER_HOST"); host != "" {
		return connect(ctx, SystemConnection{
			Name: "CONTAINER_HOST", URI: host,
			Identity: os.Getenv("CONTAINER_SSHKEY"),
		})
	}
	conns, err := GetConnections()
	if err == nil {
		return connect(ctx, defaultConnection(conns))
	}
	if runtime.GOOS != "linux" {
		return nil, err
	}

	errs := []error{err}
	for _, path := range localSocketPaths() {
		if _, statErr := os.Stat(path); statErr != nil {
			continue // socket doesn't exist, try the next one
		}
		conn, connErr := connect(ctx, SystemConnection{Name: "local", URI: "unix://" + path})
		if connErr == nil {
			return conn, nil
		}
		errs = append(errs, connErr)
	}
	errs = append(errs, errors.New("no Podman socket found; enable it with: systemctl --user enable --now podman.socket"))
	return nil, errors.Join(errs...)
}

func localSocketPaths() []string {
	var paths []string
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		paths = append(paths, filepath.Join(dir, "podman", "podman.sock"))
	}
	paths = append(
		paths,
		fmt.Sprintf("/run/user/%d/podman/podman.sock", os.Getuid()), // rootless
		"/run/podman/podman.sock",                                   // rootful
	)
	return paths
}

func connect(ctx context.Context, c SystemConnection) (*Connection, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := bindings.NewConnectionWithIdentity(ctx, c.URI, c.Identity, c.IsMachine)
	if err != nil {
		return nil, fmt.Errorf("connect to podman %q (%s): %w", c.Name, c.URI, err)
	}
	return &Connection{ctx: context.WithoutCancel(conn)}, nil
}

// call returns a context that carries the Podman connection and is
// canceled when ctx is (request closed, timeout). Always call the stop func.
func (connection *Connection) call(ctx context.Context) (context.Context, func()) {
	merged, cancel := context.WithCancelCause(connection.ctx)
	stop := context.AfterFunc(ctx,
		func() {
			cancel(context.Cause(ctx))
		})

	return merged, func() {
		stop()
		cancel(nil)
	}
}
