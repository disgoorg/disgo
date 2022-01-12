package oauth2

import "github.com/DisgoOrg/disgo/internal/ttlmap"

var (
	_ StateController = (*stateControllerImpl)(nil)
)

// StateController is responsible for generating, storing and validating states
type StateController interface {
	// GenerateNewState generates a new random state to be used as a state
	GenerateNewState(redirectURI string) string

	// ConsumeState validates a state and returns the redirect url or nil if it is invalid
	ConsumeState(state string) string
}

// NewStateController returns a new empty StateController
func NewStateController(config *StateControllerConfig) StateController {
	if config == nil {
		config = &DefaultStateControllerConfig
	}

	var states = ttlmap.New[string, string](config.MaxTTL)
	for state, url := range config.States {
		states.Put(state, url)
	}

	return &stateControllerImpl{
		states:       states,
		newStateFunc: config.NewStateFunc,
	}
}

type stateControllerImpl struct {
	states       *ttlmap.Map[string, string]
	newStateFunc func() string
}

func (c *stateControllerImpl) GenerateNewState(redirectURI string) string {
	state := c.newStateFunc()
	c.states.Put(state, redirectURI)
	return state
}

func (c *stateControllerImpl) ConsumeState(state string) string {
	uri, ok := c.states.Get(state)
	if !ok {
		return ""
	}
	c.states.Delete(state)
	return uri
}
