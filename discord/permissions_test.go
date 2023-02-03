package discord

import (
	"testing"

	"github.com/disgoorg/json"

	"github.com/stretchr/testify/assert"
)

type permissionTestStruct struct {
	Permissions Permissions `json:"permissions"`
}

func TestPermissions_MarshalJSON(t *testing.T) {
	someStruct := permissionTestStruct{
		Permissions: PermissionAddReactions | PermissionChangeNickname,
	}

	jsonPerms, err := json.Marshal(someStruct)
	assert.NoError(t, err)
	assert.Equal(t, "{\"permissions\":\"67108928\"}", string(jsonPerms))
}

func TestPermissions_UnmarshalJSON(t *testing.T) {
	var someStruct permissionTestStruct

	err := json.Unmarshal([]byte("{\"permissions\":\"67108928\"}"), &someStruct)
	assert.NoError(t, err)
	assert.Equal(t, permissionTestStruct{Permissions: PermissionAddReactions | PermissionChangeNickname}, someStruct)
}

func TestPermissions_Add(t *testing.T) {
	assert.Equal(t, PermissionAddReactions, Permissions.Add(PermissionAddReactions))
}

func TestPermissions_Remove(t *testing.T) {
	assert.Equal(t, PermissionManageChannels, (PermissionAddReactions | PermissionManageChannels).Remove(PermissionAddReactions))
}

func TestPermissions_Has(t *testing.T) {
	assert.True(t, PermissionAddReactions.Has(PermissionAddReactions))
}

func TestPermissions_Missing(t *testing.T) {
	assert.True(t, PermissionManageChannels.Missing(PermissionAddReactions))
}
