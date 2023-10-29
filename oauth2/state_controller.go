package oauth2

import "github.com/disgoorg/log"

var (
	_ StateController = (*stateControllerImpl)(nil)
)

// StateController is responsible for generating, storing and validating states.
type StateController interface {
	// NewState generates a new random state to be used as a state.
	NewState(redirectURI string) string

	// UseState validates a state and returns the redirect url or nil if it is invalid.
	UseState(state string) string
}

// NewStateController returns a new empty StateController.
func NewStateController(opts ...StateControllerConfigOpt) StateController {
	config := DefaultStateControllerConfig()
	config.Apply(opts)

	states := newTTLMap(config.MaxTTL)
	for state, url := range config.States {
		states.put(state, url)
	}

	return &stateControllerImpl{
		states:       states,
		newStateFunc: config.NewStateFunc,
		logger:       config.Logger,
	}
}

type stateControllerImpl struct {
	logger       log.Logger
	states       *ttlMap
	newStateFunc func() string
}

func (c *stateControllerImpl) NewState(redirectURI string) string {
	state := c.newStateFunc()
	c.logger.Debugf("new state: %s for redirect uri: %s", state, redirectURI)
	c.states.put(state, redirectURI)
	return state
}

func (c *stateControllerImpl) UseState(state string) string {
	uri := c.states.get(state)
	if uri == "" {
		return ""
	}
	c.logger.Debugf("using state: %s for redirect uri: %s", state, uri)
	c.states.delete(state)
	return uri
}
