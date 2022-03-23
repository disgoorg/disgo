package discord

import "github.com/disgoorg/disgo/json"

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

func (d ModalSubmitInteractionData) Component(customID CustomID) (InputComponent, bool) {
	component, ok := d.Components[customID]
	return component, ok
}

func (d ModalSubmitInteractionData) TextInputComponent(customID CustomID) (TextInputComponent, bool) {
	if component, ok := d.Component(customID); ok {
		textInputComponent, ok := component.(TextInputComponent)
		return textInputComponent, ok
	}
	return TextInputComponent{}, false
}

func (d ModalSubmitInteractionData) OptText(customID CustomID) (string, bool) {
	if textInputComponent, ok := d.TextInputComponent(customID); ok {
		return textInputComponent.Value, true
	}
	return "", false
}

func (d ModalSubmitInteractionData) Text(customID CustomID) string {
	if text, ok := d.OptText(customID); ok {
		return text
	}
	return ""
}
