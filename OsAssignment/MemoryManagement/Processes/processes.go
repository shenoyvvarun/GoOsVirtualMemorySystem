package processes

import "fmt"
import "os"
import "container/list"	// Doubly-Linked List Interface

// This file contains a list of all processes that are, otherwise, ready to run

type ProcessList struct{
	PidList * list.List	// List of Pids
	done bool	// Default value is false
}


func (p *ProcessList) Add(pid uint64){
		if !p.done{
			p.done = true
			p.PidList = list.New()	// Make the list here
		}
		var unique = true	// Indicate uniqueness of Pids
		// Uniquely add this pid to the list
		for e:= p.PidList.Front();e!=nil;e= e.Next(){
			if pid == e.Value{
				unique = false
			}
		}
		if unique{
			p.PidList.PushBack(pid)
		}
}
func (p *ProcessList) Get()(pid uint64){

	v :=p.PidList.Front()
	if v == nil{
		fmt.Println("Done with Execution OF ALL PROCESSES")
		os.Exit(7)
	}
	p.PidList.Remove(v)
	p.PidList.PushBack((v.Value).(uint64))	// simulates RoundRobin
	return (v.Value).(uint64)
}

func (p *ProcessList) Delete(pid uint64){	//This method is called only when its array of accesses is Empty
		var del bool = true
		for e:= p.PidList.Front();e!=nil;e= e.Next(){
			if pid == e.Value{
				p.PidList.Remove(e)
				del = false
			}
		}
		if del{
			fmt.Println("Cannot delete a pid that doesn't exist")
			os.Exit(11)
		}
}

func (p  *ProcessList) Print(pidRunning uint64){
	for e:= p.PidList.Front();e!=nil;e= e.Next(){
			if pidRunning != e.Value{
				fmt.Printf("%v->",(e.Value).(uint64))
			}
		}
}