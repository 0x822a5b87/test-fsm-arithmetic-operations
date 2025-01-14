package token

type charStream struct {
	data  string
	index int
}

func (cs *charStream) nextEvent() tokenizerEvent {
	return cs.getNextEvent(true)
}

func (cs *charStream) peekEvent() tokenizerEvent {
	return cs.getNextEvent(false)
}

func (cs *charStream) getNextEvent(isIncIndex bool) tokenizerEvent {
	if cs.index >= len(cs.data) {
		return null
	}
	var ch tokenizerEvent
	for {
		ch = tokenizerEvent(cs.data[cs.index])
		if ch == space || ch == tab || ch == r || ch == newLine {
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
