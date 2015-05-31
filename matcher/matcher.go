package matcher

import (
	"sort"
	"strings"
	"time"

	"github.com/spektroskop/ok/util"
)

type Entry struct {
	Text    string
	Score   float64
	Start   int
	End     int
	Matched []bool
}

type Matches []Entry

func (m Matches) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m Matches) Len() int           { return len(m) }
func (m Matches) Less(i, j int) bool { return m[j].Score < m[i].Score }

func Score(search, choice string) (score float64, matched []bool, start, end int) {
	matched = make([]bool, len(choice))

	if index := strings.Index(choice, search); index != -1 {
		start = index
		end = start + len(search)
		score += float64(len(choice)) + 1.0/float64(end-start+1)
		for i := start; i < end; i++ {
			matched[i] = true
		}
		return
	}

	for _, r := range search {
		if index := strings.IndexRune(choice[end:], r); index == -1 {
			score = 0

			return
		} else {
			end += index
			matched[end] = true

			if score == 0 {
				start = end
			}

			score += 1
			end += 1
		}
	}

	score += 1.0 / float64(end-start+1)

	next := start
	if next == 0 {
		next = 1
	}
	score1, matched1, start1, end1 := Score(search, choice[next:])
	if score1 > score {
		for i := 0; i < next; i++ {
			matched[i] = false
		}
		for i, v := range matched1 {
			matched[i+next] = v
		}
		return score1, matched, next + start1, next + end1
	}

	return
}

func Run(search string, choices []string, matchChan chan<- Matches, doneChan <-chan bool) {
	var matches Matches

	if len(search) == 0 {
		goto End
	}

	for _, choice := range choices {
		select {
		case <-doneChan:
			util.Debugf("Cancel `%s'\n", search)
			return
		default:
			if score, matched, start, end := Score(search, choice); score > 0 {
				matches = append(matches, Entry{choice, score, start, end, matched})
			}
		}
	}

	sort.Sort(matches)

End:
	select {
	case matchChan <- matches:
	case <-time.After(time.Millisecond * 100):
		panic("Timeout `matchChan <- matches'")
	}
}
