package POOL

import (
	"fmt"
	"sync"
	"time"
)

// Version 1.0.4

var TimeOutFinish = 1 * time.Second
var TimeOutError = 5 * time.Second

// NewPool : Function used to init a pool
// t  is an any of object
// n  is a Name of pool
// f is a function call foreach action to do OutType
// nbParallel is amount of simultaneous worker
// debug is a Boolean to show debug
// i is the duration of each run
func NewPool(t any, n string, f OutType, nbParallel int, debug bool, i time.Duration) (*Pool, error) {
	p := Pool{}
	p.Type = t
	p.Name = n
	p.Out = f
	p.MaxWorker = nbParallel
	p.Debug = debug
	p.Interval = i
	return &p, nil
}

// Pool Structure to build a pool
type Pool struct {
	Name      string `json:"name"`
	Type      any    `json:"type,omitempty"`
	Status    *PoolStatus
	content   sync.Map
	MaxWorker int
	Debug     bool
	Interval  time.Duration
	NbWorker  int
	Out       OutType
	Worker    map[int]bool
	// Trigger on each status
	TriggerError   OutType
	TriggerSuccess OutType
	TriggerFinish  OutType
}
type OutType func(p *PoolEntry) error

type PoolEntry struct {
	Status  string
	Content any

	Priority int
	Key      string
	Result   error
	Create   time.Time
	Update   time.Time
}

// Count Used to get current pool quantity
func (p *Pool) Count() int {
	var i int
	p.content.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}
func (p *Pool) Clean() {
	for {
		p.content.Range(func(k, v interface{}) bool {
			value, ok := v.(PoolEntry)
			if ok == true {
				if value.Status == "Finish" {
					if time.Now().Sub(value.Update) > TimeOutFinish {
						p.Context("Remove entry ", k)
						fmt.Println("Remove entry finish", k)
						p.content.Delete(k)
					}
				}
				if value.Status == "Error" {
					if time.Now().Sub(value.Update) > TimeOutError {
						p.Context("Remove entry ", k)
						fmt.Println("Remove entry in error", k)
						p.content.Delete(k)
					}
				}
			}
			return true
		})
		time.Sleep(10 * time.Millisecond)
	}
}

// CountStatus Used to get current pool quantity
func (p *Pool) CountStatus() map[string]int {
	var i = make(map[string]int)
	p.content.Range(func(k, v interface{}) bool {
		value, ok := v.(PoolEntry)
		if ok == true {
			i[value.Status]++
		}

		return true
	})
	return i
}

// Init (
func (p *Pool) Init() {
	p.Context("Init Pool", p.Name)
	p.Status = &PoolStatusInit
	p.Worker = make(map[int]bool)
	go p.Run()
	go p.Clean()
}

// Add Used to add a new entry of type T
func (p *Pool) Add(k string, t any, pr int) error {
	p.Context("Try to add new entry ", k)
	obj := PoolEntry{
		Status:   "New",
		Content:  t,
		Priority: pr,
		Key:      k,
		Create:   time.Now(),
		Update:   time.Now(),
	}
	p.content.Store(k, obj)
	return nil
}
func (p *Pool) UpdateStatus(k string, s string) error {
	p.Context("Try to add new entry ", k)
	obj, err := p.content.Load(k)
	if err != true {
		return ErrorPoolKeyNotExist
	}
	value, ok := obj.(PoolEntry)
	if ok != true {
		return ErrorPoolEntryMismatch
	}
	value.Status = s
	value.Update = time.Now()
	p.content.Swap(k, value)
	if s == "Error" && p.TriggerError != nil {
		p.Context("Entry in error and a trigger exist ", k)
		errError := p.TriggerError(&value)
		if errError == nil {
			p.Context("Trigger return no error clean ", k)
			p.content.Delete(k)
		}
	}
	if s == "Finish" && p.TriggerFinish != nil {
		p.Context("Entry in finish and a trigger exist ", k)
		errSuccess := p.TriggerFinish(&value)
		if errSuccess == nil {
			p.Context("Trigger return no error clean ", k)
			p.content.Delete(k)
		}
	}
	return nil
}

// Context Used to log or show error
func (p *Pool) Context(msg ...any) {
	if p.Debug == true {
		fmt.Println(time.Now().Format(time.DateTime), "POOL-"+p.Name, ":", msg)
	}

}
func count(p *Pool) {
	for {
		p.Context("Nb Worker ", p.NbWorker, "/", p.MaxWorker)
		time.Sleep(1 * time.Second)
	}
}

func addWorker(o *PoolEntry, p *Pool) {
	//fmt.Println(o)
	p.NbWorker++
	var i = p.NbWorker
	p.Context("Start Worker", i)
	err := p.Out(o)
	if err != nil {
		p.UpdateStatus(o.Key, "Error")
	} else {
		p.UpdateStatus(o.Key, "Finish")
	}

	p.Context("Finish Worker", i)
	p.NbWorker--
}

// Run  Used to launch worker
func (p *Pool) Run() {
	go count(p)
	p.Context("Run survey of pool")
	for {
		nb := p.Count()
		if p.MaxWorker > p.NbWorker {
			todoObj, key := p.GetNext()
			if todoObj != nil {
				p.Context("Launching a new object ", key, todoObj)
				p.UpdateStatus(key, "Lock")
				go addWorker(todoObj, p)
			}

		} else {
			p.Context("All worker are used", p.MaxWorker, p.NbWorker)
		}

		p.Context("Waiting next interval with ", nb, "entries to launch")
		p.Context("Mapping Status", p.CountStatus())
		time.Sleep(10 * time.Millisecond)
	}
}

// Remove Used to remove an entry
func (p *Pool) Remove(k string) {
	p.Context("Remove entry of pool", k)
	p.content.Delete(k)
}

// Unload Used to safe stop
func (p *Pool) Unload() {
	// Todo Save Content of pool to temporary
	p.Context("Stop survey of pool")
}

// TODO Add priority on GetNext
func (p *Pool) GetNext() (b *PoolEntry, key string) {
	p.content.Range(func(k, v any) bool {
		value, ok := v.(PoolEntry)
		if ok && value.Status == "New" {
			b = &value
			key = fmt.Sprint(k)

		}
		return true
	})
	p.Context("Found key ", key)
	return b, key
}
