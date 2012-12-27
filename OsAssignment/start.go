package main

import(
	"fmt"
	"OsAssignment/MemoryManagement"
	"OsAssignment/Frame"
	//"time"
	"OsAssignment/MemoryManagement/Access"
	)
	


func main(){
	fmt.Println("Please Enter the memory size in KiloBytes");
	fmt.Scanf("%d",&mmu.MemSize)	// mmu is a package
	var f frame.FreeFrame			//frame is a package, FreeFrame is a struct, indicates all the frames(free & non free)
	f.Populate(mmu.MemSize)
	access.Skip = make(chan bool)
	// pass all frames to the function below
	go mmu.MemoryMgmtHardwareSimulation(f);
	go mmu.OSSimulation()
	//Every 5 sec the OS(simulation) chooses a new process
	<-access.Done
	
}