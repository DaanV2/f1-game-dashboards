package randx_test

import (
	"fmt"
	"testing"

	"github.com/DaanV2/f1-game-dashboards/server/pkg/randx"
	"github.com/stretchr/testify/require"
)

func Test_GenerateBase64(t *testing.T) {
	sizes := []int{16, 32, 64, 128, 256, 512, 1024}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("size %d", size), func(t *testing.T) {
			got, err := randx.GenerateBase64(size)
			require.NoError(t, err)
			require.Len(t, got, size)
		})
	}
}

func Fuzz_GenerateBase64(f *testing.F) {
	f.Add(16)
	f.Add(1024)

	f.Fuzz(func(t *testing.T, size int) {
		got, err := randx.GenerateBase64(size)
		require.NoError(t, err)
		require.Len(t, got, size)
	})
}
