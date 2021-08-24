import java.io.IOException;
import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.DoubleWritable;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;
import java.util.HashMap;

import org.apache.hadoop.mapreduce.lib.input.TextInputFormat;
import org.apache.hadoop.mapreduce.lib.output.TextOutputFormat;



public class MatMult {
	public static class MatrixMapper extends org.apache.hadoop.mapreduce.Mapper<LongWritable, Text, Text, Text> { 		
	
		public void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {
			
			int M = Integer.parseInt(context.getConfiguration().get("M"));
			int K = Integer.parseInt(context.getConfiguration().get("K"));
			
			String line = value.toString(); 
			String[] element = line.split(" ");	//element [matrix Name, row, column, value]
			
            
            String matrixName = element[0];
            String row = element[1];
            String column = element[2];
            String val = element[3];
            
            Text outputKey = new Text();
            Text outputValue = new Text();   
            
            /* If element is from matrix A */
            if (matrixName.equals("A")) {        
            	
               	/* Repeat for each column of matrix B (K columns) */
                for (int k = 0; k < K; k++) {                	
                	outputKey.set(row + "," + k);	// key = {row,col} from matrix C
                    outputValue.set(matrixName + "," + column + "," + val);	 // Value = {Matrix, col, value} to be reduced
                    context.write(outputKey, outputValue);
                }
            }
            
            /* If element is from matrix B */
            else if (matrixName.equals("B")) {
            	            
            	/* Repeat for each row of matrix A (M rows) */
                for (int i = 0; i < M; i++) {
                    outputKey.set(i + "," + column);	// key = {row,col} from matrix C
                    outputValue.set(matrixName + "," + row + "," + val);	// Value = {Matrix, row, value} to be reduced
                    context.write(outputKey, outputValue);
                }
            }
		}
		
			
	}

	public static class MatrixReducer extends Reducer<Text, Text, Text, Text>
    {
		@Override
		public void reduce(Text key, Iterable<Text> values, Context context)
                throws IOException, InterruptedException {
	        String[] value;
	        //key=(i,k),
	        //Values = [(M/N,j,V/W),..]
	        HashMap<Integer, Float> hashA = new HashMap<Integer, Float>();
	        HashMap<Integer, Float> hashB = new HashMap<Integer, Float>();
	        for (Text val : values) {
	                value = val.toString().split(",");
	                if (value[0].equals("A")) {
	                        hashA.put(Integer.parseInt(value[1]), Float.parseFloat(value[2]));
	                } else {
	                        hashB.put(Integer.parseInt(value[1]), Float.parseFloat(value[2]));
	                }
	        }
	        int N = Integer.parseInt(context.getConfiguration().get("N"));
	        float result = 0.0f;
	        float m_ij;
	        float n_jk;
	        for (int j = 0; j < N; j++) {
	                m_ij = hashA.containsKey(j) ? hashA.get(j) : 0.0f;
	                n_jk = hashB.containsKey(j) ? hashB.get(j) : 0.0f;
	                result += m_ij * n_jk;
	        }
	        if (result != 0.0f) {
	                context.write(null, new Text(key.toString() + "," + Float.toString(result)));
	        }
		}
	
    }


	  public static void main(String[] args) throws Exception {
		  if (args.length != 5) {
	          System.err.println("Usage: MatrixMultiply <in_dir> <out_dir>");
	          System.exit(2);
	      }
		  
		  String M = args[2];
		  String N = args[3];
		  String K = args[4];
		  
	      Configuration conf = new Configuration();
	      // M is an m-by-n matrix; N is an n-by-p matrix.
	      conf.set("M", M);
	      conf.set("N", N);
	      conf.set("K", K);
	      @SuppressWarnings("deprecation")
	              Job job = new Job(conf, "MatrixMultiply");
	      job.setJarByClass(MatMult.class);
	      job.setOutputKeyClass(Text.class);
	      job.setOutputValueClass(Text.class);
	
	      job.setMapperClass(MatrixMapper.class);
	      job.setReducerClass(MatrixReducer.class);
	
	      job.setInputFormatClass(TextInputFormat.class);
	      job.setOutputFormatClass(TextOutputFormat.class);
	
	      FileInputFormat.addInputPath(job, new Path(args[0]));
	      FileOutputFormat.setOutputPath(job, new Path(args[1]));
	
	      job.waitForCompletion(true);
	  }
}