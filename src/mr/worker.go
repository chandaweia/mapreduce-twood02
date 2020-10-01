package mr

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"sort"
	"time"
)

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

type worker struct {
	id      int
	mapf    func(string, string) []KeyValue
	reducef func(string, []string) string
}

func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.
	w := worker{}
	w.mapf = mapf
	w.reducef = reducef

	w.RegWorker()
	// keep requesting tasks
	for {
		task := w.ReqTask()
		if task.inputFiles != nil {
			w.SeqMR(task)
			w.RepTask(task)
		} else {
			time.Sleep(time.Millisecond * 500)
		}
	}

}

// SeqMR runs the full sequential Map and Reduce phases for a list of input files
func (w *worker) SeqMR(task Task) {
	//
	// read each input file,
	// pass it to Map,
	// accumulate the intermediate Map output.
	//
	intermediate := []KeyValue{}
	for _, filename := range task.inputFiles {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("cannot open %v", filename)
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("cannot read %v", filename)
		}
		file.Close()
		DPrintf("Mapping file %s\n", filename)
		kva := w.mapf(filename, string(content))
		intermediate = append(intermediate, kva...)
	}

	//
	// a big difference from real MapReduce is that all the
	// intermediate data is in one place, intermediate[],
	// rather than being partitioned into NxM buckets.
	//

	sort.Sort(ByKey(intermediate))

	oname := "mr-out-0"
	ofile, _ := os.Create(oname)

	//
	// call Reduce on each distinct key in intermediate[],
	// and print the result to mr-out-0.
	//
	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := w.reducef(intermediate[i].Key, values)

		// this is the correct format for each line of Reduce output.
		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

		i = j
	}

	ofile.Close()
	DPrintf("Finished writing to %v\n", oname)
}

// RegWorker will register the worker with the Master.
// TODO: Get back an ID that is assigned to this worker
func (w *worker) RegWorker() {
	argsReg := RegisterWorkerArgs{}
	replyReg := RegisterWorkerReply{}

	// send the RPC request, wait for the reply.
	call("Master.RegisterWorker", &argsReg, &replyReg)
	DPrintf("Master reports %d workers\n", replyReg.NumWorkers)
	w.id = replyReg.NumWorkers

}

// ReqTask will ask the server for all files to be processed
// TODO: Should only get back one task (either a map or reduce task)
func (w *worker) ReqTask() Task {
	argsReq := RequestTaskArgs{}
	replyReq := RequestTaskReply{}
	task := Task{}
	if call("Master.RequestTask", &argsReq, &replyReq) == false {
		os.Exit(0)
	}
	DPrintf("Master wants me to process %v workers\n", replyReq.InputFiles)
	task.inputFiles = replyReq.InputFiles
	return task
}

func (w *worker) RepTask(t Task) {
	args := ReportTaskArgs{}
	reply := ReportTaskReply{}
	args.CompletedFiles = t.inputFiles
	args.WorkerID = w.id
	call("Master.ReportTask", &args, &reply)
	if reply.MoreTasks == false {
		DPrintf("No more tasks need to be processed")
	}
}

//
// example function to show how to make an RPC call to the master.
//
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	call("Master.Example", &args, &reply)

	// reply.Y should be 100.
	DPrintf("reply.Y %v\n", reply.Y)
}

//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	c, err := rpc.DialHTTP("unix", "mr-socket")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
