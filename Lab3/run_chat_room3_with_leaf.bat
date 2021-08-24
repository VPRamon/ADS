@ECHO OFF
SETLOCAL ENABLEDELAYEDEXPANSION
FOR /L %%x IN (1,1,3) DO (
    SET "NUM=configFile_600%%x.txt
	start cmd.exe /k "go run election.go 3_echo_room_with_leaf/!NUM!"
)