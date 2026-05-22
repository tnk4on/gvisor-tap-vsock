package transport

import (
	"errors"
	"net"
	"net/url"

	"github.com/linuxkit/virtsock/pkg/hvsock"
)

// Dial connects to a Hyper-V vsock endpoint.
// URL format: vsock://VM-GUID:SERVICE-GUID/path
func Dial(endpoint string) (net.Conn, string, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, "", err
	}
	switch parsed.Scheme {
	case "vsock":
		vmGUID, err := hvsock.GUIDFromString(parsed.Hostname())
		if err != nil {
			return nil, "", err
		}
		svcGUID, err := hvsock.GUIDFromString(parsed.Port())
		if err != nil {
			return nil, "", err
		}
		conn, err := hvsock.Dial(hvsock.Addr{
			VMID:      vmGUID,
			ServiceID: svcGUID,
		})
		if err != nil {
			return nil, "", err
		}
		return conn, parsed.Path, nil
	case "unix":
		conn, err := net.Dial("unix", parsed.Path)
		return conn, "/connect", err
	default:
		return nil, "", errors.New("unexpected scheme")
	}
}
