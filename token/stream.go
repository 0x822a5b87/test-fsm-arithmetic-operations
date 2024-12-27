package token

type charStream struct {
	data  string
	index int
}

func (cs *charStream) nextEvent() TokenizerEvent {
	return cs.getNextEvent(true)
}

func (cs *charStream) peekEvent() TokenizerEvent {
	return cs.getNextEvent(false)
}

func (cs *charStream) getNextEvent(isIncIndex bool) TokenizerEvent {
	if cs.index >= len(cs.data) {
		return Null
	}
	var ch TokenizerEvent
	for {
		ch = TokenizerEvent(cs.data[cs.index])
		if ch == Space || ch == Tab || ch == Return || ch == NewLine {
			cs.index++
		} else {
			break
		}
	}
	if isIncIndex {
		cs.index++
	}
	return ch
}
