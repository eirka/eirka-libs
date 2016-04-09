package status

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"

	e "github.com/eirka/eirka-libs/errors"
)

var (
	startTime = time.Now()
)

// Statistics holds runtime stats
type Statistics struct {
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

// StatusController is a Gin controller to display current runtime info over http
func StatusController(c *gin.Context) {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	stats := &Statistics{
		Uptime:       humanize.Time(startTime),
		NumGoroutine: runtime.NumGoroutine(),
		MemAllocated: humanize.Bytes(m.Alloc),
		MemTotal:     humanize.Bytes(m.TotalAlloc),
		MemSys:       humanize.Bytes(m.Sys),
		Lookups:      m.Lookups,
		MemMallocs:   m.Mallocs,
		MemFrees:     m.Frees,
		HeapAlloc:    humanize.Bytes(m.HeapAlloc),
		HeapSys:      humanize.Bytes(m.HeapSys),
		HeapIdle:     humanize.Bytes(m.HeapIdle),
		HeapInuse:    humanize.Bytes(m.HeapInuse),
		HeapReleased: humanize.Bytes(m.HeapReleased),
		HeapObjects:  m.HeapObjects,
		StackInuse:   humanize.Bytes(m.StackInuse),
		StackSys:     humanize.Bytes(m.StackSys),
		MSpanInuse:   humanize.Bytes(m.MSpanInuse),
		MSpanSys:     humanize.Bytes(m.MSpanSys),
		MCacheInuse:  humanize.Bytes(m.MCacheInuse),
		MCacheSys:    humanize.Bytes(m.MCacheSys),
		BuckHashSys:  humanize.Bytes(m.BuckHashSys),
		GCSys:        humanize.Bytes(m.GCSys),
		OtherSys:     humanize.Bytes(m.OtherSys),
		NextGC:       humanize.Bytes(m.NextGC),
		LastGC:       fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000),
		PauseTotalNs: fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000),
		PauseNs:      fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000),
		NumGC:        m.NumGC,
	}

	// Marshal the structs into JSON
	output, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		c.JSON(e.ErrorMessage(e.ErrInternalError))
		c.Error(err).SetMeta("StatusController.Marshal")
		return
	}

	c.Data(200, "application/json", output)

	return

}
