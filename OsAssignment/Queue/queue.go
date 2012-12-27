package queue

import "fmt"

type q struct{
	queue []uint64
	
}

func (t *q)Push(v uint64){
	t.queue = append(t.queue,v)
}

func (t *q)Pop()(v uint64){
	v = t.queue[len(t.queue)-1]
	t.queue = t.queue[0:len(t.queue)-1]	// Re-slicing
	return
}

func(t *q) Print(){
	for i := range t.queue{
		fmt.Printf("%d->",t.queue[i])
	}
}

func (t *q)Random()(v uint64){
	if len(t.queue) == 0{
		return 0
	}
	v = t.queue[0]
	t.queue = t.queue[1:]
	return
}
type ReadyQueue struct{	//This is like a typedef. This creates a q, calling it ReadyQueue
	q
}	// similar to ReadyQueue.q.methodName(); 
type BlockedQueue struct{
	q
}