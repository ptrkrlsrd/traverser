package traverser

import (
	"testing"
)

var num = 10000

func BenchmarkReverseString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        NewCacheKey("/abc")
    }
}
