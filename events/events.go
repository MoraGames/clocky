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
	"00:00": {2, false, "", time.Time{}},
	"00:12": {1, false, "", time.Time{}},
	"00:24": {1, false, "", time.Time{}},
	"01:01": {1, false, "", time.Time{}},
	"01:10": {1, false, "", time.Time{}},
	"01:11": {1, false, "", time.Time{}},
	"01:23": {2, false, "", time.Time{}},
	"01:35": {1, false, "", time.Time{}},
	"02:02": {1, false, "", time.Time{}},
	"02:10": {1, false, "", time.Time{}},
	"02:20": {1, false, "", time.Time{}},
	"02:22": {1, false, "", time.Time{}},
	"02:34": {1, false, "", time.Time{}},
	"02:46": {2, false, "", time.Time{}},
	"03:03": {1, false, "", time.Time{}},
	"03:21": {1, false, "", time.Time{}},
	"03:30": {1, false, "", time.Time{}},
	"03:33": {1, false, "", time.Time{}},
	"03:45": {1, false, "", time.Time{}},
	"03:57": {1, false, "", time.Time{}},
	"04:04": {1, false, "", time.Time{}},
	"04:20": {1, false, "", time.Time{}},
	"04:32": {1, false, "", time.Time{}},
	"04:40": {1, false, "", time.Time{}},
	"04:44": {1, false, "", time.Time{}},
	"04:56": {1, false, "", time.Time{}},
	"05:05": {1, false, "", time.Time{}},
	"05:31": {1, false, "", time.Time{}},
	"05:43": {1, false, "", time.Time{}},
	"05:50": {1, false, "", time.Time{}},
	"05:55": {1, false, "", time.Time{}},
	"06:06": {1, false, "", time.Time{}},
	"06:42": {1, false, "", time.Time{}},
	"06:54": {1, false, "", time.Time{}},
	"07:07": {1, false, "", time.Time{}},
	"07:53": {1, false, "", time.Time{}},
	"08:08": {1, false, "", time.Time{}},
	"09:09": {1, false, "", time.Time{}},
	"10:00": {1, false, "", time.Time{}},
	"10:01": {1, false, "", time.Time{}},
	"10:10": {1, false, "", time.Time{}},
	"10:12": {1, false, "", time.Time{}},
	"10:24": {1, false, "", time.Time{}},
	"11:11": {2, false, "", time.Time{}},
	"11:23": {1, false, "", time.Time{}},
	"11:35": {1, false, "", time.Time{}},
	"12:10": {1, false, "", time.Time{}},
	"12:12": {1, false, "", time.Time{}},
	"12:21": {1, false, "", time.Time{}},
	"12:22": {1, false, "", time.Time{}},
	"12:34": {2, false, "", time.Time{}},
	"12:46": {1, false, "", time.Time{}},
	"13:13": {1, false, "", time.Time{}},
	"13:21": {1, false, "", time.Time{}},
	"13:31": {1, false, "", time.Time{}},
	"13:33": {1, false, "", time.Time{}},
	"13:45": {1, false, "", time.Time{}},
	"13:57": {2, false, "", time.Time{}},
	"14:14": {1, false, "", time.Time{}},
	"14:20": {1, false, "", time.Time{}},
	"14:32": {1, false, "", time.Time{}},
	"14:41": {1, false, "", time.Time{}},
	"14:44": {1, false, "", time.Time{}},
	"14:56": {1, false, "", time.Time{}},
	"15:15": {1, false, "", time.Time{}},
	"15:31": {1, false, "", time.Time{}},
	"15:43": {1, false, "", time.Time{}},
	"15:51": {1, false, "", time.Time{}},
	"15:55": {1, false, "", time.Time{}},
	"16:16": {1, false, "", time.Time{}},
	"16:42": {1, false, "", time.Time{}},
	"16:54": {1, false, "", time.Time{}},
	"17:17": {1, false, "", time.Time{}},
	"17:53": {1, false, "", time.Time{}},
	"18:18": {1, false, "", time.Time{}},
	"19:19": {1, false, "", time.Time{}},
	"20:00": {1, false, "", time.Time{}},
	"20:02": {1, false, "", time.Time{}},
	"20:12": {1, false, "", time.Time{}},
	"20:20": {1, false, "", time.Time{}},
	"20:24": {1, false, "", time.Time{}},
	"21:11": {1, false, "", time.Time{}},
	"21:12": {1, false, "", time.Time{}},
	"21:21": {1, false, "", time.Time{}},
	"21:23": {1, false, "", time.Time{}},
	"21:35": {1, false, "", time.Time{}},
	"22:10": {1, false, "", time.Time{}},
	"22:22": {2, false, "", time.Time{}},
	"22:34": {1, false, "", time.Time{}},
	"22:46": {1, false, "", time.Time{}},
	"23:21": {1, false, "", time.Time{}},
	"23:23": {1, false, "", time.Time{}},
	"23:32": {1, false, "", time.Time{}},
	"23:33": {1, false, "", time.Time{}},
	"23:45": {2, false, "", time.Time{}},
	"23:57": {1, false, "", time.Time{}},
	//TESTS: --------------------------------------------------------------
	"23:59": {0, false, "", time.Time{}},
}
