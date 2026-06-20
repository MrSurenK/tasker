# Go-Tasker

A daily to do list written in GO to learn the syntax and gain faimilarity with the lang.

How it works:
User starts program in command line
|
User keys in taks for the day and all taks are saved in Documents/Tasks folder in markup files
|
As user completes tasks, user crosses out tasks
|
End of the day or when user uses program again old tasks are cleared completely

Potential Future Features:

1. TUI interface
2. Lightweight db

------ Task Flow --------
User starts app
|
Check for file and if no file for the day create new file
|
If file found then parse all the tasks into Task[] object for easier processing

CRUD functionality:

- All performed on Task[] object

Save

- Everytime a change is made to Task[] object it is saved to markdown file. But source of truth within application session remains the Task[] object.
