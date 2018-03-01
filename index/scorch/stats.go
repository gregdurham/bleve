//  Copyright (c) 2017 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scorch

import (
	"encoding/json"
	"reflect"
	"sync/atomic"
)

// Stats tracks statistics about the index, fields that are
// prefixed like CurXxxx are gauges (can go up and down),
// and fields that are prefixed like TotXxxx are monotonically
// increasing counters.
type Stats struct {
	TotUpdates      uint64
	TotDeletes      uint64
	TotBatches      uint64
	TotBatchesEmpty uint64
	TotOnErrors     uint64

	TotAnalysisTime uint64
	TotIndexTime    uint64

	TotIndexedPlainTextBytes uint64

	TotTermSearchersStarted  uint64
	TotTermSearchersFinished uint64

	TotIntroduceLoop       uint64
	TotIntroduceSegmentBeg uint64
	TotIntroduceSegmentEnd uint64
	TotIntroduceMergeBeg   uint64
	TotIntroduceMergeEnd   uint64
	TotIntroduceRevertBeg  uint64
	TotIntroduceRevertEnd  uint64

	TotIntroducedItems         uint64
	TotIntroducedSegmentsBatch uint64
	TotIntroducedSegmentsMerge uint64

	TotPersistLoopBeg          uint64
	TotPersistLoopErr          uint64
	TotPersistLoopProgress     uint64
	TotPersistLoopWait         uint64
	TotPersistLoopWaitNotified uint64
	TotPersistLoopEnd          uint64

	TotPersistedItems    uint64
	TotPersistedSegments uint64

	TotPersisterSlowMergerPause  uint64
	TotPersisterSlowMergerResume uint64

	TotFileMergeLoopBeg uint64
	TotFileMergeLoopErr uint64
	TotFileMergeLoopEnd uint64

	TotFileMergePlan     uint64
	TotFileMergePlanErr  uint64
	TotFileMergePlanNone uint64
	TotFileMergePlanOk   uint64

	TotFileMergePlanTasks              uint64
	TotFileMergePlanTasksDone          uint64
	TotFileMergePlanTasksErr           uint64
	TotFileMergePlanTasksSegments      uint64
	TotFileMergePlanTasksSegmentsEmpty uint64

	TotFileMergeSegmentsEmpty uint64
	TotFileMergeSegments      uint64

	TotFileMergeZapBeg uint64
	TotFileMergeZapEnd uint64

	TotFileMergeIntroductions     uint64
	TotFileMergeIntroductionsDone uint64

	TotMemMergeBeg      uint64
	TotMemMergeErr      uint64
	TotMemMergeDone     uint64
	TotMemMergeZapBeg   uint64
	TotMemMergeZapEnd   uint64
	TotMemMergeSegments uint64
}

// atomically populates the returned map
func (s *Stats) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	sve := reflect.ValueOf(s).Elem()
	svet := sve.Type()
	for i := 0; i < svet.NumField(); i++ {
		svef := sve.Field(i)
		if svef.CanAddr() {
			svefp := svef.Addr().Interface()
			m[svet.Field(i).Name] = atomic.LoadUint64(svefp.(*uint64))
		}
	}
	return m
}

// MarshalJSON implements json.Marshaler, and in contrast to standard
// json marshaling provides atomic safety
func (s *Stats) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToMap())
}
