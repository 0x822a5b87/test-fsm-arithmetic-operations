package fsm

type State int

type Event interface {
	EventType() int
}

// Action 这里应该是一个函数，暂时先写成这样
type Action func(event Event)

type Fsm[S State, E Event, A Action] interface {
	// Exec 根据当前 Fsm 收到的事件执行对应的 Action
	Exec(event E)
	AddAction(event E, action Action)
}

//func NewFms[S State, E Event](initState S) *Fsm[S, E] {
//	return &Fsm[S, E]{
//		State: initState,
//	}
//}
