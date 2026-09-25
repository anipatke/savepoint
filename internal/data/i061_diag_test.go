package data

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"testing"
	"time"
)

// Throwaway diagnostic for I-061; never merged. It always fails on Windows so
// its measurements print in the CI failure log.
func TestI061Diagnostic(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only diagnostic")
	}
	var report []string
	add := func(format string, args ...any) { report = append(report, fmt.Sprintf(format, args...)) }

	// 1. Do successive create/remove cycles of the lock path get distinct mtimes?
	dir := t.TempDir()
	path := filepath.Join(dir, taskIDLockFile)
	for _, gap := range []time.Duration{0, 2 * time.Millisecond, 10 * time.Millisecond} {
		distinct := map[time.Time]bool{}
		sameAsPrev := 0
		var prev time.Time
		for i := 0; i < 40; i++ {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
			if err != nil {
				add("mtime gap=%s create err: %v", gap, err)
				break
			}
			info, err := os.Lstat(path)
			if err == nil {
				if info.ModTime().Equal(prev) {
					sameAsPrev++
				}
				prev = info.ModTime()
				distinct[info.ModTime()] = true
			}
			f.Close()
			os.Remove(path)
			time.Sleep(gap)
		}
		add("mtime gap=%s: 40 cycles, %d distinct mtimes, %d equal to previous holder", gap, len(distinct), sameAsPrev)
	}

	// 2. How long does one uncontended allocation and one strict load take?
	root := newTaskIDProject(t)
	var loads, allocs []time.Duration
	for i := 0; i < 16; i++ {
		start := time.Now()
		if _, err := LoadV2Index(root); err != nil {
			add("load err: %v", err)
		}
		loads = append(loads, time.Since(start))
		start = time.Now()
		if _, err := AllocateTaskID(root); err != nil {
			add("alloc err: %v", err)
		}
		allocs = append(allocs, time.Since(start))
	}
	add("sequential LoadV2Index: %s", summarize(loads))
	add("sequential AllocateTaskID (preflight load + held section): %s", summarize(allocs))

	// 3. Concurrent callers: per-call latency and which ones fail.
	root = newTaskIDProject(t)
	var mu sync.Mutex
	var ok, failed []time.Duration
	var group sync.WaitGroup
	for range 16 {
		group.Add(1)
		go func() {
			defer group.Done()
			start := time.Now()
			_, err := AllocateTaskID(root)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				failed = append(failed, time.Since(start))
			} else {
				ok = append(ok, time.Since(start))
			}
		}()
	}
	group.Wait()
	add("concurrent 16: ok=%d (%s) failed=%d (%s)", len(ok), summarize(ok), len(failed), summarize(failed))

	for _, line := range report {
		t.Log(line)
	}
	t.Errorf("I-061 diagnostic report above (intentional failure)")
}

func summarize(values []time.Duration) string {
	if len(values) == 0 {
		return "none"
	}
	sorted := append([]time.Duration(nil), values...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	var total time.Duration
	for _, v := range sorted {
		total += v
	}
	return fmt.Sprintf("min=%s median=%s max=%s mean=%s", sorted[0].Round(time.Microsecond), sorted[len(sorted)/2].Round(time.Microsecond), sorted[len(sorted)-1].Round(time.Microsecond), (total / time.Duration(len(sorted))).Round(time.Microsecond))
}
