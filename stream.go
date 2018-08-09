package thyme

import (
	"bytes"
	"fmt"
	"os"
)

// Stream represents all the sampling data gathered by Thyme.
type Stream struct {
	// Snapshots is a list of window snapshots ordered by time.
	Snapshots []*Snapshot
}

// Add adds snapshot to stream
func (s *Stream) Add(snapshot *Snapshot) {
	s.Snapshots = append(s.Snapshots, snapshot)
}

// Dump writes snapshot data into csv friendly format
func (s *Stream) Dump() string {
	for _, snap := range s.Snapshots {
		AggregateSnapshot(snap)
	}

	fmt.Println(PrintSnapData())
	return PrintSnapData()
}

// Flush flushes snapshots into output file
func (s *Stream) Flush(filename string) {
	if len(s.Snapshots) > 10 {
		f, err := openOrCreate(filename)
		if err != nil {
			panic(err)
		}

		defer f.Close()
		if _, err := f.Stat(); err != nil {
			panic(err)
		}

		content := s.Dump()
		if _, err := f.WriteString(content); err != nil {
			panic(err)
		}

		//Clean snapshot content
		s.Snapshots = []*Snapshot{}
	}
}

// Print returns a pretty-printed representation of the snapshot.
func (s Stream) Print() string {
	var b bytes.Buffer
	for _, snap := range s.Snapshots {
		fmt.Fprintf(&b, "%s", snap.Print())
	}
	return string(b.Bytes())
}

func openOrCreate(filename string) (*os.File, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(filename)
			return f, err
		}
	}

	return f, err
}
