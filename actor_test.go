package goakka

import (
	"fmt"
	"testing"
	"time"
)

type GinActor struct {

}

type GinMessage struct {
	GinItem string
}

func (g *GinActor) Receive(ctx *ActorContext) {
	sender, msg := ctx.Message()
	switch msg.(type) {
	case *GinMessage:
		ctx.Log(sender, msg)
		ctx.Send(sender, &ProductMessage{
			Items: []string{"1", "2"},
		})
	}
}

type ProductActor struct {

}

type ProductMessage struct {
	Items []string
}

func (p *ProductActor) Receive(ctx *ActorContext) {
	sender, msg := ctx.Message()
	switch msg.(type) {
	case *ProductMessage:
		ctx.Log(sender, msg)
	}
}

func TestActor(t *testing.T) {
	manager := NewActorSystem()
	fmt.Println(manager)
	product_actor := manager.ActorOf(func() Actor {
		return &ProductActor{}
	})
	gin_actor := manager.ActorOf(func() Actor {
		return &GinActor{}
	})
	product_actor.Send(gin_actor.Self(), &GinMessage{
		GinItem: "gin from product",
	})
	gin_actor.Ask(product_actor.Self(), &ProductMessage{Items: []string{"ask type"}}, 2*time.Second)
	for {

	}
}
