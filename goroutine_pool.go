package executor
 
import (
    log "github.com/sirupsen/logrus"
)
 
type Task interface {
    Execute()
}
 
/* 有关协程池的定义及操作 */
//定义池类型
type Pool struct {
    //对外接收Task的入口
    EntryChannel chan Task
 
    //协程池最大worker数量,限定Goroutine的个数
    worker_num int
 
    //协程池内部的任务就绪队列
    JobsChannel chan Task
}
 
//创建一个协程池
func NewPool(cap int) *Pool {
    p := Pool{
        EntryChannel: make(chan Task),
        worker_num:   cap,
        JobsChannel:  make(chan Task),
    }
 
    return &p
}
 
//协程池创建一个worker并且开始工作
func (p *Pool) worker(work_ID int) {
    //worker不断的从JobsChannel内部任务队列中拿任务
    for task := range p.JobsChannel {
        //如果拿到任务,则执行task任务
        task.Execute()
        log.Logger.Infof("worker ID %d execute finished", work_ID)
    }
}
 
//让协程池Pool开始工作
func (p *Pool) Run() {
    //1,首先根据协程池的worker数量限定,开启固定数量的Worker,
    //  每一个Worker用一个Goroutine承载
    for i := 0; i < p.worker_num; i++ {
        go p.worker(i)
    }
 
    //2, 从EntryChannel协程池入口取外界传递过来的任务
    //   并且将任务送进JobsChannel中
    for task := range p.EntryChannel {
        p.JobsChannel <- task
    }
 
    //3, 执行完毕需要关闭JobsChannel
    close(p.JobsChannel)
 
    //4, 执行完毕需要关闭EntryChannel
    close(p.EntryChannel)
}
