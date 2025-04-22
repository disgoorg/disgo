package discord

import "testing"

func TestComponentIter(t *testing.T) {
	components := []LayoutComponent{
		NewActionRow(
			NewPrimaryButton("button1", "button1").WithID(2),
			NewSecondaryButton("button2", "button2").WithID(3),
		).WithID(1),
		NewSmallSeparator().WithID(4),
		NewActionRow(
			NewStringSelectMenu("select1", "select1",
				StringSelectMenuOption{Label: "option1", Value: "option1"},
				StringSelectMenuOption{Label: "option2", Value: "option21"},
			).WithID(6),
		).WithID(5),
		NewSection(
			NewTextDisplay("text1").WithID(8),
			NewTextDisplay("text2").WithID(9),
		).WithID(7),
		NewContainer(
			NewActionRow(
				NewPrimaryButton("button2", "button2").WithID(12),
			).WithID(11),
			NewSection(
				NewTextDisplay("text3").WithID(14),
				NewTextDisplayf("text%d", 4).WithID(15),
			).WithID(13),
		).WithID(10),
	}

	i := 1
	for c := range componentIter(components) {
		if c.GetID() != i {
			t.Errorf("expected component with ID %d, got %d", i, c.GetID())
		}
		i++
	}
}
