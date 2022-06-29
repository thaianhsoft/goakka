package goakka

import (
	"fmt"
)

type Executor struct {
	pool *Pool
}

func (d *Executor) RunnableProcessFunc(runnable *RunnableBehaviour) ThreadBehaviour {
	d.pool.GetCondVariableLocker().L.Lock()
	for d.pool.GetQueue().Empty() {
		d.pool.GetCondVariableLocker().Wait()
	}
	tb := &RunningBehaviour{
		Job: d.pool.GetQueue().PopFront(),
	}
	fmt.Println(tb.Job.(Job).getData(), d.pool.GetQueue().back)
	defer func() {
		d.pool.GetCondVariableLocker().L.Unlock()
		d.pool.GetCondVariableLocker().Signal()
	}()
	return tb
}

func (d *Executor) RunningProcessFunc(running *RunningBehaviour) ThreadBehaviour {
	switch message := running.Job.(type) {
	case *DefaultMessage:
		// write process func here
		receiverActorCtx := ActorSystemManager.getChildContext(message.receiver)
		if receiverActorCtx != nil {
			if !receiverActorCtx.mailBox.Empty() {
				receiverActorCtx.executeMail()
			}
		}
		return &RunnableBehaviour{threadId: running.GetId()}
	case *CoodinateMessage:
		message.SendSignal(signalExecutorThread)
		return &RunnableBehaviour{threadId: running.GetId()}
	default:
		return nil
	}
}

func (d *Executor) WaitCoodinateBehaviour(waitCoodinate *WaitCoodinateBehaviour) ThreadBehaviour {
	return nil
}

func (d *Executor) GetPool() *Pool {
	return d.pool
}


func NewExecutor(threadSize int) *Executor {
	d := &Executor{}
	d.pool = NewPool(threadSize, d)
	defer d.pool.spawnThread()
	fmt.Println("ok executor created")
	return d
}




