package subresources_test

import (
	"testing"

	"github.com/nalajala4naresh/chvmm-api/subresources"
	"github.com/nalajala4naresh/chvmm-api/v1beta1"
)

// The subresource group is derived from the VirtualMachine group; they must move together.
func TestGroupMatchesVirtualMachineGroup(t *testing.T) {
	if want := "subresources." + v1beta1.GroupVersion.Group; subresources.GroupName != want {
		t.Fatalf("GroupName = %q, want %q", subresources.GroupName, want)
	}
	if subresources.Version != v1beta1.GroupVersion.Version {
		t.Fatalf("Version = %q, want %q", subresources.Version, v1beta1.GroupVersion.Version)
	}
}

func TestPath(t *testing.T) {
	got := subresources.Path("ns1", "vm1", subresources.PortForward, "22")
	want := "/apis/subresources.cloudhypervisor.quill.today/v1beta1/namespaces/ns1/virtualmachines/vm1/portforward/22"
	if got != want {
		t.Fatalf("Path = %q, want %q", got, want)
	}
}
