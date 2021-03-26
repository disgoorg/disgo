package api

import (
	"encoding/json"
	"testing"

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
