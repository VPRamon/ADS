python matrixGen.py $1 $2 $3
hadoop com.sun.tools.javac.Main MatrixMultiplier.java
jar cf mm.jar MatrixMultiplier*.class
hadoop fs -rm -f -r /user/u87515/output/mm_partial/
hadoop fs -rm -f -r /user/u87515/output/mm/
time hadoop jar mm.jar MatrixMultiplier /user/u149889/input/mm /user/u149889/output/mm
#hadoop fs -cat /user/u87515/output/mm/part-r-00000

