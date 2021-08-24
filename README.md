# ADS
## Authors:
	Ramon Vallés Puig, Nia:205419
	Marcelo Sánchez Ortega, Nia:204592

We have created 2 .sh scripts to compile it in a easy way: 
To execute the scripts we must specify as arguments the size of the matrices "sh run.sh M N K"
The structure of the scripts is the following:

    Create a file with 2 matrices of size MxN and NxK (given parameters):   Create matricespython matrixGen.py $1 $2 $3
    
    Compile the java file:   hadoop com.sun.tools.javac.Main MatrixMultiplier.java
    
    Pack in a jar file:   jar cf mm.jar MatrixMultiplier*.class
    
    Delete obsolete directories:    hadoop fs -rm -f -r /user/u87515/output/mm_partial/
    
                                    hadoop fs -rm -f -r /user/u87515/output/mm/
				    
    Execute jar and compute time:   time hadoop jar mm.jar MatrixMultiplier /user/u149889/input/mm /user/u149889/output/mm
    
    Show resoulting Matrix:         #hadoop fs -cat /user/u87515/output/mm/part-r-00000
