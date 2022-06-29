package goakka

var ActorSystemManager *ActorSystem
type signalThread uint
const (
	signalDispatcherThread signalThread = iota
	signalExecutorThread signalThread = iota
)
type ActorSystem struct {
	*ActorContext
	globalProcess process
	dispatcher    *Dispatcher
	executor      *Executor
	childs        map[process]*ActorContext
}

func NewActorSystem() *ActorSystem {
	if ActorSystemManager != nil {
		return ActorSystemManager
	}
	ActorSystemManager = &ActorSystem{
		ActorContext:  newActorContext(0, nil),
		globalProcess: 1,
		dispatcher:    NewDispatcher(10),
		executor:      NewExecutor(10),
		childs: map[process]*ActorContext{},
	}
	return ActorSystemManager
}

func (a *ActorSystem) GetDispatcher() *Dispatcher {
	return a.dispatcher
}

func (a *ActorSystem) GetExecutor() *Executor {
	return a.executor
}

func (a *ActorSystem) ActorOf(fn func() Actor) *ActorContext {
	ctx := newActorContext(a.globalProcess, fn())
	if _, ok := a.childs[a.globalProcess]; !ok {
		a.childs[a.globalProcess] = ctx
		a.globalProcess++
		return ctx
	}
	return nil
}

func (a *ActorSystem) getChildContext(process process) *ActorContext {
	if _, ok := a.childs[process]; ok {
		return a.childs[process]
	}
	return nil
}


