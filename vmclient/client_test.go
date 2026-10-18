package vmclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"

	"github.com/nalajala4naresh/chvmm-api/subresources"
	"github.com/nalajala4naresh/chvmm-api/v1beta1"
)

var vmResource = schema.GroupResource{Group: "cloudhypervisor.quill.today", Resource: "virtualmachines"}

// fakeGateway serves the subresource API like ch-vmm-gateway: actions record their path,
// streams echo, and VMs named "forbidden" or "stopped" fail like RBAC and phase checks do.
type fakeGateway struct {
	mu       sync.Mutex
	requests []string
}

func (g *fakeGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mu.Lock()
	g.requests = append(g.requests, r.Method+" "+r.URL.Path)
	g.mu.Unlock()

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	// apis/<group>/<version>/namespaces/<ns>/virtualmachines/<name>/<subresource>[/<port>]
	if len(parts) < 8 || parts[1] != subresources.GroupName || parts[5] != "virtualmachines" {
		writeStatus(w, apierrors.NewNotFound(schema.GroupResource{}, r.URL.Path))
		return
	}
	name, subresource := parts[6], parts[7]
	switch name {
	case "forbidden":
		writeStatus(w, apierrors.NewForbidden(vmResource, name, errors.New("denied")))
		return
	case "stopped":
		writeStatus(w, apierrors.NewConflict(vmResource, name, errors.New("VM is not running")))
		return
	}

	switch {
	case r.Method == http.MethodGet && (subresource == subresources.Console || subresource == subresources.PortForward):
		upgrader := websocket.Upgrader{Subprotocols: []string{subresources.PlainStreamProtocol}}
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer ws.Close()
		for {
			msgType, data, err := ws.ReadMessage()
			if err != nil {
				return
			}
			if string(data) == "bye" {
				_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}
			if err := ws.WriteMessage(msgType, data); err != nil {
				return
			}
		}
	case r.Method == http.MethodPost && subresource == subresources.Migrate:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(v1beta1.VirtualMachineMigration{
			ObjectMeta: metav1.ObjectMeta{Name: name + "-migration-abcde", Namespace: parts[4]},
			Spec:       v1beta1.VirtualMachineMigrationSpec{VMName: name},
		})
	case r.Method == http.MethodPost:
		w.WriteHeader(http.StatusAccepted)
	default:
		writeStatus(w, apierrors.NewMethodNotSupported(vmResource, r.Method))
	}
}

func (g *fakeGateway) seen() []string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]string(nil), g.requests...)
}

func writeStatus(w http.ResponseWriter, err error) {
	var apiStatus apierrors.APIStatus
	errors.As(err, &apiStatus)
	status := apiStatus.Status()
	status.TypeMeta = metav1.TypeMeta{Kind: "Status", APIVersion: "v1"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status.Code))
	_ = json.NewEncoder(w).Encode(status)
}

func newTestClient(t *testing.T, prefix string) (*Client, *fakeGateway) {
	t.Helper()
	g := &fakeGateway{}
	var handler http.Handler = g
	if prefix != "" {
		handler = http.StripPrefix(prefix, g)
	}
	srv := httptest.NewTLSServer(handler)
	t.Cleanup(srv.Close)
	c, err := NewForConfig(&rest.Config{
		Host:            srv.URL + prefix,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	return c, g
}

func TestActions(t *testing.T) {
	c, g := newTestClient(t, "")
	vms := c.VirtualMachines("ns1")
	ctx := context.Background()

	for _, call := range []func(context.Context, string) error{vms.Pause, vms.Unpause, vms.Restart} {
		if err := call(ctx, "vm1"); err != nil {
			t.Fatal(err)
		}
	}
	base := "POST /apis/" + subresources.GroupName + "/v1beta1/namespaces/ns1/virtualmachines/vm1/"
	want := []string{base + "pause", base + "unpause", base + "restart"}
	if got := g.seen(); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestActionErrors(t *testing.T) {
	c, _ := newTestClient(t, "")
	vms := c.VirtualMachines("ns1")
	if err := vms.Pause(context.Background(), "forbidden"); !apierrors.IsForbidden(err) {
		t.Fatalf("got %v, want Forbidden", err)
	}
	if err := vms.Restart(context.Background(), "stopped"); !apierrors.IsConflict(err) {
		t.Fatalf("got %v, want Conflict", err)
	}
}

func TestMigrate(t *testing.T) {
	c, _ := newTestClient(t, "")
	m, err := c.VirtualMachines("ns1").Migrate(context.Background(), "vm1")
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "vm1-migration-abcde" || m.Namespace != "ns1" || m.Spec.VMName != "vm1" {
		t.Fatalf("unexpected migration %+v", m)
	}
}

func TestConsole(t *testing.T) {
	c, g := newTestClient(t, "")
	conn, err := c.VirtualMachines("ns1").Console(context.Background(), "vm1")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 5)
	if _, err := io.ReadFull(conn, buf); err != nil {
		t.Fatal(err)
	}
	if string(buf) != "hello" {
		t.Fatalf("read %q", buf)
	}
	if got := g.seen()[0]; got != "GET /apis/"+subresources.GroupName+"/v1beta1/namespaces/ns1/virtualmachines/vm1/console" {
		t.Fatalf("request %q", got)
	}

	// A normal close by the server reads as EOF.
	if _, err := conn.Write([]byte("bye")); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Read(buf); !errors.Is(err, io.EOF) {
		t.Fatalf("got %v, want EOF", err)
	}
}

func TestPortForwardLargeWrites(t *testing.T) {
	c, g := newTestClient(t, "")
	conn, err := c.VirtualMachines("ns1").PortForward(context.Background(), "vm1", 22)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Larger than one message, so it spans several messages each way.
	payload := bytes.Repeat([]byte("0123456789abcdef"), 3*maxMessageSize/16+7)
	errCh := make(chan error, 1)
	go func() {
		_, err := conn.Write(payload)
		errCh <- err
	}()
	got := make([]byte, len(payload))
	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadFull(conn, got); err != nil {
		t.Fatal(err)
	}
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatal("payload mismatch")
	}
	if got := g.seen()[0]; !strings.HasSuffix(got, "/virtualmachines/vm1/portforward/22") {
		t.Fatalf("request %q", got)
	}
}

func TestStreamErrorsAreAPIErrors(t *testing.T) {
	c, _ := newTestClient(t, "")
	vms := c.VirtualMachines("ns1")
	if _, err := vms.Console(context.Background(), "forbidden"); !apierrors.IsForbidden(err) {
		t.Fatalf("got %v, want Forbidden", err)
	}
	if _, err := vms.PortForward(context.Background(), "stopped", 22); !apierrors.IsConflict(err) {
		t.Fatalf("got %v, want Conflict", err)
	}
	if _, err := vms.PortForward(context.Background(), "vm1", 0); err == nil {
		t.Fatal("port 0 accepted")
	}
}

// A server behind a path prefix, as with Rancher or other cluster proxies.
func TestServerPathPrefix(t *testing.T) {
	c, _ := newTestClient(t, "/k8s/clusters/c1")
	vms := c.VirtualMachines("ns1")
	if err := vms.Pause(context.Background(), "vm1"); err != nil {
		t.Fatal(err)
	}
	conn, err := vms.Console(context.Background(), "vm1")
	if err != nil {
		t.Fatal(err)
	}
	_ = conn.Close()
}
