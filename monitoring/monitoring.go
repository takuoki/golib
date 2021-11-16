package monitoring

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
)

type contextKey string

const monitoringKey contextKey = "monitoring"

func monitoringFromContext(ctx context.Context) *monitoring {
	val := ctx.Value(monitoringKey)
	if m, ok := val.(*monitoring); ok {
		return m
	}
	return nil
}

type monitoring struct {
	mu       sync.Mutex
	names    []string
	memStats []runtime.MemStats
}

func (m *monitoring) record(name string) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	m.mu.Lock()
	defer m.mu.Unlock()

	m.names = append(m.names, name)
	m.memStats = append(m.memStats, mem)
}

// Record records the memory usage at the time of the call.
func Record(ctx context.Context, name string) context.Context {
	m := monitoringFromContext(ctx)
	if m == nil {
		m = &monitoring{}
		ctx = context.WithValue(ctx, monitoringKey, m)
	}

	m.record(name)

	return ctx
}

// Output outputs a record of memory usage.
// If `formatFunc` == nil, the output will be in the default format.
func Output(ctx context.Context, w io.Writer, formatFunc OutputFormatFunc) {
	m := monitoringFromContext(ctx)
	if m == nil {
		return
	}

	if formatFunc == nil {
		formatFunc = defaultOutputFormatFunc
	}

	formatFunc(w, m.names, m.memStats)
}

// OutputFile outputs a memory usage record to a file.
// The file is output to the specified directory with the file name `YYYYMMDDTHHmmss.md`.
// If `formatFunc` == nil, the output will be in the default format.
func OutputFile(ctx context.Context, dirPath string, formatFunc OutputFormatFunc) error {

	if dir, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)
	} else if !dir.IsDir() {
		return errors.New("the same file as the output directory name exists")
	}
	file, err := os.Create(fmt.Sprintf("%s/%s.md", dirPath, time.Now().Format("20060102T150405")))
	if err != nil {
		return fmt.Errorf("failed to create monitoring result file: %w", err)
	}

	Output(ctx, file, formatFunc)

	return nil
}

// OutputFormatFunc is a function that specifies the output format.
// The arguments `names` and` memStats` have the same number of elements.
type OutputFormatFunc func(w io.Writer, names []string, memStats []runtime.MemStats)

func defaultOutputFormatFunc(w io.Writer, names []string, memStats []runtime.MemStats) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "Alloc (MB)", "TotalAlloc (MB)", "Sys (MB)"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)

	for i, name := range names {
		table.Append([]string{
			name,
			toMb(memStats[i].Alloc),
			toMb(memStats[i].TotalAlloc),
			toMb(memStats[i].Sys),
		})
	}

	table.Render()
}

func toMb(b uint64) string {
	return strconv.FormatUint(b/1024/1024, 10)
}
