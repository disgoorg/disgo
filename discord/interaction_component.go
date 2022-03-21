package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
)

var (
	_ Interaction = (*ComponentInteraction)(nil)
)

type ComponentInteraction struct {
	BaseInteraction
	Data    ComponentInteractionData `json:"data"`
	Message Message                  `json:"message"`
}

func (i *ComponentInteraction) UnmarshalJSON(data []byte) error {
	type componentInteraction ComponentInteraction
	var vInteraction struct {
		Data json.RawMessage `json:"data"`
		componentInteraction
	}

	if err := json.Unmarshal(data, &vInteraction); err != nil {
		return err
	}

	var cType struct {
		Type ComponentType `json:"component_type"`
	}

	if err := json.Unmarshal(vInteraction.Data, &cType); err != nil {
		return err
	}

	var (
		interactionData ComponentInteractionData
		err             error
	)
	switch cType.Type {
	case ComponentTypeButton:
		v := ButtonInteractionData{}
		err = json.Unmarshal(vInteraction.Data, &v)
		interactionData = v

	case ComponentTypeSelectMenu:
		v := SelectMenuInteractionData{}
		err = json.Unmarshal(vInteraction.Data, &v)
		interactionData = v

	default:
		return fmt.Errorf("unkown component vInteraction data with type %d received", cType.Type)
	}
	if err != nil {
		return err
	}

	*i = ComponentInteraction(vInteraction.componentInteraction)

	i.Data = interactionData
	return nil
}

func (ComponentInteraction) Type() InteractionType {
	return InteractionTypeComponent
}

func (i ComponentInteraction) ButtonInteractionData() ButtonInteractionData {
	return i.Data.(ButtonInteractionData)
}

func (i ComponentInteraction) SelectMenuInteractionData() SelectMenuInteractionData {
	return i.Data.(SelectMenuInteractionData)
}

func (ComponentInteraction) interaction() {}

type ComponentInteractionData interface {
	Type() ComponentType
	CustomID() CustomID

	componentInteractionData()
}

type rawButtonInteractionData struct {
	Custom CustomID `json:"custom_id"`
}

type ButtonInteractionData struct {
	customID CustomID
}

func (d *ButtonInteractionData) UnmarshalJSON(data []byte) error {
	var v rawButtonInteractionData
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.customID = v.Custom
	return nil
}

func (d *ButtonInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawButtonInteractionData{
		Custom: d.customID,
	})
}

func (ButtonInteractionData) Type() ComponentType {
	return ComponentTypeButton
}

func (d ButtonInteractionData) CustomID() CustomID {
	return d.customID
}

func (ButtonInteractionData) componentInteractionData() {}

var (
	_ ComponentInteractionData = (*SelectMenuInteractionData)(nil)
)

type rawSelectMenuInteractionData struct {
	Custom CustomID `json:"custom_id"`
	Values []string `json:"values"`
}

type SelectMenuInteractionData struct {
	customID CustomID
	Values   []string
}

func (d *SelectMenuInteractionData) UnmarshalJSON(data []byte) error {
	var v rawSelectMenuInteractionData
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	d.customID = v.Custom
	d.Values = v.Values
	return nil
}

func (d *SelectMenuInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawSelectMenuInteractionData{
		Custom: d.customID,
		Values: d.Values,
	})
}

func (SelectMenuInteractionData) Type() ComponentType {
	return ComponentTypeSelectMenu
}

func (d SelectMenuInteractionData) CustomID() CustomID {
	return d.customID
}

func (SelectMenuInteractionData) componentInteractionData() {}
