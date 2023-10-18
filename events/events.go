package events

import "time"

type EventsMap map[EventKey]*EventValue

func NewEventsMap() EventsMap {
	return make(EventsMap)
}

func (events EventsMap) Add(eventKey EventKey, eventValue *EventValue) {
	events[eventKey] = eventValue
}

func (events EventsMap) Reset() {
	for _, event := range events {
		event.Activated = false
		event.ActivatedBy = ""
		event.ActivatedAt = time.Time{}
	}
}

func (events EventsMap) Clear() {
	for eventKey := range events {
		delete(events, eventKey)
	}
}

var Events = EventsMap{
	"00:00": {2, false, "", time.Time{}, make(map[int64]bool)},
	"00:12": {1, false, "", time.Time{}, make(map[int64]bool)},
	"00:24": {1, false, "", time.Time{}, make(map[int64]bool)},
	"01:01": {1, false, "", time.Time{}, make(map[int64]bool)},
	"01:10": {1, false, "", time.Time{}, make(map[int64]bool)},
	"01:11": {1, false, "", time.Time{}, make(map[int64]bool)},
	"01:23": {2, false, "", time.Time{}, make(map[int64]bool)},
	"01:35": {1, false, "", time.Time{}, make(map[int64]bool)},
	"02:02": {1, false, "", time.Time{}, make(map[int64]bool)},
	"02:10": {1, false, "", time.Time{}, make(map[int64]bool)},
	"02:20": {1, false, "", time.Time{}, make(map[int64]bool)},
	"02:22": {1, false, "", time.Time{}, make(map[int64]bool)},
	"02:34": {1, false, "", time.Time{}, make(map[int64]bool)},
	"02:46": {2, false, "", time.Time{}, make(map[int64]bool)},
	"03:03": {1, false, "", time.Time{}, make(map[int64]bool)},
	"03:21": {1, false, "", time.Time{}, make(map[int64]bool)},
	"03:30": {1, false, "", time.Time{}, make(map[int64]bool)},
	"03:33": {1, false, "", time.Time{}, make(map[int64]bool)},
	"03:45": {1, false, "", time.Time{}, make(map[int64]bool)},
	"03:57": {1, false, "", time.Time{}, make(map[int64]bool)},
	"04:04": {1, false, "", time.Time{}, make(map[int64]bool)},
	"04:20": {1, false, "", time.Time{}, make(map[int64]bool)},
	"04:32": {1, false, "", time.Time{}, make(map[int64]bool)},
	"04:40": {1, false, "", time.Time{}, make(map[int64]bool)},
	"04:44": {1, false, "", time.Time{}, make(map[int64]bool)},
	"04:56": {1, false, "", time.Time{}, make(map[int64]bool)},
	"05:05": {1, false, "", time.Time{}, make(map[int64]bool)},
	"05:31": {1, false, "", time.Time{}, make(map[int64]bool)},
	"05:43": {1, false, "", time.Time{}, make(map[int64]bool)},
	"05:50": {1, false, "", time.Time{}, make(map[int64]bool)},
	"05:55": {1, false, "", time.Time{}, make(map[int64]bool)},
	"06:06": {1, false, "", time.Time{}, make(map[int64]bool)},
	"06:42": {1, false, "", time.Time{}, make(map[int64]bool)},
	"06:54": {1, false, "", time.Time{}, make(map[int64]bool)},
	"07:07": {1, false, "", time.Time{}, make(map[int64]bool)},
	"07:53": {1, false, "", time.Time{}, make(map[int64]bool)},
	"08:08": {1, false, "", time.Time{}, make(map[int64]bool)},
	"09:09": {1, false, "", time.Time{}, make(map[int64]bool)},
	"10:00": {1, false, "", time.Time{}, make(map[int64]bool)},
	"10:01": {1, false, "", time.Time{}, make(map[int64]bool)},
	"10:10": {1, false, "", time.Time{}, make(map[int64]bool)},
	"10:12": {1, false, "", time.Time{}, make(map[int64]bool)},
	"10:24": {1, false, "", time.Time{}, make(map[int64]bool)},
	"11:11": {2, false, "", time.Time{}, make(map[int64]bool)},
	"11:23": {1, false, "", time.Time{}, make(map[int64]bool)},
	"11:35": {1, false, "", time.Time{}, make(map[int64]bool)},
	"12:10": {1, false, "", time.Time{}, make(map[int64]bool)},
	"12:12": {1, false, "", time.Time{}, make(map[int64]bool)},
	"12:21": {1, false, "", time.Time{}, make(map[int64]bool)},
	"12:22": {1, false, "", time.Time{}, make(map[int64]bool)},
	"12:34": {2, false, "", time.Time{}, make(map[int64]bool)},
	"12:46": {1, false, "", time.Time{}, make(map[int64]bool)},
	"13:13": {1, false, "", time.Time{}, make(map[int64]bool)},
	"13:21": {1, false, "", time.Time{}, make(map[int64]bool)},
	"13:31": {1, false, "", time.Time{}, make(map[int64]bool)},
	"13:33": {1, false, "", time.Time{}, make(map[int64]bool)},
	"13:45": {1, false, "", time.Time{}, make(map[int64]bool)},
	"13:57": {2, false, "", time.Time{}, make(map[int64]bool)},
	"14:14": {1, false, "", time.Time{}, make(map[int64]bool)},
	"14:20": {1, false, "", time.Time{}, make(map[int64]bool)},
	"14:32": {1, false, "", time.Time{}, make(map[int64]bool)},
	"14:41": {1, false, "", time.Time{}, make(map[int64]bool)},
	"14:44": {1, false, "", time.Time{}, make(map[int64]bool)},
	"14:56": {1, false, "", time.Time{}, make(map[int64]bool)},
	"15:15": {1, false, "", time.Time{}, make(map[int64]bool)},
	"15:31": {1, false, "", time.Time{}, make(map[int64]bool)},
	"15:43": {1, false, "", time.Time{}, make(map[int64]bool)},
	"15:51": {1, false, "", time.Time{}, make(map[int64]bool)},
	"15:55": {1, false, "", time.Time{}, make(map[int64]bool)},
	"16:16": {1, false, "", time.Time{}, make(map[int64]bool)},
	"16:42": {1, false, "", time.Time{}, make(map[int64]bool)},
	"16:54": {1, false, "", time.Time{}, make(map[int64]bool)},
	"17:17": {1, false, "", time.Time{}, make(map[int64]bool)},
	"17:53": {1, false, "", time.Time{}, make(map[int64]bool)},
	"18:18": {1, false, "", time.Time{}, make(map[int64]bool)},
	"19:19": {1, false, "", time.Time{}, make(map[int64]bool)},
	"20:00": {1, false, "", time.Time{}, make(map[int64]bool)},
	"20:02": {1, false, "", time.Time{}, make(map[int64]bool)},
	"20:12": {1, false, "", time.Time{}, make(map[int64]bool)},
	"20:20": {1, false, "", time.Time{}, make(map[int64]bool)},
	"20:24": {1, false, "", time.Time{}, make(map[int64]bool)},
	"21:11": {1, false, "", time.Time{}, make(map[int64]bool)},
	"21:12": {1, false, "", time.Time{}, make(map[int64]bool)},
	"21:21": {1, false, "", time.Time{}, make(map[int64]bool)},
	"21:23": {1, false, "", time.Time{}, make(map[int64]bool)},
	"21:35": {1, false, "", time.Time{}, make(map[int64]bool)},
	"22:10": {1, false, "", time.Time{}, make(map[int64]bool)},
	"22:22": {2, false, "", time.Time{}, make(map[int64]bool)},
	"22:34": {1, false, "", time.Time{}, make(map[int64]bool)},
	"22:46": {1, false, "", time.Time{}, make(map[int64]bool)},
	"23:21": {1, false, "", time.Time{}, make(map[int64]bool)},
	"23:23": {1, false, "", time.Time{}, make(map[int64]bool)},
	"23:32": {1, false, "", time.Time{}, make(map[int64]bool)},
	"23:33": {1, false, "", time.Time{}, make(map[int64]bool)},
	"23:45": {2, false, "", time.Time{}, make(map[int64]bool)},
	"23:57": {1, false, "", time.Time{}, make(map[int64]bool)},
	//TESTS: --------------------------------------------------------------
	//"23:59": {99, false, "", time.Time{}, make(map[int64]bool)},
}
