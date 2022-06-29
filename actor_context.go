package goakka

import (
	"log"
	"reflect"
	"time"
)

type ActorContext struct {
	process    process
	mailBox    *Queue
	actor      Actor
	latestMail any
}

func newActorContext(process process, actor Actor) *ActorContext {
	return &ActorContext{
		process: process,
		mailBox: NewQueue(),
		actor:   actor,
	}
}

func (a *ActorContext) Send(pid Pid, data any) {
	receiverProcess := pid.ToProcess()
	defaultMsg := &DefaultMessage{
		receiver: receiverProcess,
		sender: a.process,
		data: data,
	}
	ActorSystemManager.GetDispatcher().GetPool().PushBackJob(defaultMsg)
}

func (a *ActorContext) executeMail() {
	if !a.mailBox.Empty() {
		a.latestMail = a.mailBox.PopFront()
		a.actor.Receive(a) // receive call message()
	}
}

func (a *ActorContext) Message() (sender Pid, data any) {
	if a.latestMail != nil {
		sender, data = a.latestMail.(Job).getSender().ToPid(), a.latestMail.(Job).getData()
	}
	return
}

func (a *ActorContext) Self() Pid {
	return a.process.ToPid()
}

func (a *ActorContext) Log(senderPid, data any) {
	rdata := reflect.Indirect(reflect.ValueOf(data))
	ractor := reflect.Indirect(reflect.ValueOf(a.actor))
	dataType := rdata.Type().Name()
	actorName := ractor.Type().Name()
	log.Printf("LOG [%v]-Actor Received from sender [PID=%v], data [Type=%v]: %v", actorName, senderPid, dataType, data)
}

func (a *ActorContext) Ask(pid Pid, data any, timeout time.Duration){
	receiverProcess := pid.ToProcess()
	coodinateMsg := &CoodinateMessage{
		receiver:         receiverProcess,
		sender:           a.process,
		data:             data,
		SignalFromOthers: 0,
		Timeout:          timeout,
	}
	ActorSystemManager.GetDispatcher().GetPool().PushFrontJob(coodinateMsg)
}


