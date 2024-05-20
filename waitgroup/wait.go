package waitgroup

import (
	"runtime/debug"

	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wangweihong/gotoolbox/randutil"
)

type WaitGroupRoutineFunc struct {
	Name string
	Ctx  context.Context
	// 异步函数体
	Call func() Result
}

func NewWaitGroupHandleFunc(ctx context.Context, name string, call func() Result) WaitGroupRoutineFunc {
	if name == "" {
		name = "async-routine"
	}
	name = name + "-" + randutil.RandNumSets(6)
	return WaitGroupRoutineFunc{
		Name: name,
		Ctx:  ctx,
		Call: call,
	}
}

func NewWaitGroup(ctx context.Context) *Group {
	return &Group{
		ctx:     ctx,
		wg:      sync.WaitGroup{},
		results: make(map[string]Result),
		retLock: sync.Mutex{},
	}
}

// Group allows to start a group of goroutines and wait for their completion.
type Group struct {
	ctx         context.Context
	wg          sync.WaitGroup
	results     map[string]Result
	retLock     sync.Mutex
	debug       bool
	printReturn bool
}

func (g *Group) Debug() *Group {
	g.debug = true
	return g
}

func (g *Group) DebugReturn() *Group {
	g.printReturn = true
	return g
}

func (g *Group) Wait() {
	g.wg.Wait()
	g.PrintResults()
}

// Start starts f in a new goroutine in the group.
func (g *Group) Start(f WaitGroupRoutineFunc) {
	g.wg.Add(1)
	go func() {
		ret := NewResult(nil, nil)
		start := time.Now()
		defer g.wg.Done()
		defer g.setResult(f.Name, &ret, start)
		defer g.handleWaitGroupCrash(&ret)
		ret = f.Call()
	}()
}

func (g *Group) setResult(name string, ret *Result, startTime time.Time) {
	ret.Cost = time.Since(startTime)
	g.retLock.Lock()
	defer g.retLock.Unlock()
	g.results[name] = *ret
}

func (g *Group) GetResults() map[string]Result {
	g.retLock.Lock()
	defer g.retLock.Unlock()
	return g.results
}

func (g *Group) GetResultFast() ( /*total*/ int /*success*/, int /*fail*/, int) {
	g.retLock.Lock()
	defer g.retLock.Unlock()

	total := len(g.results)
	var fail int
	var success int
	for _, v := range g.results {
		if v.Error != nil {
			fail += 1
		} else {
			success += 1
		}
	}
	return total, success, fail
}

func (g *Group) PrintResults() {
	g.retLock.Lock()
	defer g.retLock.Unlock()
	if g.debug {
		fmt.Printf("has start %v waitgroup routines\n", len(g.results))
		for k, v := range g.results {
			if g.printReturn {
				fmt.Printf("waitgroup routine %v, results:%v, cost:%vs \n", k, v, v.Cost.Seconds())
			} else {
				fmt.Printf("waitgroup routine %v, cost:%vs \n", k, v.Cost.Seconds())
			}
		}
	}
}

func (g *Group) handleWaitGroupCrash(st *Result) {
	if x := recover(); x != nil {
		st.Error = fmt.Errorf("runtime panic:%v, stack:%v", x, string(debug.Stack()))
	}
}

func (g *Group) ConvertResultToBatchOutput() BatchOutput {
	g.retLock.Lock()
	defer g.retLock.Unlock()

	var bo BatchOutput
	for _, v := range g.results {
		bo.Total += 1
		if v.Error != nil {
			bo.Fail += 1
		} else {
			bo.Success += 1
		}
		bo.Results = append(bo.Results, SetOutput(v.Data, v.Error))
	}
	return bo
}

type Result struct {
	Cost  time.Duration
	Error error
	Data  interface{}
}

func NewResult(data interface{}, err error) Result {
	return Result{
		Data:  data,
		Error: err,
	}
}
