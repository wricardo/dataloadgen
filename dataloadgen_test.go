package dataloadgen_test

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/mshaeon/dataloadgen"
	"github.com/stretchr/testify/require"
)

func ExampleLoader() {
	loader := dataloadgen.NewLoader(func(keys []string) (map[string]int, error) {
		errs := make(dataloadgen.ErrorMap[string])
		ret := make(map[string]int)
		for _, key := range keys {
			num, err := strconv.ParseInt(key, 10, 32)
			ret[key] = int(num)
			if err != nil {
				errs[key] = err
			}
		}
		return ret, errs
	},
		dataloadgen.WithBatchCapacity(1),
		dataloadgen.WithWait(16*time.Millisecond),
	)
	one, err := loader.Load("1")
	if err != nil {
		panic(err)
	}
	fmt.Println(one)
	// Output: 1
}

func TestEdgeCases(t *testing.T) {
	var fetches [][]int
	var mu sync.Mutex
	dl := dataloadgen.NewLoader(func(keys []int) (map[int]string, error) {
		mu.Lock()
		fetches = append(fetches, keys)
		mu.Unlock()

		results := make(map[int]string, len(keys))
		errors := make(dataloadgen.ErrorMap[int])

		for i, key := range keys {
			if key%2 == 0 {
				errors[i] = fmt.Errorf("not found")
			} else {
				results[i] = fmt.Sprint(key)
			}
		}
		return results, errors
	},
		dataloadgen.WithBatchCapacity(5),
		dataloadgen.WithWait(1*time.Millisecond),
	)

	t.Run("load function called only once when cached", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			_, err := dl.Load(0)
			require.Error(t, err)
			require.Len(t, fetches, 1)
			require.Len(t, fetches[0], 1)
		}
		for i := 0; i < 2; i++ {
			r, err := dl.Load(1)
			require.NoError(t, err)
			require.Len(t, fetches, 2)
			require.Len(t, fetches[1], 1)
			require.Equal(t, "1", r)
		}
	})
}
