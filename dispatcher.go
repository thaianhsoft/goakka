package goakka

import (
	"fmt"
	"time"
)

type Dispatcher struct {
	pool *Pool
}

func (d *Dispatcher) RunnableProcessFunc(runnable *RunnableBehaviour) ThreadBehaviour {
	d.pool.GetCondVariableLocker().L.Lock()
	for d.pool.GetQueue().Empty() {
		d.pool.GetCondVariableLocker().Wait()
	}
	tb := &RunningBehaviour{
		Job: d.pool.GetQueue().PopFront(),
	}
	d.pool.GetCondVariableLocker().L.Unlock()
	d.pool.GetCondVariableLocker().Signal()
	return tb
}

func (d *Dispatcher) RunningProcessFunc(running *RunningBehaviour) ThreadBehaviour {
	switch message := running.Job.(type) {
	case *DefaultMessage:
		// write process func here
		receiverActorCtx := ActorSystemManager.getChildContext(message.receiver)
		if receiverActorCtx != nil {
			receiverActorCtx.mailBox.PushBack(message)
			ActorSystemManager.GetExecutor().GetPool().PushBackJob(message)
		}
		return &RunnableBehaviour{threadId: running.GetId()}
	case *CoodinateMessage:
		return &WaitCoodinateBehaviour{threadId: running.GetId(), DurationTimeout: message.GetDurationTimeout(), Job: running.Job}
	default:
		return nil
	}
}

func (d *Dispatcher) WaitCoodinateBehaviour(waitCoodinate *WaitCoodinateBehaviour) ThreadBehaviour {
	timer := time.NewTimer(*waitCoodinate.DurationTimeout)
	sentSignalExecutor := false
	for {
		select {
		case <-timer.C:
			fmt.Println("timeout")
		default:
			if waitCoodinate.Job.(*CoodinateMessage).SignalEnd() {
				receiverCtx := ActorSystemManager.getChildContext(waitCoodinate.Job.(Job).getReceiver())
				if receiverCtx != nil {
					defaultMsg := &DefaultMessage{
						receiver: waitCoodinate.Job.(Job).getReceiver(),
						sender:   waitCoodinate.Job.(Job).getSender(),
						data:     waitCoodinate.Job.(Job).getData(),
					}
					receiverCtx.mailBox.PushBack(defaultMsg)
					waitCoodinate.Job.(*CoodinateMessage).wait <- struct{}{}
					ActorSystemManager.executor.GetPool().PushBackJob(defaultMsg)
				}
				return &RunnableBehaviour{threadId: waitCoodinate.GetId()}
			} else {
				if !sentSignalExecutor {
					waitCoodinate.Job.(*CoodinateMessage).SendSignal(signalDispatcherThread)
					ActorSystemManager.executor.GetPool().PushFrontJob(waitCoodinate.Job)
					sentSignalExecutor = true
				}
			}
		}
	}
}

func (d *Dispatcher) GetPool() *Pool {
	return d.pool
}

func NewDispatcher(threadSize int) *Dispatcher {
	d := &Dispatcher{}
	d.pool = NewPool(threadSize, d)
	defer d.pool.spawnThread()
	return d
}




