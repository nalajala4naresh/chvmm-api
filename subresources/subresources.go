// Package subresources defines the wire contract of the VirtualMachine subresource API served by
// ch-vmm-gateway through the Kubernetes API aggregation layer. Clients use it through the vmclient
// package; the gateway uses it to serve the same paths.
//
// Every subresource lives under
//
//	/apis/subresources.cloudhypervisor.quill.today/v1beta1/namespaces/{namespace}/virtualmachines/{name}/{subresource}
//
// Streaming subresources (console, portforward) are opened with a WebSocket GET using the
// PlainStreamProtocol subprotocol. Action subresources are POSTs with an empty body. RBAC grants
// both kinds with the "create" verb on "virtualmachines/<subresource>" in GroupName.
package subresources

import (
	"fmt"
	"net/url"
)

const (
	// GroupName is the API group of the subresource API.
	GroupName = "subresources.cloudhypervisor.quill.today"
	// Version is the API version of the subresource API.
	Version = "v1beta1"

	// PlainStreamProtocol is the WebSocket subprotocol of streaming subresources. Each binary
	// message carries raw bytes of the stream; there is no channel multiplexing.
	PlainStreamProtocol = "plain.cloudhypervisor.quill.today"
)

// Subresource names.
const (
	// Console streams the VM's serial console.
	Console = "console"
	// PortForward streams a TCP connection to a port of the VM; the port is the last path segment.
	PortForward = "portforward"
	// Pause pauses a Running VM.
	Pause = "pause"
	// Unpause resumes a Paused VM.
	Unpause = "unpause"
	// Restart reboots a Running VM.
	Restart = "restart"
	// Migrate live-migrates a Running VM and responds with the created VirtualMachineMigration.
	Migrate = "migrate"
)

// Path returns the absolute API path of a VM subresource. extra is appended as further path
// segments, such as the port of PortForward.
func Path(namespace, name, subresource string, extra ...string) string {
	p := fmt.Sprintf("/apis/%s/%s/namespaces/%s/virtualmachines/%s/%s",
		GroupName, Version, url.PathEscape(namespace), url.PathEscape(name), subresource)
	for _, e := range extra {
		p += "/" + url.PathEscape(e)
	}
	return p
}
