package task_test

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Parse custom test flags such as -update
	flag.Parse()
	os.Exit(m.Run())
}
