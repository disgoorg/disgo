package oauth2

import "log/slog"

var (
	_ StateController = (*defaultStateController)(nil)
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
	cfg := defaultStateControllerConfig()
	cfg.apply(opts)

	states := newTTLMap(cfg.MaxTTL)
	for state, url := range cfg.States {
		states.put(state, url)
	}

	return &defaultStateController{
		states:       states,
		newStateFunc: cfg.NewStateFunc,
		logger:       cfg.Logger,
	}
}

type defaultStateController struct {
	logger       *slog.Logger
	states       *ttlMap
	newStateFunc func() string
}

func (c *defaultStateController) NewState(redirectURI string) string {
	state := c.newStateFunc()
	c.logger.Debug("new state: %s for redirect uri", slog.String("state", state), slog.String("redirect_uri", redirectURI))
	c.states.put(state, redirectURI)
	return state
}

func (c *defaultStateController) UseState(state string) string {
	uri := c.states.get(state)
	if uri == "" {
		return ""
	}
	c.logger.Debug("using state: %s for redirect uri", slog.String("state", state), slog.String("redirect_uri", uri))
	c.states.delete(state)
	return uri
}
