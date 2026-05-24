package transport

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/linuxkit/virtsock/pkg/hvsock"
)

// Dial connects to a Hyper-V vsock endpoint.
// URL format: vsock://VM-GUID/SERVICE-GUID[/path]
// Example: vsock://d3e7a096-d081-40e6-a8c6-336370555af0/00000401-FACB-11E6-BD58-64006A7986D3/connect
func Dial(endpoint string) (net.Conn, string, error) {
	if !strings.HasPrefix(endpoint, "vsock://") {
		return nil, "", fmt.Errorf("unsupported scheme in %q", endpoint)
	}
	rest := strings.TrimPrefix(endpoint, "vsock://")

	parts := strings.SplitN(rest, "/", 3)
	if len(parts) < 2 {
		return nil, "", errors.New("vsock dial URL must be vsock://VM-GUID/SERVICE-GUID[/path]")
	}

	vmGUID, err := hvsock.GUIDFromString(parts[0])
	if err != nil {
		return nil, "", fmt.Errorf("invalid VM GUID %q: %w", parts[0], err)
	}
	svcGUID, err := hvsock.GUIDFromString(parts[1])
	if err != nil {
		return nil, "", fmt.Errorf("invalid Service GUID %q: %w", parts[1], err)
	}

	var path string
	if len(parts) == 3 {
		path = "/" + parts[2]
	}

	conn, err := hvsock.Dial(hvsock.Addr{
		VMID:      vmGUID,
		ServiceID: svcGUID,
	})
	if err != nil {
		return nil, "", err
	}
	return conn, path, nil
}
