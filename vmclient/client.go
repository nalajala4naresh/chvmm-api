// Package vmclient is a Go client for the VirtualMachine subresource API served by
// ch-vmm-gateway (see package subresources): serial console, port-forward, and VM actions.
//
// Requests go through kube-apiserver like any other API call, so the caller's credentials and
// RBAC apply. A minimal example:
//
//	c, err := vmclient.NewForConfig(restConfig)
//	...
//	vms := c.VirtualMachines("default")
//	if err := vms.Pause(ctx, "my-vm"); apierrors.IsForbidden(err) { ... }
//
//	conn, err := vms.PortForward(ctx, "my-vm", 22) // a net.Conn to the VM's port 22
//
// Errors returned by the gateway or kube-apiserver are *apierrors.StatusError values, also for
// streaming subresources, so the k8s.io/apimachinery/pkg/api/errors helpers work on them.
package vmclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	wstransport "k8s.io/client-go/transport/websocket"
	"k8s.io/streaming/pkg/httpstream"

	"github.com/nalajala4naresh/chvmm-api/subresources"
	"github.com/nalajala4naresh/chvmm-api/v1beta1"
)

// Interface is implemented by Client. Depend on it to substitute a fake in tests.
type Interface interface {
	VirtualMachines(namespace string) VirtualMachineInterface
}

// VirtualMachineInterface calls the subresources of VirtualMachines in one namespace.
type VirtualMachineInterface interface {
	// Pause pauses a Running VM. It returns a Conflict error if the VM is not Running.
	Pause(ctx context.Context, name string) error
	// Unpause resumes a Paused VM. It returns a Conflict error if the VM is not Paused.
	Unpause(ctx context.Context, name string) error
	// Restart reboots a Running VM. It returns a Conflict error if the VM is not Running.
	Restart(ctx context.Context, name string) error
	// Migrate starts a live migration of a Running VM and returns the created migration. It
	// returns a Conflict error if the VM is not Running or is already migrating.
	Migrate(ctx context.Context, name string) (*v1beta1.VirtualMachineMigration, error)

	// Console connects to the VM's serial console. ctx bounds only connection setup; close
	// the returned connection to end the session.
	Console(ctx context.Context, name string) (net.Conn, error)
	// PortForward connects to a TCP port of the VM. ctx bounds only connection setup; close the
	// returned connection to end it. Each call opens one TCP connection to the VM.
	PortForward(ctx context.Context, name string, port int) (net.Conn, error)
}

// Client calls the subresource API through kube-apiserver.
type Client struct {
	config *rest.Config
	rest   rest.Interface
}

var _ Interface = &Client{}

// NewForConfig returns a Client that uses the credentials and server in config.
func NewForConfig(config *rest.Config) (*Client, error) {
	cfg := rest.CopyConfig(config)
	cfg.GroupVersion = &schema.GroupVersion{Group: subresources.GroupName, Version: subresources.Version}
	cfg.APIPath = "/apis"
	cfg.ContentType = runtime.ContentTypeJSON
	cfg.AcceptContentTypes = runtime.ContentTypeJSON
	cfg.NegotiatedSerializer = statusCodecs.WithoutConversion()
	if cfg.UserAgent == "" {
		cfg.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	rc, err := rest.RESTClientFor(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{config: cfg, rest: rc}, nil
}

// statusCodecs decodes the metav1.Status error bodies of the API.
var statusCodecs = func() serializer.CodecFactory {
	scheme := runtime.NewScheme()
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{Version: "v1"})
	return serializer.NewCodecFactory(scheme)
}()

// VirtualMachines returns a client for VMs in namespace.
func (c *Client) VirtualMachines(namespace string) VirtualMachineInterface {
	return &virtualMachines{client: c, namespace: namespace}
}

type virtualMachines struct {
	client    *Client
	namespace string
}

func (v *virtualMachines) Pause(ctx context.Context, name string) error {
	return v.action(ctx, name, subresources.Pause)
}

func (v *virtualMachines) Unpause(ctx context.Context, name string) error {
	return v.action(ctx, name, subresources.Unpause)
}

func (v *virtualMachines) Restart(ctx context.Context, name string) error {
	return v.action(ctx, name, subresources.Restart)
}

func (v *virtualMachines) Migrate(ctx context.Context, name string) (*v1beta1.VirtualMachineMigration, error) {
	body, err := v.post(ctx, name, subresources.Migrate).Raw()
	if err != nil {
		return nil, err
	}
	// Decode with encoding/json: the response may omit apiVersion and kind.
	var migration v1beta1.VirtualMachineMigration
	if err := json.Unmarshal(body, &migration); err != nil {
		return nil, fmt.Errorf("decode migration: %w", err)
	}
	return &migration, nil
}

func (v *virtualMachines) action(ctx context.Context, name, subresource string) error {
	return v.post(ctx, name, subresource).Error()
}

func (v *virtualMachines) post(ctx context.Context, name, subresource string) rest.Result {
	return v.client.rest.Post().
		Namespace(v.namespace).
		Resource("virtualmachines").
		Name(name).
		SubResource(subresource).
		Do(ctx)
}

func (v *virtualMachines) Console(ctx context.Context, name string) (net.Conn, error) {
	return v.client.dial(ctx, subresources.Path(v.namespace, name, subresources.Console))
}

func (v *virtualMachines) PortForward(ctx context.Context, name string, port int) (net.Conn, error) {
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid port %d", port)
	}
	return v.client.dial(ctx, subresources.Path(v.namespace, name, subresources.PortForward, strconv.Itoa(port)))
}

// dial opens a WebSocket to a streaming subresource at the absolute API path p.
func (c *Client) dial(ctx context.Context, p string) (net.Conn, error) {
	rt, holder, err := wstransport.RoundTripperFor(c.config)
	if err != nil {
		return nil, err
	}
	u, err := serverURL(c.config, p)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	ws, err := wstransport.Negotiate(rt, holder, req, subresources.PlainStreamProtocol)
	if err != nil {
		return nil, streamError(err)
	}
	return newConn(ws), nil
}

// serverURL joins the absolute API path p to the server URL, keeping any path prefix of the
// server (for example, a cluster behind a proxy at https://host/k8s/clusters/c1).
func serverURL(config *rest.Config, p string) (*url.URL, error) {
	u, _, err := rest.DefaultServerUrlFor(config)
	if err != nil {
		return nil, err
	}
	raw := strings.TrimSuffix(u.EscapedPath(), "/") + p
	decoded, err := url.PathUnescape(raw)
	if err != nil {
		return nil, err
	}
	u.Path, u.RawPath = decoded, raw
	return u, nil
}

// streamError returns the API error inside a failed WebSocket upgrade, so callers can use
// apierrors.IsForbidden and friends on streaming errors too.
func streamError(err error) error {
	var upgradeErr *httpstream.UpgradeFailureError
	if errors.As(err, &upgradeErr) {
		var statusErr *apierrors.StatusError
		if errors.As(upgradeErr.Cause, &statusErr) {
			return statusErr
		}
	}
	return err
}
