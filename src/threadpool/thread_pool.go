package threadpool

import(
	"sync"
)

const(
	THREAD_COMPLETE int = iota
	THREAD_CHILD_CONTIUNE
	THREAD_MAIN_CONTIUNE
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
		defer this.wg.Done()

		task.Run() 
		this.addFinshTask(task)
	}()
 }
 func(this *ThreadPool) addFinshTask(task IThreadTask){

	this.mtx.Lock()
	defer this.mtx.Unlock()

	this.finshTask = append(this.finshTask,task)
 }
 func(this *ThreadPool) Update(){

	if this.runStatus == false{
		return
	}

	if len(this.finshTask) == 0 {
		return
	}

	this.mtx.Lock()

	tempTask := make([]IThreadTask,len(this.finshTask))
	copy(tempTask,this.finshTask)
	this.finshTask = this.finshTask[0:0]

	this.mtx.Unlock()

	for _, task := range tempTask{
		ret := task.PsentMainThread()
		if THREAD_COMPLETE == ret{
			
		} else if THREAD_CHILD_CONTIUNE == ret{
			this.AddTask(task)
		} else if THREAD_MAIN_CONTIUNE == ret{
			this.addFinshTask(task)
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