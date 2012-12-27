package mmu

import(
		"strings"
		"io"
		"fmt"
		"os"
		"OsAssignment/MemoryManagement/Access"
		"OsAssignment/PageTable"
		"OsAssignment/Queue"
		"OsAssignment/Frame"
		"strconv"
		"time"
		"OsAssignment/MemoryManagement/Processes"
	)

var MemSize int64
var translationDone = make(chan bool,1)
var doTranslation   = make(chan SendData,1)
var Processes = processes.ProcessList{}
var frames frame.FreeFrame
var pTablePtr map[uint64] pageTable.PageTable
var bQ	queue.BlockedQueue
var rQ	queue.ReadyQueue
type SendData struct {
	pid uint64
	page int64	
}

func MemoryMgmtHardwareSimulation(fr frame.FreeFrame){
	frames = fr
	/* Use  the pid to get the []access.Access element of that particular pid
		, this array denotes all the remaining requests pertaining to that process */

	requestPtr := make(map[uint64][]access.Access)																	
	filled := make( map[uint64] int) // The value tells us the index in the []acess.Access that has been processed 
	pTablePtr = make(map[uint64]pageTable.PageTable) // Pagetable per Process
	

	f,err := os.Open("C:/Users/Varun/Desktop/mygo/bin/input.txt")
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	//Read All accesses from file
	var e error
	s := make([]string,3) 
	temp := make([]string,2)
	n:=0
	for e != io.EOF {
		n,e = fmt.Fscanf(f,"%s%s%s",&s[0],&s[1],&s[2])	
		if n == 3{
			temp = strings.Split(s[0],",")	// Remove commas
			processid,err := strconv.ParseUint(temp[0],10,64) // The first field is assumed to be pid
			if err != nil{
				fmt.Println(err);
				os.Exit(7);
			}
			Processes.Add(processid)
			_,ok := requestPtr[processid]
			//Denotes the 1st request of that particular process
			if ok == false{
				requestPtr[processid] = make([]access.Access,100) //Fix this	
				filled[processid] =0
			}
			// Get the [] Accesses for a particular pid, and get the latest elem
			//for that particular array element,add the REQUEST READ
			
			// Len function cannot be used because, it will always return 100
			// We, here, want to get the number of entries in the Array of Accesses
			// and hence, we do not use the len function, and so, we use the
			// filled map
			
			((requestPtr[processid])[filled[processid]]).AddEntry(strings.Join(s,""))
			filled[processid]++
			
		}
	}
		
		access.Running = Processes.Get()
	for {
			fmt.Println("---------------------------------------------------------------------");
			page,pid,access,virtualAddress :=access.GetNextMemAccess(Processes,requestPtr,filled)
			realAddress,ok := VirtAddr2RealAddr(pid,virtualAddress,page)
			if ok==false {	// Page is NOT in Memory
				// Call the below method to get a new page into memory
				go OSMemoryMgmtSignalHandler()
				doTranslation<- SendData{pid,page}
				<-translationDone;	// The whole method waits until a page has been got into memory
			}
			realAddress,ok = VirtAddr2RealAddr(pid,virtualAddress,page)
			if ok == true {
				fmt.Println("\n",realAddress," issued on adress Lines")
				_,present :=pTablePtr[pid]
				if !present{
					pTablePtr[pid] = pageTable.PageTable{make([]int64,1<<6)}
				}
				pTablePtr[pid] = pTablePtr[pid].UseEntry(pid, page, access);	// Sets M bit to 1 if Write
           		fmt.Printf("Program %v accessing %v in mode %v\n", pid, virtualAddress, access); 
			}		
		}	

}
// Here we will make use of maps of ptable for each process
func VirtAddr2RealAddr(pid uint64, virtualAddress uint64,page int64)(realAddress uint64,ok bool) {
	
	_,ok = pTablePtr[pid]	// This indicates whether a PageTable exists for a Process
	// If not, then:
	if ok == false{
		pTablePtr[pid] = pageTable.PageTable{make([]int64,1<<6)}
	}
	frame,ok := pTablePtr[pid].GetEntry(pid,page)	// Check for presence of page and return the frame number if present  	
	if ok ==  true{
		// get bottom 10 bits, add to frame numbers' 6 bits, and get a 16 bit real address 
		realAddress = (virtualAddress - ((virtualAddress >> 10)<<10)) + uint64((frame<<10))
	}/*else{
		realAddress = 0
	}*/
	return		
}

// Call this method if there is a PageFault

func OSMemoryMgmtSignalHandler( ){

	var a SendData
	a = <-doTranslation
	bQ.Push(a.pid)
	fmt.Println("****BLOCKED QUE****");
	bQ.Print()						// Print the blocked queue
	fmt.Println("\n")
	page_out,frame_out,pid_out := getSwapCandidate(a.page,a.pid)
	// If modified bit is set then only possibilty is 3
	if page_out>0{
		if((pTablePtr[pid_out]).Ptable[page_out] >> 14) == 3{
			fmt.Printf("Swapping out proc %d:pg %d from frame %d\n", pid_out, page_out, frame_out,)
			time.Sleep(1000 * time.Millisecond)		// simulate disk write
			// Set the P Bit of the to be swapped out process to 0
			pTablePtr[pid_out].Ptable[page_out] = pTablePtr[pid_out].Ptable[page_out] - (((pTablePtr[pid_out]).Ptable[page_out] >> 15)<<15)
			//fmt.Printf("\n%b\n",pTablePtr[pid_out].Ptable[page_out])
		}
	}
	
	// Loading a new Process!!
	
	fmt.Printf("Loading page proc %v :page %v into frame %v \n", a.pid, a.page, frame_out);
	time.Sleep(1000 * time.Millisecond)
	pTablePtr[a.pid] = pTablePtr[a.pid].SetEntry(a.pid, a.page, frame_out)
	
	if a.pid == bQ.Pop() {
		rQ.Push(a.pid)
	}else{
		fmt.Println("FATAL  ERROR: The blocked process was never resumed");
		os.Exit(10)
	}
	fmt.Println("****READY QUEUE****"); 
	Processes.Print(access.Running)	// Print The runnable Process
	rQ.Print() 
	// These 2 above prints,logically form the ready que
	rQ.Pop()
	fmt.Println("\n\n")
	frames.Print();
	fmt.Println("\n\n")
	translationDone <- true	
	
}

// Get the next candidate to process

func getSwapCandidate(pg int64,pi uint64)(page,frame int64,pid uint64){
	var f []int
	page = -1
	frame = 0
	f = frames.FrameList
	for i,v := range f{
		if v == 0{							// Not used frame is given higher precedence.
			frame = int64(i)				// Store frame num
			frames.FrameList[i] = int(pg) 	// Indicate That the frame is used
			frames.PidList[i] = int(pi)
			return	
		}		
	}
	
	frame = int64(frames.Next)				// Gives the next frame to be swapped
	page = int64(frames.FrameList[frames.Next])
	pid = uint64(frames.PidList[frames.Next])
	frames.FrameList[frames.Next] = int(pg) // Indicate That the frame is used
	frames.PidList[frames.Next] = int(pi)
	frames.Next+=1
	if frames.Next == MemSize{
		frames.Next = 0
	}
	return
}

// Executed 1) The Executing Process goes into completion
//			2) Every 5 sec, to give every process a chance to execute

func OSSimulation(){

	for{
	
		FiveSecTimer := time.NewTimer(time.Second * 5)
		select{
			case <-FiveSecTimer.C: access.Running = Processes.Get()
			case <-access.Skip :
		}
		fmt.Println("OS kicked off")
		
	}
}