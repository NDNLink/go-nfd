package emission_test

import (
	"testing"

	"go-nfd/core/emission"
	"go-nfd/dpdk/dpdktestenv"
)

func TestOnCancel(t *testing.T) {
	assert, _ := dpdktestenv.MakeAR(t)

	nA, nB := 0, 0
	fA := func() { nA++ }
	fB := func() { nB++ }

	emitter := emission.NewEmitter()
	cA := emitter.On(1, fA)
	cB := emitter.On(1, fB)

	emitter.EmitSync(1)
	assert.Equal(1, nA)
	assert.Equal(1, nB)

	assert.NoError(cA.Close())
	emitter.EmitSync(1)
	assert.Equal(1, nA)
	assert.Equal(2, nB)

	assert.NoError(cA.Close())
	emitter.EmitSync(1)
	assert.Equal(1, nA)
	assert.Equal(3, nB)

	assert.NoError(cB.Close())
	emitter.EmitSync(1)
	assert.Equal(1, nA)
	assert.Equal(3, nB)
}
