package threadpool

import(
	"sync"
)

const(
	THREAD_COMPLETE = iota
)

 type ITask interface{
	Run() bool
 }

 type IThreadTask interface{
	ITask
	PsentMainThread() int
 }

 type ThreadPool struct{
	finshTask	[]IThreadTask
	runStatus	bool
	mtx 		sync.Mutex
	wg			sync.WaitGroup
 }

 func New() *ThreadPool{
	return &ThreadPool{
		finshTask : make([]IThreadTask, 0),
		runStatus : true,
	}
 }

 func(this *ThreadPool) AddTask(task IThreadTask){

	if this.runStatus == false{
		return
	}

	go func(){

		this.wg.Add(1)

		task.Run()

		this.mtx.Lock()

		defer func(){
			this.mtx.Unlock()
			this.wg.Done()
		}() 

		this.finshTask = append(this.finshTask,task)

	}()
 }
 func(this *ThreadPool) Update(){

	if this.runStatus == false{
		return
	}

	if len(this.finshTask) == 0 {
		return
	}

	this.mtx.Lock()

	temp := make([]IThreadTask,0)
	temp = this.finshTask
	this.finshTask = this.finshTask[0:0]

	this.mtx.Unlock()

	for _, task := range temp{
		ret := task.PsentMainThread()
		if THREAD_COMPLETE == ret{
			this.AddTask(task)
		}	
	}
	
 }

 func(this* ThreadPool) Stop(){

	if this.runStatus == false{
		return
	}
	
	this.runStatus = false
	this.wg.Wait()
	this.finshTask = this.finshTask[0:0]
 }