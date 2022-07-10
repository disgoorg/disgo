package oauth2

var (
	_ StateController = (*stateControllerImpl)(nil)
)

// StateController is responsible for generating, storing and validating states.
type StateController interface {
	// GenerateNewState generates a new random state to be used as a state.
	GenerateNewState(redirectURI string) string

	// ConsumeState validates a state and returns the redirect url or nil if it is invalid.
	ConsumeState(state string) string
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
	}
}

type stateControllerImpl struct {
	states       *ttlMap
	newStateFunc func() string
}

func (c *stateControllerImpl) GenerateNewState(redirectURI string) string {
	state := c.newStateFunc()
	c.states.put(state, redirectURI)
	return state
}

func (c *stateControllerImpl) ConsumeState(state string) string {
	uri := c.states.get(state)
	if uri == "" {
		return ""
	}
	c.states.delete(state)
	return uri
}
