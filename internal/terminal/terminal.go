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
Function to get user instruction on what to do next? 
1. Create a new item on list 
2. Edit an existing item on list
3. Delete an existing item on list
4. Show current items on list
*/
// func tellMeWhatToDo()













