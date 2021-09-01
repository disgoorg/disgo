package oauth2

var _ StateController = (*stateControllerImpl)(nil)

type StateController interface {
	GenerateNewState(redirectURI string) string
	ConsumeState(state string) *string
}

func NewStateController() StateController {
	return &stateControllerImpl{states: map[string]string{}}
}

type stateControllerImpl struct {
	states map[string]string
}

func (c *stateControllerImpl) GenerateNewState(redirectURI string) string {
	state := RandStr(32)
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
