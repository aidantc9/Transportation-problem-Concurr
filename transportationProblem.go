package main//Aidan Charles

import (
"fmt"
"container/list"
"sync"
"os"
"log"
"bufio"
"strings"
"strconv"
)
type Boards struct{//holds all the different arrays representing the different boards
	boardC [][]int//this board holds the cost of each square
	supply []int //this board holds the supply of each factory
	demand []int//this board holds the demand of each warehouse 
	boardA [][]Cell//the is a bit of a all rounded board made up of cells but its main purpose is to hold the solution

}
type Point struct{//this structure just holds a position of a cell on a board in its row and column coordinates 
	row int
	col int 
}
type Cell struct {//this is the block in the board that holds all the main values like the amount of product it holds 
	cellC int
	Point 
	amount int

}
type Path struct{//This objects is used for paths across the board when finding the optimal solution
	path []Cell//this holds the (in order) the cells visited by the path 
	margC int//this is the maginal cost of the path
	changeRate int// this holds the amount the path if taken will reduce or increase the given cells on the path 
}
var wg sync.WaitGroup//Global waitgroup to wait for all the go routines to finish 
var chCounter int//kind of like a wait group used to check when all the paths have been calculated  
var pathCost int//Global variable for the total cost of a path 


func (bd *Boards) steppingStone(result chan Path){//this method goes through the board and finds a empty cell and calls the marginal cost method to find the path and its marginal cost 
	wg.Add(1)

	for i:=0;i<len(bd.supply);i++{//this loop goes through the board
		for j:=0;j<len(bd.demand);j++{
			if bd.boardA[i][j].amount==0 {
				test:=Cell{bd.boardC[i][j],Point{i,j},0}
				chCounter++//increase go routine counter 
					
				go bd.marginalCost(test,result)//concurently call marginal cost 
				

			}
		}
	}
		
	cheapestCost:=0//the cost of the cheapest path 
	cheapestPath:=Path{}//the cheapest path 
	posL:=list.New()//this list holds all the solutions found by the marginal cost method 
	for {
		
		sl:=<-result
		posL.PushBack(sl)
		if chCounter==0{
			break
		}

	}
		
			
	for x:=posL.Front();x!=nil;x=x.Next(){//goes through all the paths and finds the best one 
		temp:=x.Value.(Path)
		if temp.margC<cheapestCost{
			cheapestCost=temp.margC
			cheapestPath=temp
				

		}
	}
		
			
		if cheapestCost<0{//if a path are a negative maginal cost then it takes that path and edits the board accordingly 
			for i:=0;i<len(cheapestPath.path);i++{
					
				if i%2==0{
					bd.boardA[cheapestPath.path[i].row][cheapestPath.path[i].col].amount+=cheapestPath.changeRate
				}else{
					bd.boardA[cheapestPath.path[i].row][cheapestPath.path[i].col].amount-=cheapestPath.changeRate
				}
			}
				
		}

	
	
		
	

	
	if cheapestCost>=0{//restarts the algorith until it cant make the board any better 
		bd.steppingStone(result)
			
	}else{
		pathCost=bd.findTotalCost()//finds the total cost of the final board 
	}

	wg.Done()

}






func (bd *Boards) marginalCost (test Cell,result chan Path){//generates a path and finds it marginal cost it passes this information over the result channel 
	temp:= list.New()
	temp.PushBack(test)
	
	
	for i:=0;i<len(bd.boardA);i++{//adds all the possible cells on the path 
		for j:=0;j<len(bd.boardA[0]);j++{
			if bd.boardA[i][j].amount!=0{
				temp.PushBack(bd.boardA[i][j])
			}

		}
	}

	var sRow Cell//holds the cell in the same row as the current cell being worked on
   	var sCol Cell//same as sRow but for columns
	for {//this loop removes all the unessary parts of the path 
		count:=0
		for cell:=temp.Front();cell!=nil;cell=cell.Next(){
			c:=cell.Value.(Cell)
			sRow,sCol=bd.checkRC(c,temp)
			if sCol == (Cell{}) || sRow == (Cell{}){//if it does not have a cell beside the current cell being worked on then its not part of the path and is removed
				temp.Remove(cell)
				count++
			}
			
		}

		if count==0{//no more changes made therefore done and exits the loop
			break
		}
		
	
		
	}
	ans:=bd.fullPath(test,temp)//call the full path method which reorders the path found so it follows the order needed for this alg
	changeRate:=1000000//how much the cells will change 
	p:=Path{ans,0,0}//the path struct used to hold margc and change rate 
	for i:=0;i<len(p.path);i++{
 		if i%2==0{
 			p.margC+=p.path[i].cellC
 		}else{
 			p.margC-=p.path[i].cellC
 		}//finds marginalcost
 	
 		if changeRate>p.path[i].amount&&p.path[i].amount!=0{//finds the min amount and uses that as the change rate 
 			changeRate=p.path[i].amount
 			p.changeRate=changeRate
 		}
 	}
 	
 	
 	result<-p//passes the path over the channel
 	chCounter--//decrements to number of routines counter 




}
func (bd *Boards)fullPath(test Cell,pPath *list.List) []Cell {//finds the full path ie orders an unordered path into the right order 
	start:=test
	
	cPath:= make([]Cell, pPath.Len())
	cPath[0]=start
	var tempR Cell
	var tempC Cell
	for i:=0;i<len(cPath);i++{
		cPath[i]=start
		tempR,tempC=bd.checkRC(start,pPath)//finds the next element in the path this is possible because the order always alternates from row to col changes 
		if (i%2==0){
			start=tempR
		}else{
			start=tempC
		}
		
	}
	return cPath


}
func (bd *Boards) print(){//prints the board 
	for i:=0;i<len(bd.boardA);i++{
		for j:=0;j<len(bd.boardA[0]);j++{
			fmt.Print(" ")
			fmt.Print(bd.boardA[i][j].amount) 
			fmt.Print(" ")
		}
		fmt.Println()
	}
}
func (bd *Boards)findTotalCost() int {//finds the total cost of a board 
	sum:=0
	for i:=0;i<len(bd.boardC);i++{
		for j:=0;j<len(bd.boardC[0]);j++{
			sum+=bd.boardC[i][j]*bd.boardA[i][j].amount
		}
		
	}
	return sum
}




func (bd *Boards) checkRC(first Cell, path *list.List) (Cell,Cell) {//checks the row and col shared by the cell and finds the cells with elements in them same row, and same col
   	 var sRow Cell//stores the element in the same row 
   	 var sCol Cell//stores the element in the same col 
   	var currVal Cell//current cell being worked on
    for cell := path.Front(); cell != nil; cell = cell.Next() {//goes through the path and finds a cell that either share the same row or same col and stores that in the 2 variables above and returns them
        currVal= cell.Value.(Cell)
        if currVal == first {
        	continue
        }
        if currVal.row == first.row && sRow == (Cell{}) {
            sRow = currVal
        } else if currVal.col == first.col && sCol == (Cell{}) {
            sCol = currVal
        }
        if sCol != (Cell{}) && sRow != (Cell{}) {
            break
        }
        
    
    }
    return sRow,sCol
}
func (bd *Boards) checkDegn() bool{
	m:=len(bd.boardC)
	n:=len(bd.boardC[0])
	test:=m+n-1
	count:=0
	for i:=0;i<m;i++{
		for j:=0;j<n;j++{
			if bd.boardA[i][j].amount!=0{
				count++
			}
		}
	}
	if count<test{
		return true
	}
	return false
}


func main() {
	var inputData string//the string holding the file name of the input board 
	var initial string //holds the file name of the intial solution file 
	

	fmt.Printf("Please enter in the input data file name  ")
	fmt.Scanf("%s", &inputData)
	fmt.Printf("Please enter in the initial solution file name  ")
	fmt.Scanf("%s", &initial)
	file1, err1 := os.Open(inputData)
	file2, err2 := os.Open(initial)
	if err1 != nil {
    	log.Fatal(err1)
	}
	if err2 != nil {
     	log.Fatal(err2)
	}
	defer file1.Close()
	defer file2.Close()

	scanner1 := bufio.NewScanner(file1)
	scanner2 := bufio.NewScanner(file2)
	counter:=0
	jSize:=0//the number of rows
	var iSize int//num of cols
	check:=false
	costList:= list.New()//list that holds all the elements in the file that represent a cost value
	supplyList:= list.New()//hows all the supply values from the file
	demandList:=list.New()//same as above but for demand 
	amountList:=list.New()//used to hold cell data from the intial solution file 
	var num int //temp variable
	
	
	for scanner1.Scan() {//goes through the first file and stores all the values in their respective lists
		temp:=strings.Fields(scanner1.Text())
		if counter==0{
			jSize+= len(temp)-2
		}else if strings.TrimSpace(temp[0])=="DEMAND"{
			check=true
			iSize=counter-1
		}
		for i:=0;i<len(temp);i++{
			if i>0&&i<len(temp)-1&&counter>0&&check==false{
				num, _=strconv.Atoi(strings.TrimSpace(temp[i]))
				costList.PushBack(num)

			}
			if i==len(temp)-1&&counter>0&&check==false{
				num, _=strconv.Atoi(strings.TrimSpace(temp[i]))
				supplyList.PushBack(num)
			}
			if check&&i>0{
				num, _=strconv.Atoi(strings.TrimSpace(temp[i]))
				demandList.PushBack(num)

			}
		}

    	counter++
	}
	if err := scanner1.Err(); err != nil {
    log.Fatal(err)
	}
	
	var sup []int//all of the slices are used when converting the list into arrays or 2d arrays that can be used for the alg
	var dem []int
	var tmp [] int
	var tmp2 [] int
	var cst [][]int
	
	
	for cell := supplyList.Front(); cell != nil; cell = cell.Next(){//stores all the supplies in a slice 
		sup=append(sup,cell.Value.(int))

		
	}
	
	for cell := demandList.Front(); cell != nil; cell = cell.Next(){//same as above for demand 
		dem=append(dem,cell.Value.(int))
		
	}
	cnt:=0
	for cell := costList.Front(); cell != nil; cell = cell.Next(){//stores the values from the cost list into a 2d slice
		tmp=append(tmp,cell.Value.(int))//this temp is used because the cost slice is 2d
		cnt++
		if cnt==jSize{
			cst=append(cst,tmp)
			cnt=0
			tmp = tmp2
		}
		
	}
	
	cntI:=0
	check=false
	for scanner2.Scan() {//goes through the second file and creates the board with all the cell structures 
		temp:=strings.Fields(scanner2.Text())
		if strings.TrimSpace(temp[0])=="DEMAND"{
				check=true
		}
		for k:=0;k<len(temp);k++{//stores all the cells in a list that is later processed into a 2d array 
			str:=strings.TrimSpace(temp[k])
			
			if k>0&&k<len(temp)-1&&cntI>0&&check==false&&str!="-"{//non empty cells 
				num, _=strconv.Atoi(str)
				j:=k-1
				p:=Point{cntI-1,j}
				price:=cst[cntI-1][j]
				cell:=Cell{price,p,num}
				amountList.PushBack(cell)

		
			}else if k>0&&k<len(temp)-1&&cntI>0&&check==false&&str=="-"{//empty cells 
				amountList.PushBack(Cell{})
				
			}
		}
	
		cntI++

		
	}
	var amt [][]Cell//this is the board that holds all the cells 
	var tmpC [] Cell
	var tmpC2 [] Cell
	cnt=0
	for cell := amountList.Front(); cell != nil; cell = cell.Next(){//takes the list and turns it into a usable 2d array 
		tmpC=append(tmpC,cell.Value.(Cell))
		cnt++
		if cnt==jSize{
			amt=append(amt,tmpC)
			cnt=0
			tmpC = tmpC2
		}
		
	}
	


	bd:=Boards{cst,sup,dem,amt}//creates boards struct used for the alg 
	if bd.checkDegn(){
		panic("DEGEN CASE")
	}

	c := make(chan Path,10)//result chan for paths 
	bd.print()
	bd.steppingStone(c)//call the stepping stone algorithm
	wg.Wait()
	fmt.Println(pathCost)
	bd.print()



	letter:=[]string{"A","B","C","D","E","F","G","H","I","J","K"}//the rest of these lines are used just to write to the file 
	f, err := os.Create("solution.txt")
    if err != nil {
        fmt.Println(err)
    }
    str:="COST"
    for i:=0;i<jSize;i++{
    	str=str+" "+letter[i]
    }
    str=str+" SUPPLY"
    f.WriteString(str+"\n")

    temps:=""
    for i:=0;i<iSize;i++{
    	temps=temps+"Source"+strconv.Itoa(i+1)+" "
    	for j:=0;j<jSize;j++{
    		if bd.boardA[i][j].amount!=0{
    			temps=temps+ strconv.Itoa(bd.boardA[i][j].amount)+" "
    		}else{
    			temps=temps+"- "
    		}

    	
    	}
    	temps=temps+strconv.Itoa(bd.supply[i])
    	f.WriteString(temps+"\n")
    	temps=""
    }
    temps="DEMAND "
    for i:=0;i<len(bd.demand);i++{
    	temps=temps+strconv.Itoa(bd.demand[i])+" "
    }
   	f.WriteString(temps+"\n")
   	f.WriteString("TotalCost: "+strconv.Itoa(pathCost))




}



