package widgets

import (
	"sort"
	"sync"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/ui"
)

type RecordsWidget struct {
	*ui.Widget
	Records map[string]models.Record
}

var (
	rwInstance *RecordsWidget
	rwOnce     sync.Once
)

func GetRecordsWidget() *RecordsWidget {
	rwOnce.Do(func() {
		widget := ui.NewWidget("widgets/records.jet")

		widget.Pos = app.UIPos{
			X: -160,
			Y: 86,
		}

		rwInstance = &RecordsWidget{
			Widget:  widget,
			Records: make(map[string]models.Record, 0),
		}
	})

	return rwInstance
}

func (rw *RecordsWidget) reload() {
	records := rw.getSortedRecords()

	rw.Data = map[string]any{
		"Records": records[:min(len(records), 10)],
		"Count":   len(records),
	}

	rw.Display()
}

func (rw *RecordsWidget) getSortedRecords() []models.Record {
	records := make([]models.Record, 0, len(rw.Records))
	for _, record := range rw.Records {
		records = append(records, record)
	}

	sort.Slice(records, func(i, j int) bool {
		if records[i].Time == records[j].Time {
			return records[i].CreatedAt.Before(records[j].CreatedAt)
		}
		return records[i].Time < records[j].Time
	})

	return records
}

func (rw *RecordsWidget) SetRecords(records map[string]models.Record) {
	rw.Records = records
	rw.reload()
}
