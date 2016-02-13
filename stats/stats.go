package utils

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"runtime"
	"strings"
	"time"
)

var (
	startTime = time.Now()
)

var sysStatus struct {
	Uptime       string
	NumGoroutine int

	// General statistics.
	MemAllocated string // bytes allocated and still in use
	MemTotal     string // bytes allocated (even if freed)
	MemSys       string // bytes obtained from system (sum of XxxSys below)
	Lookups      uint64 // number of pointer lookups
	MemMallocs   uint64 // number of mallocs
	MemFrees     uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    string // bytes allocated and still in use
	HeapSys      string // bytes obtained from system
	HeapIdle     string // bytes in idle spans
	HeapInuse    string // bytes in non-idle span
	HeapReleased string // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	//	Inuse is bytes used now.
	//	Sys is bytes obtained from system.
	StackInuse  string // bootstrap stacks
	StackSys    string
	MSpanInuse  string // mspan structures
	MSpanSys    string
	MCacheInuse string // mcache structures
	MCacheSys   string
	BuckHashSys string // profiling bucket hash table
	GCSys       string // GC metadata
	OtherSys    string // other system allocations

	// Garbage collector statistics.
	NextGC       string // next run in HeapAlloc time (bytes)
	LastGC       string // last run in absolute time (ns)
	PauseTotalNs string
	PauseNs      string // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
	NumGC        uint32
}

func updateSystemStatus() {
	sysStatus.Uptime = humanize.Time(startTime)

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	sysStatus.NumGoroutine = runtime.NumGoroutine()

	sysStatus.MemAllocated = humanize.Bytes(m.Alloc)
	sysStatus.MemTotal = humanize.Bytes(m.TotalAlloc)
	sysStatus.MemSys = humanize.Bytes(m.Sys)
	sysStatus.Lookups = m.Lookups
	sysStatus.MemMallocs = m.Mallocs
	sysStatus.MemFrees = m.Frees

	sysStatus.HeapAlloc = humanize.Bytes(m.HeapAlloc)
	sysStatus.HeapSys = humanize.Bytes(m.HeapSys)
	sysStatus.HeapIdle = humanize.Bytes(m.HeapIdle)
	sysStatus.HeapInuse = humanize.Bytes(m.HeapInuse)
	sysStatus.HeapReleased = humanize.Bytes(m.HeapReleased)
	sysStatus.HeapObjects = m.HeapObjects

	sysStatus.StackInuse = humanize.Bytes(m.StackInuse)
	sysStatus.StackSys = humanize.Bytes(m.StackSys)
	sysStatus.MSpanInuse = humanize.Bytes(m.MSpanInuse)
	sysStatus.MSpanSys = humanize.Bytes(m.MSpanSys)
	sysStatus.MCacheInuse = humanize.Bytes(m.MCacheInuse)
	sysStatus.MCacheSys = humanize.Bytes(m.MCacheSys)
	sysStatus.BuckHashSys = humanize.Bytes(m.BuckHashSys)
	sysStatus.GCSys = humanize.Bytes(m.GCSys)
	sysStatus.OtherSys = humanize.Bytes(m.OtherSys)

	sysStatus.NextGC = humanize.Bytes(m.NextGC)
	sysStatus.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	sysStatus.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
	sysStatus.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	sysStatus.NumGC = m.NumGC
}
