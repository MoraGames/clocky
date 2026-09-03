package events

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/MoraGames/clockyuwu/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type (
	EventsData struct {
		Map        EventsMap
		Keys       EventsKeys
		Stats      EventsStats
		Expiration time.Time
	}

	EventsMap   map[string]*Event
	EventsKeys  []string
	EventsStats struct {
		TotalSetsNum      int
		EnabledSetsNum    int
		EnabledSets       []string
		TotalEventsNum    int
		EnabledEventsNum  int
		EnabledPointsSum  int
		EnabledEffectsNum int
		EnabledEffects    map[string]int
	}

	EventsResetPinnedMessage struct {
		Exist     bool
		ChatID    int64
		MessageID int
	}

	DailyRewardedUser struct {
		User *structs.UserMinimal
		Sets []string
	}

	EffectPresence struct {
		Name   string
		Amount int
	}
)

var (
	PinnedResetMessage      EventsResetPinnedMessage
	HintRewardedUsers       = make(map[string][]DailyRewardedUser)
	Events                  *EventsData
	AssignEventsWithDefault = func(utils types.Utils) {
		Events = NewEventsData(true, utils)
	}
	EventsFileValid = func(utils types.Utils) bool {
		return Events != nil && !Events.Expiration.IsZero() && time.Now().Before(Events.Expiration)
	}
	NormalizeEventsData = func(utils types.Utils) {
		if Events != nil && Events.Expiration.IsZero() {
			Events.Expiration = currentDailyExpiration(time.Now())
		}
	}
)

func NewEventsData(newEffects bool, utils types.Utils) *EventsData {
	ed := &EventsData{
		make(EventsMap),
		make(EventsKeys, 0),
		EventsStats{0, 0, nil, 0, 0, 0, 0, make(map[string]int)},
		currentDailyExpiration(time.Now()),
	}

	ed.EnabledRandomSets(types.Interval{Min: 0.65, Max: 1.00}, types.Interval{Min: 0.10, Max: 0.20}, utils)

	now := time.Now()
	for i := 0; i < 24*60; i++ {
		time := time.Date(now.Year(), now.Month(), now.Day(), i/60, i%60, 0, 0, now.Location())

		if CalculateValid(time) {
			event := NewEvent(time)
			ed.Map[event.Name] = event
			ed.Keys = append(ed.Keys, event.Name)

			ed.Stats.TotalEventsNum++
			if event.Enabled {
				ed.Stats.EnabledEventsNum++
				ed.Stats.EnabledPointsSum += event.Points
			}
		}
	}
	ed.AssignJokerFormats()

	if newEffects {
		ed.AssignRandomEffects(
			utils,
			structs.EffectPresence{Effect: structs.QuintupleNegativePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}}, // "Mul-5" ->  10% of (95E: 01-02 | 218E: 02-04)
			structs.EffectPresence{Effect: structs.QuadrupleNegativePoints, Possible: 0.40, Amount: types.Interval{Min: 0.02, Max: 0.03}}, // "Mul-4" ->  40% of (95E: 02-03 | 218E: 04-07)
			structs.EffectPresence{Effect: structs.TripleNegativePoints, Possible: 0.70, Amount: types.Interval{Min: 0.03, Max: 0.05}},    // "Mul-3" ->  70% of (95E: 03-05 | 218E: 07-11)
			structs.EffectPresence{Effect: structs.DoubleNegativePoints, Possible: 1.00, Amount: types.Interval{Min: 0.05, Max: 0.10}},    // "Mul-2" -> 100% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.DoublePositivePoints, Possible: 1.00, Amount: types.Interval{Min: 0.16, Max: 0.30}},    // "Mul+2" -> 100% of (95E: 08-14 | 218E: 17-33)
			structs.EffectPresence{Effect: structs.TriplePositivePoints, Possible: 0.75, Amount: types.Interval{Min: 0.10, Max: 0.20}},    // "Mul+3" ->  75% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.QuadruplePositivePoints, Possible: 0.50, Amount: types.Interval{Min: 0.06, Max: 0.10}}, // "Mul+4" ->  50% of (95E: 03-05 | 218E: 07-11)
			structs.EffectPresence{Effect: structs.QuintuplePositivePoints, Possible: 0.25, Amount: types.Interval{Min: 0.04, Max: 0.06}}, // "Mul+5" ->  25% of (95E: 02-03 | 218E: 04-07)
			structs.EffectPresence{Effect: structs.SixtuplePositivePoints, Possible: 0.10, Amount: types.Interval{Min: 0.02, Max: 0.04}},  // "Mul+6" ->  10% of (95E: 01-02 | 218E: 02-04)
			structs.EffectPresence{Effect: structs.SubFourPoints, Possible: 0.25, Amount: types.Interval{Min: 0.02, Max: 0.05}},           // "Sub 4" ->  25% of (95E: 02-05 | 218E: 04-11)
			structs.EffectPresence{Effect: structs.SubThreePoints, Possible: 0.50, Amount: types.Interval{Min: 0.05, Max: 0.10}},          // "Sub 3" ->  50% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.SubTwoPoints, Possible: 1.00, Amount: types.Interval{Min: 0.10, Max: 0.20}},            // "Sub 2" -> 100% of (95E: 10-19 | 218E: 22-44)
			structs.EffectPresence{Effect: structs.AddTwoPoints, Possible: 1.00, Amount: types.Interval{Min: 0.20, Max: 0.40}},            // "Add 2" -> 100% of (95E: 10-19 | 218E: 22-44)
			structs.EffectPresence{Effect: structs.AddThreePoints, Possible: 0.50, Amount: types.Interval{Min: 0.10, Max: 0.20}},          // "Add 3" ->  50% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.AddFourPoints, Possible: 0.25, Amount: types.Interval{Min: 0.04, Max: 0.10}},           // "Add 4" ->  25% of (95E: 02-05 | 218E: 04-11)
			structs.EffectPresence{Effect: structs.AddFivePoints, Possible: 0.10, Amount: types.Interval{Min: 0.02, Max: 0.04}},           // "Add 5" ->  10% of (95E: 01-02 | 218E: 02-04)
		)
	}

	return ed
}

func (ed *EventsData) Reset(newEffects bool, writeMsgData *types.WriteMessageData, utils types.Utils) {
	ed.Stats = EventsStats{0, 0, nil, 0, 0, 0, 0, make(map[string]int)}
	ed.Expiration = currentDailyExpiration(time.Now())
	ed.EnabledRandomSets(types.Interval{Min: 0.65, Max: 1.0}, types.Interval{Min: 0.10, Max: 0.20}, utils)

	for eventName := range ed.Map {
		ed.Map[eventName].Reset()

		ed.Stats.TotalEventsNum++
		if ed.Map[eventName].Enabled {
			ed.Stats.EnabledEventsNum++
			ed.Stats.EnabledPointsSum += ed.Map[eventName].Points
		}
	}
	ed.AssignJokerFormats()

	if newEffects {
		ed.AssignRandomEffects(
			utils,
			structs.EffectPresence{Effect: structs.QuintupleNegativePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}}, // "Mul-5" ->  10% of (95E: 01-02 | 218E: 02-04)
			structs.EffectPresence{Effect: structs.QuadrupleNegativePoints, Possible: 0.40, Amount: types.Interval{Min: 0.02, Max: 0.03}}, // "Mul-4" ->  40% of (95E: 02-03 | 218E: 04-07)
			structs.EffectPresence{Effect: structs.TripleNegativePoints, Possible: 0.70, Amount: types.Interval{Min: 0.03, Max: 0.05}},    // "Mul-3" ->  70% of (95E: 03-05 | 218E: 07-11)
			structs.EffectPresence{Effect: structs.DoubleNegativePoints, Possible: 1.00, Amount: types.Interval{Min: 0.05, Max: 0.10}},    // "Mul-2" -> 100% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.DoublePositivePoints, Possible: 1.00, Amount: types.Interval{Min: 0.16, Max: 0.30}},    // "Mul+2" -> 100% of (95E: 08-14 | 218E: 17-33)
			structs.EffectPresence{Effect: structs.TriplePositivePoints, Possible: 0.75, Amount: types.Interval{Min: 0.10, Max: 0.20}},    // "Mul+3" ->  75% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.QuadruplePositivePoints, Possible: 0.50, Amount: types.Interval{Min: 0.06, Max: 0.10}}, // "Mul+4" ->  50% of (95E: 03-05 | 218E: 07-11)
			structs.EffectPresence{Effect: structs.QuintuplePositivePoints, Possible: 0.25, Amount: types.Interval{Min: 0.04, Max: 0.06}}, // "Mul+5" ->  25% of (95E: 02-03 | 218E: 04-07)
			structs.EffectPresence{Effect: structs.SixtuplePositivePoints, Possible: 0.10, Amount: types.Interval{Min: 0.02, Max: 0.04}},  // "Mul+6" ->  10% of (95E: 01-02 | 218E: 02-04)
			structs.EffectPresence{Effect: structs.SubFourPoints, Possible: 0.25, Amount: types.Interval{Min: 0.02, Max: 0.05}},           // "Sub 4" ->  25% of (95E: 02-05 | 218E: 04-11)
			structs.EffectPresence{Effect: structs.SubThreePoints, Possible: 0.50, Amount: types.Interval{Min: 0.05, Max: 0.10}},          // "Sub 3" ->  50% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.SubTwoPoints, Possible: 1.00, Amount: types.Interval{Min: 0.10, Max: 0.20}},            // "Sub 2" -> 100% of (95E: 10-19 | 218E: 22-44)
			structs.EffectPresence{Effect: structs.AddTwoPoints, Possible: 1.00, Amount: types.Interval{Min: 0.20, Max: 0.40}},            // "Add 2" -> 100% of (95E: 10-19 | 218E: 22-44)
			structs.EffectPresence{Effect: structs.AddThreePoints, Possible: 0.50, Amount: types.Interval{Min: 0.10, Max: 0.20}},          // "Add 3" ->  50% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.AddFourPoints, Possible: 0.25, Amount: types.Interval{Min: 0.04, Max: 0.10}},           // "Add 4" ->  25% of (95E: 02-05 | 218E: 04-11)
			structs.EffectPresence{Effect: structs.AddFivePoints, Possible: 0.10, Amount: types.Interval{Min: 0.02, Max: 0.04}},           // "Add 5" ->  10% of (95E: 01-02 | 218E: 02-04)
		)
	}

	// Save on file the new data
	ed.SaveOnFile(utils)

	// Write Reset Message
	if writeMsgData != nil {
		ed.WriteResetMessage(writeMsgData, utils)
	}
}

func (ed *EventsData) EnabledRandomSets(percentage types.Interval, jokerPercentage types.Interval, utils types.Utils) error {
	if percentage.Min < 0 {
		return fmt.Errorf("minPercentage must be >= 0")
	} else if percentage.Max > 1 {
		return fmt.Errorf("maxPercentage must be <= 1")
	} else if percentage.Min > percentage.Max {
		return fmt.Errorf("minPercentage must be <= maxPercentage")
	}
	if jokerPercentage.Min < 0 {
		return fmt.Errorf("minJokerPercentage must be >= 0")
	} else if jokerPercentage.Max > 1 {
		return fmt.Errorf("maxJokerPercentage must be <= 1")
	} else if jokerPercentage.Min > jokerPercentage.Max {
		return fmt.Errorf("minJokerPercentage must be <= maxJokerPercentage")
	}

	ed.Stats.TotalSetsNum = len(Sets)
	for _, set := range Sets {
		set.Enabled = false
		set.Joker = false
	}

	min, max := int(math.Round(percentage.Min*float64(ed.Stats.TotalSetsNum))), int(math.Round(percentage.Max*float64(ed.Stats.TotalSetsNum)))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	setToActivate := r.Intn(max-min+1) + min

	minJoker, maxJoker := int(math.Round(jokerPercentage.Min*float64(setToActivate))), int(math.Round(jokerPercentage.Max*float64(setToActivate)))
	setToActivateJoker := r.Intn(maxJoker-minJoker+1) + minJoker

	for i, j := 0, 0; i < setToActivate; {
		setIndex := r.Intn(ed.Stats.TotalSetsNum)
		if !Sets[setIndex].Enabled {
			Sets[setIndex].Enabled = true
			ed.Stats.EnabledSetsNum++
			ed.Stats.EnabledSets = append(ed.Stats.EnabledSets, Sets[setIndex].Name)
			if j < setToActivateJoker {
				Sets[setIndex].Joker = true
				j++
			}
			i++
		}
	}

	utils.Logger.WithFields(logrus.Fields{
		"tot": ed.Stats.TotalSetsNum,
		"num": ed.Stats.EnabledSetsNum,
		"set": ed.Stats.EnabledSets,
	}).Debug("EnabledSets")

	return nil
}

func (ed *EventsData) AssignRandomEffects(utils types.Utils, effects ...structs.EffectPresence) {
	var r *rand.Rand
	effectsAmountToApply, effectsNamesToApply, totalEffectsAmount := make(map[string]int), make([]string, 0), 0

	for _, effect := range effects {

		// fmt.Println("DEBUG] Testing for effect ", effect.Effect.Name, "with probability", effect.Possible)

		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		if r.Float64() < effect.Possible {
			// Effects will be assigned
			minEventsEffected, maxEventsEffected := int(math.Round(effect.Amount.Min*float64(ed.Stats.EnabledEventsNum))), int(math.Round(effect.Amount.Max*float64(ed.Stats.EnabledEventsNum)))
			if minEventsEffected == maxEventsEffected {
				maxEventsEffected++
			}

			// fmt.Println("DEBUG] Effect ", effect.Effect.Name, "will be assigned to ", minEventsEffected, "-", maxEventsEffected, "events")

			numEventsWithEffect := r.Intn(maxEventsEffected-minEventsEffected) + minEventsEffected
			effectsNamesToApply = append(effectsNamesToApply, effect.Effect.Name)
			effectsAmountToApply[effect.Effect.Name] += numEventsWithEffect
			totalEffectsAmount += numEventsWithEffect

			// fmt.Println("DEBUG] Effect ", effect.Effect.Name, "assigned to ", numEventsWithEffect, "events")
			// fmt.Println("     |- effectsNamesToApply: ", effectsNamesToApply)
			// fmt.Println("     |- effectsAmountToApply: ", effectsAmountToApply)
			// fmt.Println("     |- totalEffectsAmount: ", totalEffectsAmount)
		}
	}

	// Check if are applicable all effects calculated
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for totalEffectsAmount > ed.Stats.EnabledEventsNum*3 {
		// fmt.Println("DEBUG] Reducing effects amount from", totalEffectsAmount, "to", ed.Stats.EnabledEventsNum*3)
		// Remove a random effect
		effectToDecrease := effectsNamesToApply[r.Intn(len(effectsNamesToApply))]
		effectsAmountToApply[effectToDecrease]--
		totalEffectsAmount--
		if effectsAmountToApply[effectToDecrease] == 0 {
			delete(effectsAmountToApply, effectToDecrease)
			effectsNamesToApply = RemoveValue(effectsNamesToApply, effectToDecrease)
		}
		// fmt.Println("DEBUG] Effect decreased is ", effectToDecrease)
		// fmt.Println("     |- effectsNamesToApply: ", effectsNamesToApply)
		// fmt.Println("     |- effectsAmountToApply: ", effectsAmountToApply)
		// fmt.Println("     |- totalEffectsAmount: ", totalEffectsAmount)
	}

	utils.Logger.WithFields(logrus.Fields{
		"toApp": effectsAmountToApply,
	}).Debug("Effects to enable")

	// Apply all effects
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, effectName := range effectsNamesToApply {
		for i := 0; i < effectsAmountToApply[effectName]; {
			eventName := ed.Keys[r.Intn(len(ed.Keys))]
			if ed.Map[eventName].Enabled && len(ed.Map[eventName].Effects) < 3 {
				ed.Map[eventName].AddEffect(structs.Effects[effectName])
				ed.Stats.EnabledEffectsNum++
				ed.Stats.EnabledEffects[effectName]++
				i++
			}
		}
	}

	utils.Logger.WithFields(logrus.Fields{
		"num": ed.Stats.EnabledEffectsNum,
		"map": ed.Stats.EnabledEffects,
	}).Debug("EnabledEffects")
}

func EventsOf(setFunc func(int, int, int, int) bool) []*Event {
	eventsOfSet := make([]*Event, 0)
	for i := 0; i < 24*60; i++ {
		h := i / 60
		m := i % 60
		if setFunc(h/10, h%10, m/10, m%10) {
			eventsOfSet = append(eventsOfSet, Events.Map[fmt.Sprintf("%02d:%02d", h, m)])
		}
	}
	return eventsOfSet
}

func RemoveValue(s []string, value string) []string {
	newS := make([]string, 0, len(s)-1)
	for _, v := range s {
		if v != value {
			newS = append(newS, v)
		}
	}
	return newS
}

func (ed *EventsData) SaveOnFile(utils types.Utils) {
	//Save Sets
	SetsJson = SetFile{
		Slice:      Sets.ToJsonSlice(),
		Expiration: currentDailyExpiration(time.Now()),
	}
	setsFile, err := json.MarshalIndent(SetsJson, "", " ")
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while marshalling Sets data")
	}
	err = os.WriteFile("files/sets.json", setsFile, 0644)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while writing Sets data")
	}

	//Save Events
	ed.Expiration = currentDailyExpiration(time.Now())
	eventsFile, err := json.MarshalIndent(Events, "", " ")
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while marshalling Events data")
	}
	err = os.WriteFile("files/events.json", eventsFile, 0644)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while writing Events data")
	}
}

func (ed *EventsData) WriteResetMessage(writeMsgData *types.WriteMessageData, utils types.Utils) {
	// Sort the data contained by Stats.EnabledSets and Stats.EnabledEffects
	sortedActiveSets := make([]string, len(ed.Stats.EnabledSets))
	copy(sortedActiveSets, ed.Stats.EnabledSets)
	sort.Slice(sortedActiveSets, func(i, j int) bool {
		return sortedActiveSets[i] < sortedActiveSets[j]
	})

	sortedEnabledEffects := make([]EffectPresence, 0, len(ed.Stats.EnabledEffects))
	for effectName, effectNum := range ed.Stats.EnabledEffects {
		sortedEnabledEffects = append(sortedEnabledEffects, EffectPresence{
			Name:   effectName,
			Amount: effectNum,
		})
	}
	sort.Slice(sortedEnabledEffects, func(i, j int) bool {
		return sortedEnabledEffects[i].Name < sortedEnabledEffects[j].Name
	})

	// Generate text
	text := "Gli eventi son stati resettati.\nEcco alcune informazioni:\n\n"
	text += fmt.Sprintf("Schemi: %v/%v\nEventi: %v/%v\nPunti ottenibili: %v\n", ed.Stats.EnabledSetsNum, ed.Stats.TotalSetsNum, ed.Stats.EnabledEventsNum, ed.Stats.TotalEventsNum, ed.Stats.EnabledPointsSum)

	text += fmt.Sprintf("\nSchemi Attivi (%v):\n", ed.Stats.EnabledSetsNum)
	for _, setName := range sortedActiveSets {
		jokerMessage := ""
		if setFound := Sets.Find(setName); setFound != nil && setFound.Joker {
			jokerMessage = "🃏 "
		}
		text += fmt.Sprintf(" | %s%q\n", jokerMessage, setName)
	}

	text += fmt.Sprintf("\nEffetti Attivi (%v):\n", ed.Stats.EnabledEffectsNum)
	for _, effect := range sortedEnabledEffects {
		text += fmt.Sprintf(" | %q = %v\n", effect.Name, effect.Amount)
	}

	text += "\nBuona fortuna!"

	// Send message
	message := tgbotapi.NewMessage(writeMsgData.ChatID, text)
	if writeMsgData.ReplyMessageID != -1 {
		message.ReplyToMessageID = writeMsgData.ReplyMessageID
	}
	msg, err := writeMsgData.Bot.Send(message)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": msg,
		}).Error("Error while sending message")
	}

	// Update the pinned Message
	UpdatePinnedMessage(writeMsgData, utils, msg)
}

func UpdatePinnedMessage(writeMsgData *types.WriteMessageData, utils types.Utils, msgToPin tgbotapi.Message) {
	// Unpin the old reset message if exists
	if PinnedResetMessage.Exist {
		msg, err := writeMsgData.Bot.Send(tgbotapi.UnpinChatMessageConfig{
			ChatID:    PinnedResetMessage.ChatID,
			MessageID: PinnedResetMessage.MessageID,
		})
		if err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"err": err,
				"msg": msg,
			}).Error("Error while unpinning message")
		}
	}

	// Update the pinned reset message
	PinnedResetMessage = EventsResetPinnedMessage{
		true,
		msgToPin.Chat.ID,
		msgToPin.MessageID,
	}

	// Save PinnedResetMessage
	pinnedMessageFile, err := json.MarshalIndent(PinnedResetMessage, "", " ")
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while marshalling Events data")
	}
	err = os.WriteFile("files/pinnedMessage.json", pinnedMessageFile, 0644)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while writing Events data")
	}

	// Pin the new reset message if exists
	if PinnedResetMessage.Exist {
		msg, err := writeMsgData.Bot.Send(tgbotapi.PinChatMessageConfig{
			ChatID:              PinnedResetMessage.ChatID,
			MessageID:           PinnedResetMessage.MessageID,
			DisableNotification: true,
		})
		if err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"err": err,
				"msg": msg,
			}).Error("Error while pinning message")
		}
	}
}
