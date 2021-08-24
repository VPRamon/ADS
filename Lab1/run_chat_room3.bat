@ECHO OFF
SETLOCAL ENABLEDELAYEDEXPANSION
FOR /L %%x IN (1,1,3) DO (
    SET "NUM=configFile_600%%x.txt
	start cmd.exe /k "go run program_ex2.go 3_chat_room/!NUM!"
)