package discord

import (
	"testing"

	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/snowflake/v2"
)

const unknownGuildChannelJSON = `{
	"id": "1000000000000000001",
	"type": 99,
	"guild_id": "1000000000000000000",
	"position": 3,
	"name": "new-kind",
	"parent_id": "1000000000000000002",
	"last_message_id": "1000000000000000003",
	"flags": 2,
	"permission_overwrites": [
		{"id": "1000000000000000000", "type": 0, "allow": "1024", "deny": "0"},
		{"id": "1000000000000000004", "type": 9, "allow": "0", "deny": "2048"}
	],
	"new_field": {"nested": true}
}`

func TestUnmarshalChannel_UnknownType(t *testing.T) {
	var v UnmarshalChannel
	if err := json.Unmarshal([]byte(unknownGuildChannelJSON), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c, ok := v.Channel.(UnknownChannel)
	if !ok {
		t.Fatalf("expected UnknownChannel, got %T", v.Channel)
	}
	if c.Type() != 99 || c.ID() != 1000000000000000001 || c.GuildID() != 1000000000000000000 {
		t.Errorf("unexpected identity: type=%d id=%s guild=%s", c.Type(), c.ID(), c.GuildID())
	}
	if c.Name() != "new-kind" || c.Position() != 3 || c.Flags() != ChannelFlagPinned {
		t.Errorf("unexpected fields: name=%q position=%d flags=%d", c.Name(), c.Position(), c.Flags())
	}
	if c.ParentID() == nil || *c.ParentID() != 1000000000000000002 {
		t.Errorf("unexpected parent id: %v", c.ParentID())
	}
	if c.LastMessageID() == nil || *c.LastMessageID() != 1000000000000000003 {
		t.Errorf("unexpected last message id: %v", c.LastMessageID())
	}

	overwrites := c.PermissionOverwrites()
	if len(overwrites) != 2 {
		t.Fatalf("expected 2 overwrites, got %d", len(overwrites))
	}
	if _, ok := overwrites.Role(1000000000000000000); !ok {
		t.Error("expected role overwrite")
	}
	o, ok := overwrites[1].(UnknownPermissionOverwrite)
	if !ok {
		t.Fatalf("expected UnknownPermissionOverwrite, got %T", overwrites[1])
	}
	if o.Type() != 9 || o.ID() != 1000000000000000004 || o.Deny != PermissionSendMessages {
		t.Errorf("unexpected overwrite: %+v", o)
	}
}

func TestUnknownChannel_MarshalJSON(t *testing.T) {
	var v UnmarshalChannel
	if err := json.Unmarshal([]byte(`{"id":"1","type":99,"new_field":"kept"}`), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c := ApplyGuildIDToChannel(v.Channel.(GuildChannel), 2)

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]any
	if err = json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["new_field"] != "kept" || out["guild_id"] != "2" || out["id"] != "1" || out["type"] != float64(99) {
		t.Errorf("unexpected json: %s", data)
	}
}

func TestGatewayGuild_UnknownChannelType(t *testing.T) {
	data := `{"id": "1000000000000000000", "channels": [
		{"id": "1000000000000000005", "type": 0, "name": "general"},
		` + unknownGuildChannelJSON + `
	]}`

	var g GatewayGuild
	if err := json.Unmarshal([]byte(data), &g); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Channels) != 2 {
		t.Fatalf("expected 2 channels, got %d", len(g.Channels))
	}
	if _, ok := g.Channels[0].(GuildTextChannel); !ok {
		t.Errorf("expected GuildTextChannel, got %T", g.Channels[0])
	}
	if _, ok := g.Channels[1].(UnknownChannel); !ok {
		t.Errorf("expected UnknownChannel, got %T", g.Channels[1])
	}
}

func TestInteractionChannel_UnknownType(t *testing.T) {
	var c InteractionChannel
	if err := json.Unmarshal([]byte(`{"id":"1","type":99,"permissions":"2048"}`), &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := c.MessageChannel.(UnknownChannel); !ok {
		t.Errorf("expected UnknownChannel, got %T", c.MessageChannel)
	}
	if c.Permissions != PermissionSendMessages {
		t.Errorf("unexpected permissions: %d", c.Permissions)
	}
}

func TestUnmarshalWebhook_UnknownType(t *testing.T) {
	var v UnmarshalWebhook
	if err := json.Unmarshal([]byte(`{"id":"1","type":99,"name":"hook","channel_id":"2","guild_id":"3","new_field":1}`), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w, ok := v.Webhook.(UnknownWebhook)
	if !ok {
		t.Fatalf("expected UnknownWebhook, got %T", v.Webhook)
	}
	if w.Type() != 99 || w.ID() != 1 || w.Name() != "hook" || w.ChannelID != 2 || w.GuildID != 3 {
		t.Errorf("unexpected webhook: %+v", w)
	}

	data, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]any
	if err = json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["new_field"] != float64(1) || out["type"] != float64(99) {
		t.Errorf("unexpected json: %s", data)
	}
}

func TestUnmarshalIntegration_UnknownType(t *testing.T) {
	var v UnmarshalIntegration
	if err := json.Unmarshal([]byte(`{"id":"1","type":"new_kind","name":"thing","enabled":true}`), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	i, ok := v.Integration.(UnknownIntegration)
	if !ok {
		t.Fatalf("expected UnknownIntegration, got %T", v.Integration)
	}
	if i.Type() != "new_kind" || i.ID() != 1 || i.Name != "thing" || !i.Enabled {
		t.Errorf("unexpected integration: %+v", i)
	}
}

func TestUnmarshalApplicationCommand_UnknownType(t *testing.T) {
	var v UnmarshalApplicationCommand
	if err := json.Unmarshal([]byte(`{"id":"1","type":99,"application_id":"2","name":"launch","version":"3","default_member_permissions":"8"}`), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c, ok := v.ApplicationCommand.(UnknownApplicationCommand)
	if !ok {
		t.Fatalf("expected UnknownApplicationCommand, got %T", v.ApplicationCommand)
	}
	if c.Type() != 99 || c.ID() != 1 || c.ApplicationID() != 2 || c.Name() != "launch" || c.Version() != 3 || c.DefaultMemberPermissions() != PermissionAdministrator {
		t.Errorf("unexpected command: %+v", c)
	}
}

func TestApplicationCommandPermissions_UnknownType(t *testing.T) {
	var p ApplicationCommandPermissions
	if err := json.Unmarshal([]byte(`{"id":"1","application_id":"2","guild_id":"3","permissions":[
		{"id":"4","type":1,"permission":true},
		{"id":"5","type":99,"permission":false}
	]}`), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Permissions) != 2 {
		t.Fatalf("expected 2 permissions, got %d", len(p.Permissions))
	}
	u, ok := p.Permissions[1].(UnknownApplicationCommandPermission)
	if !ok {
		t.Fatalf("expected UnknownApplicationCommandPermission, got %T", p.Permissions[1])
	}
	if u.Type() != 99 || u.ID() != snowflake.ID(5) || u.Permission {
		t.Errorf("unexpected permission: %+v", u)
	}
}
