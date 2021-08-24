import random
import sys

M = int(sys.argv[1])
N = int(sys.argv[2])
K = sys.argv[3]
filename = "matrices.txt"
f = open(filename,"w+")
f.write("size  " + str(M) +" "+ str(N) +" "+ str(K) + "\n")

for j in range(int(M)):
	for i in range(int(N)):
		val = random.randint(0,5)
		f.write("A " + str(j) + " "+  str(i) + " " + str(val) + "\n" ) 

for j in range(int(N)):
	for i in range(int(K)):
		val = random.randint(0,5)
		f.write("B " + str(j) + " "+  str(i) + " " + str(val) + "\n" ) 

f.close() 
