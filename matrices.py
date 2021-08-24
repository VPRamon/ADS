import random
import sys

row = int(sys.argv[1])
col = int(sys.argv[2])
num = sys.argv[3]
filename = str(num) + ".txt"
f = open(filename,"w+") 

for i in range(col):
	for j in range(row):
		val = random.randint(0,500)
		f.write(str(i) + " "+  str(j) + " " + str(val) + "\n" ) 

f.close() 


