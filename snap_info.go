package thyme

import (
	"encoding/json"
	"time"
)

//SnapInfo represents aggregated snapshot information
type SnapInfo struct {
	//LastTimestamp is last time the program create snapshot for window
	LastTimestamp time.Time
	//TotalTime is time in Millisecond format
	TotalTime int64
	//WindowName is name of the open window
	WindowName string
}

//SnapData represents window id as key with SnapInfo as information
type SnapData map[int64]*SnapInfo

var snapData SnapData

//GetOrCreate returns snap info with specified key. If key not found, return new SnapInfo object
func (s SnapData) GetOrCreate(key int64) *SnapInfo {
	if snap, found := snapData[key]; found {
		return snap
	}

	return &SnapInfo{
		LastTimestamp: time.Time{},
		TotalTime:     0,
	}
}

//AggregateSnapshot aggregates snapshot data into SnapData
func AggregateSnapshot(snapshot *Snapshot) {
	if snapData == nil {
		snapData = make(SnapData)
	}

	aggregate(snapshot)
}

//PrintSnapData returns snap data as string
func PrintSnapData() string {
	pretty, err := json.Marshal(snapData)
	if err != nil {
		return err.Error()
	}

	return string(pretty)
}

// aggregate aggregates windows time and save total time and window meta info to SnapInfo
func aggregate(snapshot *Snapshot) {
	if snapshot == nil {
		return
	}

	for _, w := range snapshot.Windows {
		info := snapData.GetOrCreate(w.ID)
		lastTime := info.LastTimestamp
		if lastTime.IsZero() {
			lastTime = snapshot.Time
		}

		activeTime := snapshot.Time.Sub(lastTime)
		currentTime := info.TotalTime

		info.LastTimestamp = snapshot.Time
		info.TotalTime = currentTime + int64(activeTime/time.Millisecond)
		info.WindowName = w.Name

		snapData[w.ID] = info
	}
}
