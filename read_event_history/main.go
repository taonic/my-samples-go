package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gogo/protobuf/jsonpb"
	historypb "go.temporal.io/api/history/v1"
	"go.temporal.io/sdk/client"
)

func main() {
	filename := "/Users/tao/code/temporal/my-samples-go/read_event_history/history.json"

	hist, err := extractHistoryFromFile(filename, 10)
	if err != nil {
		panic(err)
	}

	reader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	hist2, err := client.HistoryFromJSON(reader, client.HistoryJSONOptions{LastEventID: 10})
	if closeErr := reader.Close(); closeErr != nil && err == nil {
		err = closeErr
	}

	fmt.Println(hist.Events[0].EventTime)
	fmt.Println(hist2.Events[0].EventTime)
}

// HistoryFromJSON deserializes history from a reader of JSON bytes. This does
// not close the reader if it is closeable.
func HistoryFromJSON(r io.Reader, lastEventID int64) (*historypb.History, error) {
	unmarshaler := jsonpb.Unmarshaler{
		// Allow unknown fields because if the histroy was generated with a different version of the proto
		// fields may have been added/removed.
		AllowUnknownFields: true,
	}
	var hist historypb.History
	if err := unmarshaler.Unmarshal(r, &hist); err != nil {
		return nil, err
	}

	// If there is a last event ID, slice the rest off
	if lastEventID > 0 {
		for i, event := range hist.Events {
			if event.EventId == lastEventID {
				// Inclusive
				hist.Events = hist.Events[:i+1]
				break
			}
		}
	}
	return &hist, nil
}

func extractHistoryFromFile(jsonfileName string, lastEventID int64) (*historypb.History, error) {
	reader, err := os.Open(jsonfileName)
	if err != nil {
		return nil, err
	}
	hist, err := client.HistoryFromJSON(reader, client.HistoryJSONOptions{LastEventID: lastEventID})
	if closeErr := reader.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	return hist, err
}
