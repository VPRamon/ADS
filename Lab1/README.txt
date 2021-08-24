Author: Ramon Vallés Puig, Nia:205419

The compilation of both programs are the same: (you can also skip this step and directly run with "go run")
	- go build program_ex(i).go

To run the executable:
	- Ex1 : just type "program_ex1.exe
	- Ex2 : you need to specify each configFile, to do so type "program_ex2.exe configFile_600(i).txt
		*Note: In case that the Config File is not in the same directory as the executable, when passing the config file, make sure to put the full path where it is stored.

To make it simpler I made a .bat file to run every node in a faster way (only for exercise 2).
I prepered 2 options for windows: run_chat_room(3 or 4).bat, the difference between them is the number of processes you execute. 3 or 4 are their respective number of executables.
You will find their ConfigFiles in their respective folders.

Once the program is running:
	- Ex1: It will ask for some information such as listening port and target adress
	- Ex2: It will immediately read the configFile and start the connection
	
In both cases, the program will wait until the connections are properly established. After that, every node can send any message just by typing
"send:" + the text they want to send. Although this is not strictly necessary, I think it makes the interactive part more intuitive. 
To finish the program you can just type stop, and after a countdown of 3 seconds, the program will exit.