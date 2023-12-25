package events

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/MoraGames/clockyuwu/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type EventsMap map[EventKey]*EventValue

type EventsNumbers struct {
	Total      int
	Active     int
	Uneffected int
	Effected   map[string]int
}

func NewEventsMap() EventsMap {
	return make(EventsMap)
}

func (events EventsMap) Add(eventKey EventKey, eventValue *EventValue) {
	events[eventKey] = eventValue
}

func (events EventsMap) Reset(writeMessage bool, wrtMsgData types.WriteMessageData) {
	for _, event := range events {
		event.Activated = false
		event.ActivatedBy = ""
		event.ActivatedAt = time.Time{}
		event.ArrivedAt = time.Time{}
		event.Partecipations = make(map[int64]bool)
		event.Effects = make([]structs.Effect, 0)
	}
	evntsNums := events.RandomizeEffects()

	if writeMessage {
		text := fmt.Sprintf("Gli eventi son stati resettati.\nEcco alcune informazioni:\n\nNumero eventi %v/%v (%v senza effetti).\n\nEffetti Attivi:\n", evntsNums.Active, evntsNums.Total, evntsNums.Uneffected)
		for key, value := range evntsNums.Effected {
			text += fmt.Sprintf("  %v = %v\n", key, value)
		}
		text += "\nBuona fortuna!"

		message := tgbotapi.NewMessage(wrtMsgData.ChatID, text)
		if wrtMsgData.ReplyMessageID != -1 {
			message.ReplyToMessageID = wrtMsgData.ReplyMessageID
		}
		wrtMsgData.Bot.Send(message)
	}
}

func (events EventsMap) RandomizeEffects() EventsNumbers {
	eventsNumber := len(events)
	eventsNumberEditable := eventsNumber

	// 20% of events without other effects will have x2pts effect.
	doublePtsEvents := int(float64(eventsNumber) * 0.2)

	if doublePtsEvents > eventsNumberEditable {
		doublePtsEvents = eventsNumberEditable
	}
	for i := 0; i < doublePtsEvents; {
		eventIndex := rand.Intn(eventsNumber)
		event := events[EventsKeys[eventIndex]]
		if len(event.Effects) == 0 {
			event.Effects = append(event.Effects, structs.Effect{Name: "Double Points", Key: "x", Value: 2})
			i++
			eventsNumberEditable--
		}
	}

	// 05% of events without other effects will have x3pts effect.
	triplePtsEvents := int(float64(eventsNumber) * 0.05)

	if triplePtsEvents > eventsNumberEditable {
		triplePtsEvents = eventsNumberEditable
	}
	for i := 0; i < triplePtsEvents; {
		eventIndex := rand.Intn(eventsNumber)
		event := events[EventsKeys[eventIndex]]
		if len(event.Effects) == 0 {
			event.Effects = append(event.Effects, structs.Effect{Name: "Triple Points", Key: "x", Value: 3})
			i++
			eventsNumberEditable--
		}
	}

	// 10% of events without other effects will have x-1pts effect.
	negativePtsEvents := int(float64(eventsNumber) * 0.1)

	if negativePtsEvents > eventsNumberEditable {
		negativePtsEvents = eventsNumberEditable
	}
	for i := 0; i < negativePtsEvents; {
		eventIndex := rand.Intn(eventsNumber)
		event := events[EventsKeys[eventIndex]]
		if len(event.Effects) == 0 {
			event.Effects = append(event.Effects, structs.Effect{Name: "Negative Points", Key: "x", Value: -1})
			i++
			eventsNumberEditable--
		}
	}

	// With 30% of probability, one event without other effect will have x5pts effect.
	quintuplesPtsEvents := 0
	if rand.Float64() < 0.3 {
		quintuplesPtsEvents = 1

		if quintuplesPtsEvents > eventsNumberEditable {
			quintuplesPtsEvents = eventsNumberEditable
		}

		for i := 0; i < quintuplesPtsEvents; {
			eventIndex := rand.Intn(eventsNumber)
			event := events[EventsKeys[eventIndex]]
			if len(event.Effects) == 0 {
				event.Effects = append(event.Effects, structs.Effect{Name: "Quintuples Points", Key: "x", Value: 5})
				i++
				eventsNumberEditable--
			}
		}
	}

	// With 30% of probability, one event without other effect will have x-2pts effect.
	doubleNegativePtsEvents := 0
	if rand.Float64() < 0.3 {
		doubleNegativePtsEvents = 1

		if doubleNegativePtsEvents > eventsNumberEditable {
			doubleNegativePtsEvents = eventsNumberEditable
		}

		for i := 0; i < doubleNegativePtsEvents; {
			eventIndex := rand.Intn(eventsNumber)
			event := events[EventsKeys[eventIndex]]
			if len(event.Effects) == 0 {
				event.Effects = append(event.Effects, structs.Effect{Name: "Double Negative Points", Key: "x", Value: -2})
				i++
				eventsNumberEditable--
			}
		}
	}

	return EventsNumbers{
		Total:      eventsNumber,
		Active:     eventsNumber,
		Uneffected: eventsNumberEditable,
		Effected: map[string]int{
			"x(-2)": doubleNegativePtsEvents,
			"x(-1)": negativePtsEvents,
			"x(+2)": doublePtsEvents,
			"x(+3)": triplePtsEvents,
			"x(+5)": quintuplesPtsEvents,
		},
	}
}

func (events EventsMap) Clear() {
	for eventKey := range events {
		delete(events, eventKey)
	}
}

var Events = EventsMap{
	"00:00": {3, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"00:12": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"00:24": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"01:01": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"01:02": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"01:10": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"01:11": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"01:23": {2, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"01:35": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"02:02": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"02:04": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"02:10": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"02:20": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"02:22": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"02:34": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"02:46": {2, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"03:03": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"03:06": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"03:21": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"03:30": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"03:33": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"03:45": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"03:57": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"04:04": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"04:08": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"04:20": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"04:32": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"04:40": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"04:44": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"04:56": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"05:05": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"05:10": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"05:31": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"05:43": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"05:50": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"05:55": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"06:06": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"06:12": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"06:42": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"06:54": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"07:07": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"07:14": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"07:53": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"08:08": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"08:16": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"09:09": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"08:18": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"10:00": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"10:01": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"10:10": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"10:12": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"10:20": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"10:24": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"11:11": {2, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"11:22": {2, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"11:23": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"11:35": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"12:10": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"12:12": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"12:21": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"12:22": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"12:24": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"12:34": {2, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"12:46": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"13:13": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"13:21": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"13:26": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"13:31": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"13:33": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"13:45": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"13:57": {2, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"14:14": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"14:20": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"14:28": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"14:32": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"14:41": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"14:44": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"14:56": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"15:15": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"15:30": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"15:31": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"15:43": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"15:51": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"15:55": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"16:16": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"16:32": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"16:42": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"16:54": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"17:17": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"17:34": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"17:53": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"18:18": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"18:56": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"19:19": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"19:38": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"20:00": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"20:02": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"20:12": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"20:20": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"20:40": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"20:24": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"21:11": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"21:12": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"21:21": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"21:23": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"21:35": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"21:42": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"22:10": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"22:22": {2, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"22:34": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"22:44": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"22:46": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"23:21": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"23:23": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"23:32": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"23:33": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"23:45": {2, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"23:46": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	"23:57": {1, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	//TESTS: --------------------------------------------------------------
	//"21:37": {0, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	//"21:38": {0, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
	//"23:59": {0, false, "", time.Time{}, time.Time{}, make(map[int64]bool), make([]structs.Effect, 0)},
}

var EventsKeys = []EventKey{
	"00:00",
	"00:12",
	"00:24",
	"01:01",
	"01:02",
	"01:10",
	"01:11",
	"01:23",
	"01:35",
	"02:02",
	"02:04",
	"02:10",
	"02:20",
	"02:22",
	"02:34",
	"02:46",
	"03:03",
	"03:06",
	"03:21",
	"03:30",
	"03:33",
	"03:45",
	"03:57",
	"04:04",
	"04:08",
	"04:20",
	"04:32",
	"04:40",
	"04:44",
	"04:56",
	"05:05",
	"05:10",
	"05:31",
	"05:43",
	"05:50",
	"05:55",
	"06:06",
	"06:12",
	"06:42",
	"06:54",
	"07:07",
	"07:14",
	"07:53",
	"08:08",
	"08:16",
	"09:09",
	"08:18",
	"10:00",
	"10:01",
	"10:10",
	"10:12",
	"10:20",
	"10:24",
	"11:11",
	"11:22",
	"11:23",
	"11:35",
	"12:10",
	"12:12",
	"12:21",
	"12:22",
	"12:24",
	"12:34",
	"12:46",
	"13:13",
	"13:21",
	"13:26",
	"13:31",
	"13:33",
	"13:45",
	"13:57",
	"14:14",
	"14:20",
	"14:28",
	"14:32",
	"14:41",
	"14:44",
	"14:56",
	"15:15",
	"15:30",
	"15:31",
	"15:43",
	"15:51",
	"15:55",
	"16:16",
	"16:32",
	"16:42",
	"16:54",
	"17:17",
	"17:34",
	"17:53",
	"18:18",
	"18:56",
	"19:19",
	"19:38",
	"20:00",
	"20:02",
	"20:12",
	"20:20",
	"20:40",
	"20:24",
	"21:11",
	"21:12",
	"21:21",
	"21:23",
	"21:35",
	"21:42",
	"22:10",
	"22:22",
	"22:34",
	"22:44",
	"22:46",
	"23:21",
	"23:23",
	"23:32",
	"23:33",
	"23:45",
	"23:46",
	"23:57",
	//TESTS: --------------------------------------------------------------
	//"21:37",
	//"21:38",
	//"23:59",
}
