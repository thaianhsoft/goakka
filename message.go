package goakka

import (
	"time"
)

type Job interface{
	getSender() process
	getReceiver() process
	getData() any
}

type DefaultMessage struct {
	receiver process
	sender   process
	data     any
}

func (d *DefaultMessage) getSender() process{
	return d.sender
}

func (d *DefaultMessage) getReceiver() process {
	return d.receiver
}

func (d *DefaultMessage) getData() any {
	return d.data
}

type CoodinateMessage struct {
	receiver         process
	sender           process
	data             any
	SignalFromOthers signalThread
	Timeout          time.Duration
	wait             chan struct{}
}

func (c *CoodinateMessage) getSender() process {
	return c.sender
}

func (c *CoodinateMessage) getReceiver() process {
	return c.receiver
}

func (c *CoodinateMessage) getData() any {
	return c.data
}

func (c *CoodinateMessage) SendSignal(signal signalThread) {
	if c.SignalFromOthers < 3 && (c.SignalFromOthers & (1<<signal)) == 0 {
		c.SignalFromOthers |= 1 << signal
	}
}


func (c *CoodinateMessage) SignalEnd() bool {
	return c.SignalFromOthers == 3
}

func (c *CoodinateMessage) OpenCoodinateChan() {
}


func (c *CoodinateMessage) GetDurationTimeout() *time.Duration {
	return &c.Timeout
}