package threadpool
import
(
	"sync/atomic"
	"time"
)

type ThreadPool struct {
	thread int //线程数
	maxqueue int //最大队列
	nreq int32 //一共执行了多少任务
	jobs chan func()interface{}  //任务集
	results chan interface{} //返回集
}

func (this *ThreadPool)worker(id int,){
	for{
		j,ok:=<-this.jobs
		if ok==false{
			break
		}
		this.results <- j()
	}
}
func (this *ThreadPool)Reset(){
	atomic.StoreInt32(&this.nreq,0)
	this.jobs=make(chan func()interface{},this.maxqueue)
	this.results=make(chan interface{},this.maxqueue)
}
func (this *ThreadPool)Req(_f func()interface{}){
	atomic.AddInt32(&this.nreq,1)
	this.jobs<-_f
}
func (this *ThreadPool)Close(){
	close(this.jobs)
	close(this.results)
}
func (this *ThreadPool)Wait(){
	for{
		if atomic.LoadInt32(&this.nreq)==int32(len(this.results)){
			break
		}
		time.Sleep(time.Millisecond*10)
	}
}
func NewThreadPool(_thread int,_maxqueue int)*ThreadPool{
	pl:=&ThreadPool{thread:_thread}
	pl.jobs=make(chan func()interface{},_maxqueue)
	pl.results=make(chan interface{},_maxqueue)
	for i:=0;i<_thread;i++{
		go pl.worker(i,)
	}
	return pl
}
