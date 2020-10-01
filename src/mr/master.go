package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

type Master struct {
	// Your definitions here.
	inputFiles []string
	numWorkers int
	isDone     bool
	sentTasks  bool
	lock       sync.Mutex
}

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
func (m *Master) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	DPrintf("Worker has called the Example RPC\n")
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	os.Remove("mr-socket")
	l, e := net.Listen("unix", "mr-socket")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {

	return m.isDone
}

//RegisterWorker is an RPC method that is called by workers after they have started
// up to report that they are ready to receive tasks.
func (m *Master) RegisterWorker(args *RegisterWorkerArgs, reply *RegisterWorkerReply) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.numWorkers++
	reply.NumWorkers = m.numWorkers
	DPrintf("Registered worker %d\n", m.numWorkers)

	return nil
}

//RequestTask is an RPC method that is called by workers to request a map or reduce task
func (m *Master) RequestTask(args *RequestTaskArgs, reply *RequestTaskReply) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	if m.sentTasks == false {
		reply.InputFiles = m.inputFiles
		DPrintf("Sending file list: %v\n", reply.InputFiles)
		m.sentTasks = true
	} else {
		reply.InputFiles = nil
		DPrintf("Already sent files to be processed\n")
	}

	return nil
}

//ReportTask is an RPC method that is called by workers to report a task's status
//whenever a task is finished or failed
//HINT: when a task is failed, master should reschedule it.
func (m *Master) ReportTask(args *ReportTaskArgs, reply *ReportTaskReply) error {
	DPrintf("Worker %d finished files %v", args.WorkerID, args.CompletedFiles)
	// TODO: Check if there are more tasks to be processed!
	reply.MoreTasks = false
	m.isDone = true
	return nil
}

//
// create a Master.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	m.inputFiles = files
	m.numWorkers = 0
	m.isDone = false
	m.sentTasks = false

	// Your code here.

	go m.server()

	return &m
}
