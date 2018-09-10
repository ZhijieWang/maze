package participants

// DummyBot serves as a reference implementation for robot type participants
type DummyBot struct {
	observationChannel
}

// Announce borad cast its location, movement and observations to others via pre-defined channels.
func (b DummyBot) Announce() {

}

// Observe takes information and other observations from pre-defined channels
func (b DummyBot) Observe() {

}
