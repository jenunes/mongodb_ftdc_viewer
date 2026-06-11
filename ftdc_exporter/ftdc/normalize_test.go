package ftdc

import (
	"testing"

	"github.com/evergreen-ci/birch"
)

func TestNormalizeMetricsDocument_flatAndNestedKeys(t *testing.T) {
	doc := birch.NewDocument(
		birch.EC.Int64("start", 1700000000000),
		birch.EC.Int64("serverStatus.wiredTiger.cache.bytes currently in the cache", 4096),
		birch.EC.SubDocument("serverStatus", birch.NewDocument(
			birch.EC.SubDocument("wiredTiger", birch.NewDocument(
				birch.EC.SubDocument("concurrentTransactions", birch.NewDocument(
					birch.EC.SubDocument("monitor", birch.NewDocument(
						birch.EC.Int32("available", 100),
						birch.EC.Int32("totalTickets", 128),
						birch.EC.Int32("queueLength", 0),
					)),
				)),
			)),
		)),
	)

	include := map[string]struct{}{
		"start": {},
		"serverStatus.wiredTiger.cache.bytes currently in the cache":          {},
		"serverStatus.wiredTiger.concurrentTransactions.monitor.available":      {},
		"serverStatus.wiredTiger.concurrentTransactions.monitor.totalTickets":   {},
		"serverStatus.wiredTiger.concurrentTransactions.monitor.queueLength":    {},
	}

	got := normalizeMetricsDocument(doc, include)

	want := map[string]interface{}{
		"start": int64(1700000000000),
		"serverStatus.wiredTiger.cache.bytes currently in the cache":          int64(4096),
		"serverStatus.wiredTiger.concurrentTransactions.monitor.available":    int32(100),
		"serverStatus.wiredTiger.concurrentTransactions.monitor.totalTickets": int32(128),
		"serverStatus.wiredTiger.concurrentTransactions.monitor.queueLength":  int32(0),
	}

	for key, wantVal := range want {
		gotVal, ok := got[key]
		if !ok {
			t.Fatalf("missing key %q; got %#v", key, got)
		}
		if gotVal != wantVal {
			t.Fatalf("key %q: got %v (%T), want %v (%T)", key, gotVal, gotVal, wantVal, wantVal)
		}
	}

	if len(got) != len(want) {
		t.Fatalf("unexpected extra keys: %#v", got)
	}
}

func TestNormalizeMetricsDocument_flatQueueExecutionKeys(t *testing.T) {
	doc := birch.NewDocument(
		birch.EC.Int32("serverStatus.queues.execution.read.available", 100),
		birch.EC.Int32("serverStatus.queues.execution.write.available", 95),
		birch.EC.Int32("serverStatus.queues.execution.read.totalTickets", 128),
	)

	include := map[string]struct{}{
		"serverStatus.queues.execution.read.available":  {},
		"serverStatus.queues.execution.write.available": {},
		"serverStatus.queues.execution.read.totalTickets": {},
	}

	got := normalizeMetricsDocument(doc, include)
	if got["serverStatus.queues.execution.read.available"] != int32(100) {
		t.Fatalf("read.available = %#v", got)
	}
}

func TestNormalizeDocument_preservesNestedMetadata(t *testing.T) {
	doc := birch.NewDocument(
		birch.EC.SubDocument("doc", birch.NewDocument(
			birch.EC.SubDocument("hostInfo", birch.NewDocument(
				birch.EC.SubDocument("system", birch.NewDocument(
					birch.EC.String("hostname", "mongod1.example:27017"),
				)),
			)),
			birch.EC.SubDocument("buildInfo", birch.NewDocument(
				birch.EC.String("version", "8.0.17-6"),
			)),
		)),
	)

	got := normalizeDocument(doc, map[string]struct{}{})
	hostname := getNestedString(got, "doc.hostInfo.system.hostname")
	if hostname != "mongod1.example:27017" {
		t.Fatalf("hostname = %q, want mongod1.example:27017", hostname)
	}
	version := getNestedString(got, "doc.buildInfo.version")
	if version != "8.0.17-6" {
		t.Fatalf("version = %q, want 8.0.17-6", version)
	}
}
