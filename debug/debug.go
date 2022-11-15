package debug

import (
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/divan/inspect/metrics"
	"github.com/divan/inspect/os/memstat"
)

var MainMemStat = NewMemStatWithMetricContext()
var logMemStat = logger.GetOrCreate("arwen/memstat")

func NewMemStatWithMetricContext() *memstat.MemStat {
	mc := metrics.NewMetricContext("system")
	stat := memstat.New(mc, time.Millisecond*3600)
	return stat
}

func PrintMemStat(stat *memstat.MemStat, label string) {
	stat.Collect()
	logMemStat.Trace("memstat",
		"label", label,
		"total", stat.MemTotal.Get(),
		"free", stat.MemFree.Get(),
		"anonActive", stat.Active_anon.Get(),
		"anonInactive", stat.Inactive_anon.Get(),
		"anonPages", stat.AnonPages.Get(),
		"mapped", stat.Mapped.Get(),
		"usage", stat.Usage(),
	)
}
