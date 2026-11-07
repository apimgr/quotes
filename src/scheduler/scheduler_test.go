package scheduler

import (
	"errors"
	"testing"
	"time"
)

func TestSchedulerLifecycle(t *testing.T) {
	s := New()
	calls := 0
	s.AddTask("task", time.Hour, func() error { calls++; return nil })
	if got := s.GetTasks(); len(got) != 1 || !got[0].Enabled {
		t.Fatalf("tasks = %#v", got)
	}
	if err := s.RunNow("missing"); err != nil {
		t.Fatalf("missing task: %v", err)
	}
	if err := s.RunNow("task"); err != nil || calls != 1 {
		t.Fatalf("RunNow: calls=%d err=%v", calls, err)
	}
	s.DisableTask("task")
	if s.GetTasks()[0].Enabled {
		t.Fatal("DisableTask left task enabled")
	}
	s.EnableTask("task")
	if !s.GetTasks()[0].Enabled {
		t.Fatal("EnableTask left task disabled")
	}
	s.RemoveTask("task")
	if got := s.GetTasks(); len(got) != 0 {
		t.Fatalf("RemoveTask left tasks: %#v", got)
	}
	s.Start()
	s.Start()
	s.Stop()
	s.Stop()
}

func TestSchedulerDueTaskAndErrors(t *testing.T) {
	s := New()
	completed := make(chan struct{}, 1)
	s.AddTask("ok", time.Second, func() error { completed <- struct{}{}; return nil })
	s.AddTask("bad", time.Second, func() error { return errors.New("expected") })
	s.AddTask("disabled", time.Second, func() error { t.Error("disabled task ran"); return nil })
	s.DisableTask("disabled")
	s.mu.Lock()
	for _, task := range s.tasks {
		task.NextRun = time.Now().Add(-time.Second)
	}
	s.mu.Unlock()
	s.runDueTasks()
	select {
	case <-completed:
	case <-time.After(time.Second):
		t.Fatal("due task did not run")
	}
}

func TestParseInterval(t *testing.T) {
	for key, want := range map[string]time.Duration{"minutely": time.Minute, "hourly": time.Hour, "daily": 24 * time.Hour, "weekly": 7 * 24 * time.Hour, "monthly": 30 * 24 * time.Hour, "2h": 2 * time.Hour, "bad": 24 * time.Hour} {
		if got := ParseInterval(key); got != want {
			t.Errorf("ParseInterval(%q) = %s, want %s", key, got, want)
		}
	}
}
