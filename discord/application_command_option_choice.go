package discord

import (
	"math"

	"github.com/disgoorg/validate"
)

type ApplicationCommandOptionChoice interface {
	ChoiceName() string
	ChoiceNameLocalizations() map[Locale]string
	validate.Validator
	applicationCommandOptionChoice()
}

var _ ApplicationCommandOptionChoice = (*ApplicationCommandOptionChoiceInt)(nil)

type ApplicationCommandOptionChoiceInt struct {
	Name              string            `json:"name"`
	NameLocalizations map[Locale]string `json:"name_localizations,omitempty"`
	Value             int               `json:"value"`
}

func (ApplicationCommandOptionChoiceInt) applicationCommandOptionChoice() {}
func (c ApplicationCommandOptionChoiceInt) ChoiceName() string {
	return c.Name
}
func (c ApplicationCommandOptionChoiceInt) ChoiceNameLocalizations() map[Locale]string {
	return c.NameLocalizations
}
func (c ApplicationCommandOptionChoiceInt) Validate() error {
	min := math.SmallestNonzeroFloat64
	max := math.MaxFloat64
	return validate.Validate(
		validate.Value(c.Value, validate.NumberRange[int](int(min), int(max))),
		applicationCommandOptionChoiceValidator(c),
	)
}

var _ ApplicationCommandOptionChoice = (*ApplicationCommandOptionChoiceString)(nil)

type ApplicationCommandOptionChoiceString struct {
	Name              string            `json:"name"`
	NameLocalizations map[Locale]string `json:"name_localizations,omitempty"`
	Value             string            `json:"value"`
}

func (ApplicationCommandOptionChoiceString) applicationCommandOptionChoice() {}

func (c ApplicationCommandOptionChoiceString) ChoiceName() string {
	return c.Name
}

func (c ApplicationCommandOptionChoiceString) ChoiceNameLocalizations() map[Locale]string {
	return c.NameLocalizations
}

func (c ApplicationCommandOptionChoiceString) Validate() error {
	return validate.Validate(
		validate.Value(c.Value, validate.Required[string], validate.StringRange(1, ApplicationCommandOptionStringChoiceValueMaxLength)),
		applicationCommandOptionChoiceValidator(c),
	)
}

var _ ApplicationCommandOptionChoice = (*ApplicationCommandOptionChoiceInt)(nil)

type ApplicationCommandOptionChoiceFloat struct {
	Name              string            `json:"name"`
	NameLocalizations map[Locale]string `json:"name_localizations,omitempty"`
	Value             float64           `json:"value"`
}

func (ApplicationCommandOptionChoiceFloat) applicationCommandOptionChoice() {}

func (c ApplicationCommandOptionChoiceFloat) ChoiceName() string {
	return c.Name
}

func (c ApplicationCommandOptionChoiceFloat) ChoiceNameLocalizations() map[Locale]string {
	return c.NameLocalizations
}

func (c ApplicationCommandOptionChoiceFloat) Validate() error {
	return applicationCommandOptionChoiceValidator(c).Validate()
}

func applicationCommandOptionChoiceValidator(c ApplicationCommandOptionChoice) validate.Validator {
	return validate.Combine(
		validate.Value(c.ChoiceName(), validate.Required[string], validate.StringRange(1, ApplicationCommandOptionChoiceNameMaxLength)),
		validate.Map(c.ChoiceNameLocalizations(), validateLocalizations(1, ApplicationCommandOptionChoiceNameMaxLength)),
	)
}

const (
	ApplicationCommandOptionChoiceNameMaxLength        = 100
	ApplicationCommandOptionStringChoiceValueMaxLength = 100
)
