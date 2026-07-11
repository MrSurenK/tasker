package terminal

import (
	"os"

	"golang.org/x/term"
)

//set up terminal in raw mode
func SetUpTerminal() *term.State {
	// ---- User input to start application flow ------ // 
	//Put terminal in raw mode 
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	return oldState
}

//Return terminal back to cooked state from raw mode 
func Restore(oldState *term.State) error{
	return term.Restore(int(os.Stdin.Fd()), oldState)
}

//Show welcome msg
func ShowWelcome(){
	os.Stdout.Write([]byte("Welcome to TASKER!\r\n"))
}

func StartToDo() bool{
	//Give use prompt to proceed. Keep asking until a valid input is given via while loop
	for{
		os.Stdout.Write([]byte("\r\nShall we get productive today?...(y/n)\r\n"))

		key,read :=readKeyBoardInput()
		if !read {
			continue
		}
		//if user says no then exit program and log 
		switch key[0] {
			case 'y', 'Y':
			os.Stdout.Write([]byte("\r\nLet's go!\r\n"))
			return true //break out of while loop and set up environment
		case 'n', 'N':
			os.Stdout.Write([]byte("\r\nMaybe next time!\r\n"))
			return false; 
		case 3: // Ctrl+C
    		os.Stdout.Write([]byte("\r\nInterrupted.Closing Program!\r\n"))
			return false
		case 4: // Ctrl+D
    		os.Stdout.Write([]byte("\r\nEOF. Exiting.\r\n"))
    		return false
		default:
			os.Stdout.Write([]byte("\033[H\033[2J"))
   			os.Stdout.Write([]byte("\r\nWelcome to TASKER!\r\n"))
    		os.Stdout.Write([]byte("\r\nPlease key in a valid key \r\n"))
    		continue
		}
	}
}



func readKeyBoardInput() ([]byte, bool) {
		//read user response
		key := make([]byte, 8)
		if _, err := os.Stdin.Read(key); err != nil {
			//handle error
			os.Stdout.Write([]byte("\r\n Could not read input. Please try again"))
			return nil, false
		}
		return key, true
}

/*
Function to get user instruction on what to do next. Output wil be an int that will call the different api accordingly
1. Create a new item on list 
2. Edit an existing item on list
3. Delete an existing item on list
4. Show current items on list

if false returned then in main application just return and exit appliation
*/
func TellMeWhatToDo() (byte, bool){

	//Print msg in terminal to request for user input
	for{
		os.Stdout.Write([]byte("\r\n What would you like to do?\r\n1.Add a task\r\n2.Edit task\r\n3.Delete a task\r\n4.Show current items on list\r\n5.Exit application\r\n"))
	
		key,read :=readKeyBoardInput()
			if !read {
				continue
			}
		//Handle edge cases: like invalid inputs and return user input for future flexibility of logging 
		switch key[0]{
		case '5':
			os.Stdout.Write([]byte("\r\nExiting...\r\n"))
			return '5',true			
		case 3: // Ctrl+C
			os.Stdout.Write([]byte("\r\nInterrupted.Closing Program!\r\n"))
			return 3, false
		case 4: // Ctrl+D
			os.Stdout.Write([]byte("\r\nEOF. Exiting.\r\n"))
			return 4, false
		case '1','2','3','4':
			return key[0], true
		default: 
			os.Stdout.Write([]byte ("\r\nInvalid option\r\n"))
		}
	}	
}

//Helper method to handle  user input in terminal raw mode 
func ReadLine(prompt string)(string){
	os.Stdout.Write([]byte("\r\n" + prompt))
	var line []byte
	for{
		key := make([]byte, 8)
		n, _ := os.Stdin.Read(key)
		switch key[0]{
		case 13: //Enter
			if len(line) == 0 {
				continue // ignore stray Enter at the start
			}
			os.Stdout.Write([]byte("\r\n"))
			return string(line)
		case 127: //Backspace
			if len(line) > 0 {
				line = line[:len(line)-1]
				//erase the character: backspace, space, backspace beacause when we echo the character the cursor moves to the right so we need to delete the space first and then move the curor to the last character
				os.Stdout.Write([]byte("\b \b"))
			}
		case 3: //Ctrl + C - abort 
			return ""
		default: 
			line = append(line, key[:n]...)
			os.Stdout.Write(key[:n]) //echo the character so user can see what they typed in terminal
		}
	}
}

/*
Function to call the various crud operations to perform on the todo list file
Input: take todo list file
*/
func executeCommand(){

}


//Function to write tasks to markdown file
func SaveTask(){
	
}














