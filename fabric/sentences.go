package fabric

import (
	"context"
	"fmt"
	_ "log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

const sentence_dir = "/opt/agentsentences"
const sentence_pattern = "as-*.bin"

// GetSentences calls GetSentences in a loop unless stopped by the passed-in
// context. If desired, the caller can run this on a goroutine or in a loop
// with an expiring context to get periodic interruptions.
func GetSentences(ctx context.Context, handler func(time.Time, []byte) error) error {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			//log.Printf ("getting sentence")
			if e := GetSentence(ctx, handler); e != nil {
				time.Sleep(200 * time.Millisecond)
			}
		}
	}
	return nil
}

// GetSentence accesses the first available local sentence file
// and passes the data in the file to the supplied handler function.
// If the handler function returns no-error, then the file is removed.
// Note that agent sentences come to us with a prepended timestamp,
// so strip that off and pass it separately to the handler.
func GetSentence(ctx context.Context, handler func(time.Time, []byte) error) error {

	if fm, e := first_matching_file(sentence_dir, sentence_pattern); e == nil {
		if d, e := os.ReadFile(fm); e == nil {

			var ts time.Time
			var data []byte

			if len(d) >= 16 {
				ms, _ := strconv.ParseInt(string(d[0:16]), 10, 64)
				ts = time.UnixMilli(ms)
				data = d[16:]
			} else {
				data = d
			}

			if e := handler(ts, data); e == nil {
				// handler no-error
				os.Remove(fm)
				return nil
			} else {
				// handler error
				// handler errors do not remove the file, which means it
				// will be returned again on the next call to GetSentence.
				return e
			}
		} else {
			// unable to read the contents of the file
			os.Remove(fm)
			return e
		}
	} else {
		// unable to read the first file
		return e
	}

}

// first_matching_file finds files in dir that match the given glob pattern,
// sorts them lexically, and returns the first one.
func first_matching_file(dir, pattern string) (string, error) {
	// Build a full glob pattern scoped to the directory
	fullPattern := filepath.Join(dir, pattern)

	matches, err := filepath.Glob(fullPattern)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no files match pattern %q in %q", pattern, dir)
	}

	sort.Strings(matches)
	return matches[0], nil
}
