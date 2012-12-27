package frame

import(
	"fmt"	//this imports a package "fmt"
)

//all var by Default initialization to zero
 
type FreeFrame struct{	//Creating a type, which is a structure type, named "FreeFrame"
	FrameList []int		//An array, whose indices indicate the frame numbers; obj.FrameList[0] => this gives page present in frame number 0
	Next int64		//This indicates which frame number is to be replaced next . 
	PidList []int	//This gives the pid of the process present in given frame, obj.PidList[0] => This gives pid of process which has its page in FrameList[0]
}
func (f *FreeFrame) Populate(mem int64){	//This is a method of the structure defined above. mem is the size of memory in KB
	f.FrameList = make([]int, mem)		//We know that the page size is 1K, number of frames = mem/1 = mem. make is like new. So, this creates an array of size mem
										//This implies that memory now has mem number of frames
	f.PidList = make([]int, mem)		//This creates an int array of size mem
}
func (f *FreeFrame) Print(){	//This method will print the pid list

	for i := range f.FrameList{	//For all processes in the Framelist, or in other words, while the FrameList is not empty, do:
		fmt.Printf("|%d",f.PidList[i])	//Print the process id that has its page in this frame
	}
	fmt.Printf("|");
}

//make is applicable only to objects

//new to primitive types
//no while, do while loops
// range is an operator , which returns two values.(index and value)
