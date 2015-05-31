package matcher

import "testing"

func BenchmarkScore(b *testing.B) {
	const choice = "/foo/bar/blatti/foobar"
	for i := 0; i < b.N; i++ {
		Score("foo", choice)
		Score("foobar", choice)
		Score("/bar", choice)
	}
}
