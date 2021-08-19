package tests

import (
	"testing"
)

func benchmarkPut(count int, b *testing.B) {

}

func BenchmarkPut(b *testing.B) {
	cluster := TestCluster{}
	logDir := SetupDataDirectory()
	err := CreateNewClusterEx(3, &cluster, logDir)
	if err != nil {
		b.Fail()
	}
	b.Run("ABC", func(b *testing.B) {
		benchmarkPut(1, b)
	})

	StopCluster(&cluster)
	Cleanup(logDir)
}
