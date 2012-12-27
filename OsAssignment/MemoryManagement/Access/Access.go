package access
import(
	"strconv"
	"strings"
	"os"
	"fmt"
	//"time"
	"OsAssignment/MemoryManagement/Processes"
)
type Access struct{
	process uint64	//pid
	access  string	//r/w
	virtualAddress uint64	//address
}
var Done (chan bool)
var Skip (chan bool)
var Running uint64
func (t* Access) AddEntry(acc string){	//Adds an entry to the Access structure. It takes a line (from the input file), and sets all the structure's fields
	entry := strings.Split(acc,",")	//Splits the line into 3 parts (on comma), returns an array of strings: pid, r/w, address
	if len(entry)>3{	//If we get more than 3 fields on a line, then, ERROR
		fmt.Println("Entry had more than 3 fields")
		os.Exit(2)
	}
	var err error = nil	// err is a variable of type error , which is a primitive type
	t.process,err = strconv.ParseUint(entry[0],10,64)	//Converts the pid, which is taken as decimal(second parameter), whose length is restricted to 64 bits(3rd parameter) to uint
	
	if err != nil{	//If it is unable to convert, ERROR
		fmt.Println("Invalid Process ID",err)
		os.Exit(3)
	}
	
	if strings.Contains(entry[1],"W") || strings.Contains(entry[1],"w"){	//contains searches for second param, in first param
		t.access = "w"														// for simplicity, we're making all the chars to lowercase
	}else if strings.Contains(entry[1],"r") || strings.Contains(entry[1],"R"){
		t.access = "r"
	}else{
		fmt.Println("Invalid Access Specifier",entry[1])
		os.Exit(4)
	}

	t.virtualAddress , err  = strconv.ParseUint(entry[2], 16, 16)	//Takes the address, treats it as hexa, and restricts it to 16 bits. Returns uint 
	if err != nil{
		fmt.Println("I cannot understand the address ",err)
		os.Exit(5)
	}		
}


func GetNextMemAccess(p processes.ProcessList,requestPtr map[uint64][]Access, filled map[uint64] int)(page int64 ,process uint64,access string,virtualAddress uint64){

		
		var r = Running		//Gets the pid of the currently executing process
		var a Access	
		a= (requestPtr[r])[0]	//requestPtr is a key-value pair. Given a key, which is a pid, its value will be an array of accesses for that process.
								// Max number of accesses per process would be 100
								//Here, we are getting the first access of that array
		for a.process == 0{	//If a process is done with all its accesses, then:
				p.Delete(r)	//Delete this process from the system
				fmt.Println(r," DONE WITH EXECUTION")
				Running = p.Get() // Get a new Process
				Skip <- true
				r = Running
				a= (requestPtr[r])[0]
				
			}
		page = int64(a.virtualAddress >>10) // Page size is 2^10; Returning the page number
		process = a.process
		access = a.access
		virtualAddress = a.virtualAddress
		requestPtr[r] = (requestPtr[r])[1:]	//Then, we consider the remaining processes, leaving out the first
		return

}