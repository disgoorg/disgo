package discord

type ApplicationRoleConnectionMetadata struct {
	Type                     ApplicationRoleConnectionMetadataType `json:"type"`
	Key                      string                                `json:"key"`
	Name                     string                                `json:"name"`
	NameLocalizations        map[Locale]string                     `json:"name_localizations,omitempty"`
	Description              string                                `json:"description"`
	DescriptionLocalizations map[Locale]string                     `json:"description_localizations,omitempty"`
}

type ApplicationRoleConnectionMetadataType int

const (
	IntegerLessThanOrEqual ApplicationRoleConnectionMetadataType = iota + 1
	IntegerGreaterThanOrEqual
	IntegerEqual
	IntegerNotEqual
	DateTimeLessThanOrEqual
	DateTimeGreaterThanOrEqual
	BooleanEqual
	BooleanNotEqual
)
