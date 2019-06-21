package audit

import (
	"fmt"
)

// THIS IS TMP IMPL

func StartWorker() {
	auditor, err := NewAuditor(false)
	if err != nil {
		panic(err)
	}

	if err := auditor.WatchMount("/"); err != nil {
		panic(err)
	}

	for ev := range auditor.Event {
		fmt.Println(">>>>>", ev.Acts, ev.FileName, "by", ev.ProcessInfo)
	}
}
