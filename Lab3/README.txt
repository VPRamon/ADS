Authors:
	Ramon Vallés Puig, Nia:205419
	Marcelo Ortega, Nia:204592


The compilation of the program is: (you can also skip this step and directly run with "go run")
	- go build echo.go

To run the executable: you need to specify each configFile, to do so type "echo.exe configFile_600(i).txt
*Note: In case that the Config File is not in the same directory as the executable, when passing the config file, make sure to put the full path where it is stored.

To make it simpler we made a .bat file to run every node in a faster way. Each config file is stored in the folder '4_echo_room'

Once the program is running it will immediately read the configFile and start the connection
	
If the node is an initiator it will send a message to all the neighbours specified in the configFile.
Otherwise, if it is not an initiator it will wait until the first message is received, it will make the requesting node its parent and it will proceed with the echo.