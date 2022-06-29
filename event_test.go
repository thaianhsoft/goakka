package goakka

import (
	"github.com/thaianhsoft/goakka/behaviour"
	"github.com/thaianhsoft/goakka/example"
	"testing"
)

func TestEvent(t *testing.T) {
	b := &behaviour.BehaviourGen{}
	b.InitActorBehaviour(&example.ProductActor{})
}
