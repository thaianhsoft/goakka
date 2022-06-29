package example

type ProductActor struct {
}


type ProductActorSignUpCommand struct {
}

func (p *ProductActorSignUpCommand) OnCommandHandler() any {
	return &ProductActorSignUpEvent{}
}

type ProductActorUpdateCommand struct {
}

func (p *ProductActorUpdateCommand) OnCommandHandler() any {
	return &ProductActorUpdateEvent{}
}

type ProductActorSignUpEvent struct {
}

type ProductActorUpdateEvent struct {

}

