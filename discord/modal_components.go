package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
)

type ModalComponent interface {
	Type() ComponentType
	modalComponent()
}

type ModalContainerComponent interface {
	ModalComponent
	Components() []ModalInteractiveComponent
	modalContainerComponent()
}

type ModalInteractiveComponent interface {
	ModalComponent
	ID() CustomID
	Value() string
	modalInteractiveComponent()
}

type UnmarshalModalComponent struct {
	ModalComponent
}

func (u *UnmarshalModalComponent) UnmarshalJSON(data []byte) error {
	var cType struct {
		Type ComponentType `json:"type"`
	}

	if err := json.Unmarshal(data, &cType); err != nil {
		return err
	}

	var (
		modalSubmitComponent ModalComponent
		err                  error
	)

	switch cType.Type {
	case ComponentTypeActionRow:
		v := ModalActionRowComponent{}
		err = json.Unmarshal(data, &v)
		modalSubmitComponent = v

	case ComponentTypeInputText:
		v := ModalInputTextComponent{}
		err = json.Unmarshal(data, &v)
		modalSubmitComponent = v

	default:
		err = fmt.Errorf("unkown component with type %d received", cType.Type)
	}
	if err != nil {
		return err
	}

	u.ModalComponent = modalSubmitComponent
	return nil
}

var (
	_ ModalComponent          = (*ModalActionRowComponent)(nil)
	_ ModalContainerComponent = (*ModalActionRowComponent)(nil)
)

type ModalActionRowComponent []ModalInteractiveComponent

func (c *ModalActionRowComponent) UnmarshalJSON(data []byte) error {
	var actionRow struct {
		Components []UnmarshalModalComponent `json:"components"`
	}

	if err := json.Unmarshal(data, &actionRow); err != nil {
		return err
	}

	if len(actionRow.Components) > 0 {
		*c = make([]ModalInteractiveComponent, len(actionRow.Components))
		for i := range actionRow.Components {
			(*c)[i] = actionRow.Components[i].ModalComponent.(ModalInteractiveComponent)
		}
	}

	return nil
}

func (ModalActionRowComponent) Type() ComponentType {
	return ComponentTypeInputText
}

func (c ModalActionRowComponent) Components() []ModalInteractiveComponent {
	return c
}

func (ModalActionRowComponent) modalComponent()          {}
func (ModalActionRowComponent) modalContainerComponent() {}

var (
	_ ModalComponent            = (*ModalInputTextComponent)(nil)
	_ ModalInteractiveComponent = (*ModalInputTextComponent)(nil)
)

type ModalInputTextComponent struct {
	CustomID  CustomID `json:"custom_id"`
	TextValue string   `json:"value"`
}

func (ModalInputTextComponent) Type() ComponentType {
	return ComponentTypeInputText
}

func (c ModalInputTextComponent) ID() CustomID {
	return c.CustomID
}
func (c ModalInputTextComponent) Value() string {
	return c.TextValue
}

func (ModalInputTextComponent) modalComponent()            {}
func (ModalInputTextComponent) modalInteractiveComponent() {}
