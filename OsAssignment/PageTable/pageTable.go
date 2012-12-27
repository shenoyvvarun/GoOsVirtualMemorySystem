package pageTable

//Every process has its own version of the PageTable

type PageTable struct{	
	Ptable []int64 //PTE : frame no:[0:13], m bit is 14th bit, p bit is 15th bit 
}

// Set the Entry in The PTE of that process

func (p PageTable) SetEntry(pid uint64,page int64,frame int64)(PageTable){	//This puts frame number "frame" into "page" within PageTable of the process with pid
	p.Ptable[page] =0; //clear previous entry
	p.Ptable[page] = (1<<15) + frame	//First, it sets the P bit, which is the leftmost bit. Then, the frame number is put onto the remaining bits
	return p	//Return the page table
}

func (p PageTable) GetEntry(pid uint64,page int64) (frame int64,ok bool) {	//This returns the frame number of given page
		frame = p.Ptable[page] - ((p.Ptable[page]>>14)<<14)	//This returns the frame number, i.e. the bits from 0 t0 13 
		ok = ((p.Ptable[page]>>15) ==1)	//this is true if present bit is 1, i.e., if the page currently exists in some frame, i.e. no page fault
		return
}

func (p PageTable) UseEntry(pid uint64,page int64,read string)(PageTable){
	if read == "w"{		//If we are writing, then, the M bit is set to 1
		p.Ptable[page]=  p.Ptable[page] | 1<<14
	}
	// Once Modified bit wil be set to 1 and will be maintained. If the next time,
	// on the same page, if there is a read access, this bit will still
	// be 1 and hence, when swapping out, write-back happens.
	return p
}