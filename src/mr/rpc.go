package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import "os"
import "strconv"

const (
	mapTask = iota
	reduceTask
	waitTask
	exitTask
)

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

// Add your RPC definitions here.

type RequestTaskArgs struct{}

type RequestTaskReply struct {
	TaskType     int
	InputFiles   []string
	NReduce      int
	Index        int
	CountOfSends int
}

type SubmitTaskArgs struct {
	TaskType    int
	OutputFiles []string
	Index       int
}

type SubmitTaskReply struct{}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
