package models

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/logger"
)

var ErrTimestamp = errors.New("cannot parse Timestamp")

func MarshalTimestamp(t time.Time) graphql.Marshaler {
	if t.IsZero() {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, err := io.WriteString(w, strconv.Quote(t.Format(time.RFC3339Nano)))
		if err != nil {
			logger.Warnf("could not marshal timestamp: %v", err)
		}
	})
}

func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(string); ok {
		if len(tmpStr) == 0 {
			return time.Time{}, fmt.Errorf("%w: empty string", ErrTimestamp)
		}

		switch tmpStr[0] {
		case '@':
			i, err := strconv.ParseInt(tmpStr[1:], 10, 64)
			if err != nil {
				return time.Time{}, fmt.Errorf("%w: cannot parse @UNIXTIME: %v", ErrTimestamp, err)
			}
			return time.Unix(i, 0), nil
		case '>', '<':
			d, err := time.ParseDuration(tmpStr[1:])
			if err != nil {
				return time.Time{}, fmt.Errorf("%w: cannot parse %v-duration: %v", ErrTimestamp, tmpStr[0], err)
			}
			t := time.Now()
			// Compute point in time:
			if tmpStr[0] == '<' {
				t = t.Add(-d)
			} else {
				t = t.Add(d)
			}

			return t, nil
		}

		return time.Parse(time.RFC3339Nano, tmpStr)
	}

	return time.Time{}, fmt.Errorf("%w: not a string", ErrTimestamp)
}
