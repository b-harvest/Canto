package types_test

import (
	"math"
	"sort"
	"testing"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"
)

type IntegerTestStruct struct {
	Id  uint64
	val int64
}

func genRandomInt64() int64 {
	return tmrand.Int63n(math.MaxInt64)
}

func (t IntegerTestStruct) genRandom(id uint64) IntegerTestStruct {
	return IntegerTestStruct{Id: id, val: genRandomInt64()}
}

func (t IntegerTestStruct) hasSameId(other IntegerTestStruct) bool {
	return t.Id == other.Id
}

func (t IntegerTestStruct) less(other IntegerTestStruct) bool {
	return t.val < other.val
}

func testFilterSlice[T interface {
	genRandom(uint64) T
}](t *testing.T, filter func(T) bool) {
	const iteration = 1000
	var expected []T
	var input []T

	for i := uint64(0); i < iteration; i++ {
		var t T
		t = t.genRandom(i)
		input = append(input, t)
		if filter(t) {
			expected = append(expected, t)
		}
	}
	result := types.FilterSlice(input, filter)
	require.Equal(t, expected, result)
}

func TestFilterSlice(t *testing.T) {
	t.Run("test integer", func(t *testing.T) {
		testFilterSlice(t, func(t IntegerTestStruct) bool {
			return true
		})
		testFilterSlice(t, func(t IntegerTestStruct) bool {
			return t.Id%2 == 0
		})
		testFilterSlice(t, func(t IntegerTestStruct) bool {
			return t.val%2 == 0
		})
		testFilterSlice(t, func(t IntegerTestStruct) bool {
			return t.Id%4 == 0
		})
		testFilterSlice(t, func(t IntegerTestStruct) bool {
			return t.val%4 == 0
		})
	})
}

func testMapToSortedSlice[T interface {
	genRandom(uint64) T
	hasSameId(T) bool
	less(T) bool
}](t *testing.T) {
	const iteration = 1000
	var expected []T
	input := make(map[uint64]T)

	for i := uint64(0); i < iteration; i++ {
		id := uint64(tmrand.Int63n(100))
		var t T
		t = t.genRandom(id)
		input[id] = t
		expected = types.FilterSlice(expected, func(elem T) bool {
			return !elem.hasSameId(t)
		})
		expected = append(expected, t)
	}
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].less(expected[j])
	})
	result := types.MapToSortedSlice(input, func(i, j T) bool {
		return i.less(j)
	})
	require.Equal(t, expected, result)
}

func TestMapToSortedSlice(t *testing.T) {
	t.Run("integer", testMapToSortedSlice[IntegerTestStruct])
}
