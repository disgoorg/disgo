package discord

import (
	"testing"

	"github.com/disgoorg/json/v2"
)

type permissionTestStruct struct {
	Permissions Permissions `json:"permissions"`
}

func TestPermissions_MarshalJSON(t *testing.T) {
	someStruct := permissionTestStruct{
		Permissions: PermissionAddReactions | PermissionChangeNickname,
	}

	jsonPerms, err := json.Marshal(someStruct)
	if err != nil {
		t.Fatalf("unexpected error marshaling: %v", err)
	}

	expected := `{"permissions":"67108928"}`
	if string(jsonPerms) != expected {
		t.Errorf("expected %s, got %s", expected, string(jsonPerms))
	}
}

func TestPermissions_UnmarshalJSON(t *testing.T) {
	var someStruct permissionTestStruct

	err := json.Unmarshal([]byte(`{"permissions":"67108928"}`), &someStruct)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling: %v", err)
	}

	expected := permissionTestStruct{Permissions: PermissionAddReactions | PermissionChangeNickname}
	if someStruct != expected {
		t.Errorf("expected %+v, got %+v", expected, someStruct)
	}
}

func TestPermissions_Add(t *testing.T) {
	got := Permissions.Add(PermissionAddReactions)
	if got != PermissionAddReactions {
		t.Errorf("expected %v, got %v", PermissionAddReactions, got)
	}
}

func TestPermissions_Remove(t *testing.T) {
	got := (PermissionAddReactions | PermissionManageChannels).Remove(PermissionAddReactions)
	if got != PermissionManageChannels {
		t.Errorf("expected %v, got %v", PermissionManageChannels, got)
	}
}

func TestPermissions_Has(t *testing.T) {
	if !PermissionAddReactions.Has(PermissionAddReactions) {
		t.Errorf("expected PermissionAddReactions to have PermissionAddReactions")
	}
}

func TestPermissions_Missing(t *testing.T) {
	if !PermissionManageChannels.Missing(PermissionAddReactions) {
		t.Errorf("expected PermissionManageChannels to be missing PermissionAddReactions")
	}
}
