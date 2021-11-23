package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
)

type ModalSubmitComponent interface {
	Type() ComponentType
	modalSubmitComponent()
}

type ModalSubmitContainerComponent interface {
	ModalSubmitComponent
	Components() []ModalSubmitInteractiveComponent
	modalSubmitContainerComponent()
}

type ModalSubmitInteractiveComponent interface {
	ModalSubmitComponent
	ID() CustomID
	Value() string
	modalSubmitInteractiveComponent()
}

type UnmarshalModalSubmitComponent struct {
	ModalSubmitComponent
}

func (u *UnmarshalModalSubmitComponent) UnmarshalJSON(data []byte) error {
	var cType struct {
		Type ComponentType `json:"type"`
	}

	if err := json.Unmarshal(data, &cType); err != nil {
		return err
	}

	var (
		modalSubmitComponent ModalSubmitComponent
		err                  error
	)

	switch cType.Type {
	case ComponentTypeActionRow:
		v := ModalSubmitActionRowComponent{}
		err = json.Unmarshal(data, &v)
		modalSubmitComponent = v

	case ComponentTypeInputText:
		v := ModalSubmitInputTextComponent{}
		err = json.Unmarshal(data, &v)
		modalSubmitComponent = v

	default:
		err = fmt.Errorf("unkown component with type %d received", cType.Type)
	}
	if err != nil {
		return err
	}

	u.ModalSubmitComponent = modalSubmitComponent
	return nil
}

var (
	_ ModalSubmitComponent          = (*ModalSubmitActionRowComponent)(nil)
	_ ModalSubmitContainerComponent = (*ModalSubmitActionRowComponent)(nil)
)

type ModalSubmitActionRowComponent []ModalSubmitInteractiveComponent

func (c *ModalSubmitActionRowComponent) UnmarshalJSON(data []byte) error {
	var actionRow struct {
		Components []UnmarshalModalSubmitComponent `json:"components"`
	}

	if err := json.Unmarshal(data, &actionRow); err != nil {
		return err
	}

	if len(actionRow.Components) > 0 {
		*c = make([]ModalSubmitInteractiveComponent, len(actionRow.Components))
		for i := range actionRow.Components {
			(*c)[i] = actionRow.Components[i].ModalSubmitComponent.(ModalSubmitInteractiveComponent)
		}
	}

	return nil
}

func (ModalSubmitActionRowComponent) Type() ComponentType {
	return ComponentTypeInputText
}

func (c ModalSubmitActionRowComponent) Components() []ModalSubmitInteractiveComponent {
	return c
}

func (ModalSubmitActionRowComponent) modalSubmitComponent()          {}
func (ModalSubmitActionRowComponent) modalSubmitContainerComponent() {}

var (
	_ ModalSubmitComponent            = (*ModalSubmitInputTextComponent)(nil)
	_ ModalSubmitInteractiveComponent = (*ModalSubmitInputTextComponent)(nil)
)

type ModalSubmitInputTextComponent struct {
	CustomID  CustomID `json:"custom_id"`
	TextValue string   `json:"value"`
}

func (ModalSubmitInputTextComponent) Type() ComponentType {
	return ComponentTypeInputText
}

func (c ModalSubmitInputTextComponent) ID() CustomID {
	return c.CustomID
}
func (c ModalSubmitInputTextComponent) Value() string {
	return c.TextValue
}

func (ModalSubmitInputTextComponent) modalSubmitComponent()            {}
func (ModalSubmitInputTextComponent) modalSubmitInteractiveComponent() {}
