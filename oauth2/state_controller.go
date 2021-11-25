package oauth2

import "github.com/DisgoOrg/disgo/internal/insecurerandstr"

var _ StateController = (*stateControllerImpl)(nil)

// StateController is responsible for generating, storing and validating states
type StateController interface {
	// GenerateNewState generates a new random state to be used as a state
	GenerateNewState(redirectURI string) string

	// ConsumeState validates a state and returns the redirect url or nil if it is invalid
	ConsumeState(state string) *string
}

// NewStateController returns a new empty StateController
func NewStateController() StateController {
	return NewStateControllerWithStates(map[string]string{})
}

// NewStateControllerWithStates returns a new StateController with the given states
func NewStateControllerWithStates(states map[string]string) StateController {
	return &stateControllerImpl{states: states}
}

type stateControllerImpl struct {
	states map[string]string
}

func (c *stateControllerImpl) GenerateNewState(redirectURI string) string {
	state := insecurerandstr.RandStr(32)
	c.states[state] = redirectURI
	return state
}

func (c *stateControllerImpl) ConsumeState(state string) *string {
	uri, ok := c.states[state]
	if !ok {
		return nil
	}
	delete(c.states, state)
	return &uri
}
