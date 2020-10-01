package mr

//
// RPC definitions.
//

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

type RegisterWorkerArgs struct {
}

// RegisterWorkerReply tells the client how many workers are currently online
type RegisterWorkerReply struct {
	NumWorkers int
}

type RequestTaskArgs struct {
}

// RequestTaskReply returns a list of files to be processed
// TODO: only send one small task to do instead of all files
type RequestTaskReply struct {
	InputFiles []string
}

type ReportTaskArgs struct {
	CompletedFiles []string
	WorkerID       int
}

type ReportTaskReply struct {
	MoreTasks bool
}

// Add your RPC definitions here.
