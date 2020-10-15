package coretest


import "testing"
import (
	"testing"

	_ "go-nfd/core"
	"go-nfd/dpdk/dpdktestenv"
)

func testSipHash(t *testing.T) {
	assert, _ := dpdktestenv.MakeAR(t)

	assert.EqualValues(0, C.TestSipHash())
}
