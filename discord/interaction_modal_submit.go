package discord

import "github.com/DisgoOrg/disgo/json"

var (
	_ Interaction = (*ModalSubmitInteraction)(nil)
)

type ModalSubmitInteraction struct {
	baseInteractionImpl
	Data ModalSubmitInteractionData `json:"data"`
}

func (ModalSubmitInteraction) Type() InteractionType {
	return InteractionTypeModalSubmit
}

type ModalSubmitInteractionData struct {
	CustomID   CustomID                    `json:"custom_id"`
	Components map[CustomID]InputComponent `json:"components"`
}

func (d *ModalSubmitInteractionData) UnmarshalJSON(data []byte) error {
	type modalSubmitInteractionData ModalSubmitInteractionData
	var iData struct {
		Components []UnmarshalComponent `json:"components"`
		modalSubmitInteractionData
	}

	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}

	*d = ModalSubmitInteractionData(iData.modalSubmitInteractionData)

	if len(iData.Components) > 0 {
		d.Components = make(map[CustomID]InputComponent, len(iData.Components))
		for _, containerComponent := range iData.Components {
			for _, component := range containerComponent.Component.(ContainerComponent).Components() {
				if inputComponent, ok := component.(InputComponent); ok {
					d.Components[inputComponent.ID()] = inputComponent
				}
			}
		}
	}
	return nil
}

func (d ModalSubmitInteractionData) Get(customID CustomID) InputComponent {
	if component, ok := d.Components[customID]; ok {
		return component
	}
	return nil
}

func (d ModalSubmitInteractionData) TextComponent(customID CustomID) *TextInputComponent {
	component := d.Get(customID)
	if component == nil {
		return nil
	}
	if cmp, ok := component.(TextInputComponent); ok {
		return &cmp
	}
	return nil
}

func (d ModalSubmitInteractionData) Text(customID CustomID) *string {
	if component := d.TextComponent(customID); component != nil {
		return &component.Value
	}
	return nil
}
