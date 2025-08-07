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
				NewStringSelectMenuOption("option1", "option1"),
				NewStringSelectMenuOption("option2", "option21"),
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
		NewLabel("label1", NewStringSelectMenu("select2", "select2",
			NewStringSelectMenuOption("option1", "option1"),
			NewStringSelectMenuOption("option2", "option2"),
		).WithID(17),
		).WithID(16),
	}

	var i int
	for c := range componentIter(components) {
		i++
		if c.GetID() != i {
			t.Errorf("expected component with ID %d, got %d", i, c.GetID())
		}
	}
	if i != 17 {
		t.Errorf("expected 17 components, got %d", i)
	}
}
