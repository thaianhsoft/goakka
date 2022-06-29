package goakka

type Actor interface{
	Receive(ctx *ActorContext)
}

type Receivefunc func(ctx *ActorContext)
func (r Receivefunc) Receive(ctx *ActorContext) {
	r(ctx)
}
