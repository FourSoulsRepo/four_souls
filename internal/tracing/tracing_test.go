package tracing

import (
	"context"
	"flag"
	"testing"
)

func TestFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	o := Flags(fs)
	if err := fs.Parse([]string{"-trace_dir", "traces", "-pprof_addr", "127.0.0.1:6060"}); err != nil {
		t.Fatal(err)
	}
	if o.Dir != "traces" || o.Addr != "127.0.0.1:6060" {
		t.Errorf("flags = %+v", *o)
	}
}

func TestFromEnv(t *testing.T) {
	t.Setenv(EnvDir, "/var/traces")
	t.Setenv(EnvAddr, "127.0.0.1:6061")
	if o := FromEnv(); o.Dir != "/var/traces" || o.Addr != "127.0.0.1:6061" {
		t.Errorf("FromEnv = %+v", o)
	}
}

func TestStartZeroRecordsNothing(t *testing.T) {
	stop, err := Start(context.Background(), Options{})
	if err != nil {
		t.Fatalf("zero options: %v", err)
	}
	if err := stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
