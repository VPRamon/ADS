@ECHO OFF
SETLOCAL ENABLEDELAYEDEXPANSION
FOR /L %%x IN (1,1,2) DO (
    SET "NUM=configFile_600%%x.txt
	start cmd.exe /k "go run election.go 2 2_echo_room/!NUM!"
)