package events

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/MoraGames/clockyuwu/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type (
	EventsData struct {
		Map   EventsMap
		Keys  EventsKeys
		Stats EventsStats
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

	EventsResetMessage struct {
		Exist     bool
		ChatID    int64
		MessageID int
	}
)

var (
	PinnedResetMessage      EventsResetMessage
	HintRewardedUsers       = make(map[string][]*structs.UserMinimal)
	Events                  *EventsData
	AssignEventsWithDefault = func(utils types.Utils) {
		Events = NewEventsData(true, utils)
	}
)

func NewEventsData(newEffects bool, utils types.Utils) *EventsData {
	ed := &EventsData{
		make(EventsMap),
		make(EventsKeys, 0),
		EventsStats{0, 0, nil, 0, 0, 0, 0, make(map[string]int)},
	}

	ed.EnabledRandomSets(types.Interval{Min: 0.65, Max: 1.00}, utils)

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

	if newEffects {
		ed.AssignRandomEffects(
			utils,
			structs.EffectPresence{Effect: structs.QuintupleNegativePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}}, //71E ->  10% of 00-01 effects.  |  119E ->  10% of 01-02 effects.
			structs.EffectPresence{Effect: structs.TripleNegativePoints, Possible: 0.20, Amount: types.Interval{Min: 0.02, Max: 0.03}},    //71E ->  20% of 01-02 effects.  |  119E ->  20% of 02-04 effects.
			structs.EffectPresence{Effect: structs.DoubleNegativePoints, Possible: 0.30, Amount: types.Interval{Min: 0.03, Max: 0.05}},    //71E ->  30% of 02-03 effects.  |  119E ->  30% of 04-05 effects.
			structs.EffectPresence{Effect: structs.SingleNegativePoints, Possible: 1.00, Amount: types.Interval{Min: 0.05, Max: 0.15}},    //71E -> 100% of 03-10 effects.  |  119E -> 100% of 05-17 effects.
			structs.EffectPresence{Effect: structs.DoublePositivePoints, Possible: 0.95, Amount: types.Interval{Min: 0.10, Max: 0.20}},    //71E ->  95% of 07-14 effects.  |  119E ->  95% of 11-23 effects.
			structs.EffectPresence{Effect: structs.TriplePositivePoints, Possible: 0.25, Amount: types.Interval{Min: 0.02, Max: 0.05}},    //71E ->  25% of 01-03 effects.  |  119E ->  25% of 02-05 effects.
			structs.EffectPresence{Effect: structs.QuintuplePositivePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}}, //71E ->  10% of 00-01 effects.  |  119E ->  10% of 01-02 effects.
			structs.EffectPresence{Effect: structs.SixtuplePositivePoints, Possible: 0.05, Amount: types.Interval{Min: 0.01, Max: 0.02}},  //71E ->   5% of 00-01 effects.  |  119E ->   5% of 01-02 effects.
			structs.EffectPresence{Effect: structs.SubFourPoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}},           //71E ->  10% of 00-01 effects.  |  119E ->  10% of 01-02 effects.
			structs.EffectPresence{Effect: structs.SubThreePoints, Possible: 0.30, Amount: types.Interval{Min: 0.03, Max: 0.05}},          //71E ->  30% of 02-03 effects.  |  119E ->  30% of 04-05 effects.
			structs.EffectPresence{Effect: structs.SubTwoPoints, Possible: 0.50, Amount: types.Interval{Min: 0.05, Max: 0.15}},            //71E ->  50% of 03-10 effects.  |  119E ->  50% of 05-17 effects.
			structs.EffectPresence{Effect: structs.SubOnePoint, Possible: 0.95, Amount: types.Interval{Min: 0.10, Max: 0.20}},             //71E ->  95% of 07-14 effects.  |  119E ->  95% of 11-23 effects.
			structs.EffectPresence{Effect: structs.AddOnePoint, Possible: 1.00, Amount: types.Interval{Min: 0.10, Max: 0.25}},             //71E -> 100% of 07-17 effects.  |  119E -> 100% of 11-29 effects.
			structs.EffectPresence{Effect: structs.AddTwoPoints, Possible: 0.95, Amount: types.Interval{Min: 0.10, Max: 0.20}},            //71E ->  95% of 07-14 effects.  |  119E ->  95% of 11-23 effects.
			structs.EffectPresence{Effect: structs.AddThreePoints, Possible: 0.50, Amount: types.Interval{Min: 0.05, Max: 0.15}},          //71E ->  50% of 03-10 effects.  |  119E ->  50% of 05-17 effects.
			structs.EffectPresence{Effect: structs.AddFourPoints, Possible: 0.30, Amount: types.Interval{Min: 0.03, Max: 0.05}},           //71E ->  30% of 02-03 effects.  |  119E ->  30% of 04-05 effects.
			structs.EffectPresence{Effect: structs.AddFivePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}},           //71E ->  10% of 00-01 effects.  |  119E ->  10% of 01-02 effects.
		)
	}

	return ed
}

func (ed *EventsData) Reset(newEffects bool, writeMsgData *types.WriteMessageData, utils types.Utils) {
	ed.Stats = EventsStats{0, 0, nil, 0, 0, 0, 0, make(map[string]int)}
	ed.EnabledRandomSets(types.Interval{Min: 0.65, Max: 1.0}, utils)

	for eventName := range ed.Map {
		ed.Map[eventName].Reset()

		ed.Stats.TotalEventsNum++
		if ed.Map[eventName].Enabled {
			ed.Stats.EnabledEventsNum++
			ed.Stats.EnabledPointsSum += ed.Map[eventName].Points
		}
	}

	if newEffects {
		ed.AssignRandomEffects(
			utils,
			structs.EffectPresence{Effect: structs.QuintupleNegativePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}}, //71E ->  10% of 00-01 effects.  |  119E ->  10% of 01-02 effects.
			structs.EffectPresence{Effect: structs.TripleNegativePoints, Possible: 0.20, Amount: types.Interval{Min: 0.02, Max: 0.03}},    //71E ->  20% of 01-02 effects.  |  119E ->  20% of 02-04 effects.
			structs.EffectPresence{Effect: structs.DoubleNegativePoints, Possible: 0.30, Amount: types.Interval{Min: 0.03, Max: 0.05}},    //71E ->  30% of 02-03 effects.  |  119E ->  30% of 04-05 effects.
			structs.EffectPresence{Effect: structs.SingleNegativePoints, Possible: 1.00, Amount: types.Interval{Min: 0.05, Max: 0.15}},    //71E -> 100% of 03-10 effects.  |  119E -> 100% of 05-17 effects.
			structs.EffectPresence{Effect: structs.DoublePositivePoints, Possible: 0.95, Amount: types.Interval{Min: 0.10, Max: 0.20}},    //71E ->  95% of 07-14 effects.  |  119E ->  95% of 11-23 effects.
			structs.EffectPresence{Effect: structs.TriplePositivePoints, Possible: 0.25, Amount: types.Interval{Min: 0.02, Max: 0.05}},    //71E ->  25% of 01-03 effects.  |  119E ->  25% of 02-05 effects.
			structs.EffectPresence{Effect: structs.QuintuplePositivePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}}, //71E ->  10% of 00-01 effects.  |  119E ->  10% of 01-02 effects.
			structs.EffectPresence{Effect: structs.SixtuplePositivePoints, Possible: 0.05, Amount: types.Interval{Min: 0.01, Max: 0.02}},  //71E ->   5% of 00-01 effects.  |  119E ->   5% of 01-02 effects.
			structs.EffectPresence{Effect: structs.SubFourPoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}},           //71E ->  10% of 00-01 effects.  |  119E ->  10% of 01-02 effects.
			structs.EffectPresence{Effect: structs.SubThreePoints, Possible: 0.30, Amount: types.Interval{Min: 0.03, Max: 0.05}},          //71E ->  30% of 02-03 effects.  |  119E ->  30% of 04-05 effects.
			structs.EffectPresence{Effect: structs.SubTwoPoints, Possible: 0.50, Amount: types.Interval{Min: 0.05, Max: 0.15}},            //71E ->  50% of 03-10 effects.  |  119E ->  50% of 05-17 effects.
			structs.EffectPresence{Effect: structs.SubOnePoint, Possible: 0.95, Amount: types.Interval{Min: 0.10, Max: 0.20}},             //71E ->  95% of 07-14 effects.  |  119E ->  95% of 11-23 effects.
			structs.EffectPresence{Effect: structs.AddOnePoint, Possible: 1.00, Amount: types.Interval{Min: 0.10, Max: 0.25}},             //71E -> 100% of 07-17 effects.  |  119E -> 100% of 11-29 effects.
			structs.EffectPresence{Effect: structs.AddTwoPoints, Possible: 0.95, Amount: types.Interval{Min: 0.10, Max: 0.20}},            //71E ->  95% of 07-14 effects.  |  119E ->  95% of 11-23 effects.
			structs.EffectPresence{Effect: structs.AddThreePoints, Possible: 0.50, Amount: types.Interval{Min: 0.05, Max: 0.15}},          //71E ->  50% of 03-10 effects.  |  119E ->  50% of 05-17 effects.
			structs.EffectPresence{Effect: structs.AddFourPoints, Possible: 0.30, Amount: types.Interval{Min: 0.03, Max: 0.05}},           //71E ->  30% of 02-03 effects.  |  119E ->  30% of 04-05 effects.
			structs.EffectPresence{Effect: structs.AddFivePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}},           //71E ->  10% of 00-01 effects.  |  119E ->  10% of 01-02 effects.
		)
	}

	// Save on file the new data
	ed.SaveOnFile(utils)

	// Write Reset Message
	if writeMsgData != nil {
		ed.WriteResetMessage(writeMsgData, utils)
	}
}

func (ed *EventsData) EnabledRandomSets(percentage types.Interval, utils types.Utils) error {
	if percentage.Min < 0 {
		return fmt.Errorf("minPercentage must be >= 0")
	} else if percentage.Max > 1 {
		return fmt.Errorf("maxPercentage must be <= 1")
	} else if percentage.Min > percentage.Max {
		return fmt.Errorf("minPercentage must be <= maxPercentage")
	}

	ed.Stats.TotalSetsNum = len(Sets)
	for _, set := range Sets {
		set.Enabled = false
	}

	min, max := int(math.Round(percentage.Min*float64(ed.Stats.TotalSetsNum))), int(math.Round(percentage.Max*float64(ed.Stats.TotalSetsNum)))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	setToActivate := r.Intn(max-min) + min

	for i := 0; i < setToActivate; {
		setIndex := r.Intn(ed.Stats.TotalSetsNum)
		if !Sets[setIndex].Enabled {
			Sets[setIndex].Enabled = true
			ed.Stats.EnabledSetsNum++
			ed.Stats.EnabledSets = append(ed.Stats.EnabledSets, Sets[setIndex].Name)
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
	multiplierEffectsNames, additiveEffectsNames := make([]string, 0), make([]string, 0)
	effectsAmountToApply, effectsToApply, multiplierToApplyNum, additiveToApplyNum := make(map[string]int), make(map[string]*structs.Effect), 0, 0

	for _, effect := range effects {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		if r.Float64() < effect.Possible {
			// Effects will be assigned
			minEventsEffected, maxEventsEffected := int(math.Round(effect.Amount.Min*float64(ed.Stats.EnabledEventsNum))), int(math.Round(effect.Amount.Max*float64(ed.Stats.EnabledEventsNum)))
			if minEventsEffected == maxEventsEffected {
				maxEventsEffected++
			}

			eventsEffected := r.Intn(maxEventsEffected-minEventsEffected) + minEventsEffected
			effectsAmountToApply[effect.Effect.Name] += eventsEffected
			effectsToApply[effect.Effect.Name] = effect.Effect
			if effect.Effect.Key == "*" {
				multiplierEffectsNames = append(multiplierEffectsNames, effect.Effect.Name)
				multiplierToApplyNum += eventsEffected
			} else if effect.Effect.Key == "+" || effect.Effect.Key == "-" {
				additiveEffectsNames = append(additiveEffectsNames, effect.Effect.Name)
				additiveToApplyNum += eventsEffected
			}
		}
	}

	// Check if are applicable all effects calculated
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for additiveToApplyNum > ed.Stats.EnabledEventsNum {
		// Remove a random additive effect
		effectToDecrease := additiveEffectsNames[r.Intn(len(additiveEffectsNames))]
		effectsAmountToApply[effectToDecrease]--
		if effectsAmountToApply[effectToDecrease] == 0 {
			delete(effectsAmountToApply, effectToDecrease)
			additiveEffectsNames = RemoveValue(additiveEffectsNames, effectToDecrease)
		}
	}
	for multiplierToApplyNum > ed.Stats.EnabledEventsNum {
		// Remove a random multiplier effect
		effectToDecrease := multiplierEffectsNames[r.Intn(len(multiplierEffectsNames))]
		effectsAmountToApply[effectToDecrease]--
		if effectsAmountToApply[effectToDecrease] == 0 {
			delete(effectsAmountToApply, effectToDecrease)
			multiplierEffectsNames = RemoveValue(multiplierEffectsNames, effectToDecrease)
		}
	}

	utils.Logger.WithFields(logrus.Fields{
		"toApp": effectsAmountToApply,
	}).Debug("Effects to enable")

	// Apply all effects (multiplier before, additive after)
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, effectName := range multiplierEffectsNames {
		for i := 0; i < effectsAmountToApply[effectName]; {
			eventName := ed.Keys[r.Intn(len(ed.Keys))]
			if ed.Map[eventName].Enabled && len(ed.Map[eventName].Effects) == 0 {
				ed.Map[eventName].AddEffect(effectsToApply[effectName])
				ed.Stats.EnabledEffectsNum++
				ed.Stats.EnabledEffects[effectName]++
				i++
			}
		}
	}
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, effectName := range additiveEffectsNames {
		for i := 0; i < effectsAmountToApply[effectName]; {
			eventName := ed.Keys[r.Intn(len(ed.Keys))]
			if ed.Map[eventName].Enabled {
				if len(ed.Map[eventName].Effects) == 0 {
					//Apply effect if there are no other effects
					ed.Map[eventName].AddEffect(effectsToApply[effectName])
					ed.Stats.EnabledEffectsNum++
					ed.Stats.EnabledEffects[effectName]++
					i++
				} else if len(ed.Map[eventName].Effects) == 1 && ed.Map[eventName].Effects[0].Key == "*" {
					//Apply effect if there is only one effects and it's a multiplier
					ed.Map[eventName].AddEffect(effectsToApply[effectName])
					ed.Stats.EnabledEffectsNum++
					ed.Stats.EnabledEffects[effectName]++
					i++
				}
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
	newS := make([]string, len(s)-1)
	for _, v := range s {
		if v != value {
			newS = append(newS, v)
		}
	}
	return newS
}

func (ed *EventsData) SaveOnFile(utils types.Utils) {
	//Save Sets
	SetsJson = Sets.ToJsonSlice()
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
	// Generate text
	text := "Gli eventi son stati resettati.\nEcco alcune informazioni:\n\n"
	text += fmt.Sprintf("Schemi: %v/%v\nEventi: %v/%v\nPunti ottenibili: %v\n", ed.Stats.EnabledSetsNum, ed.Stats.TotalSetsNum, ed.Stats.EnabledEventsNum, ed.Stats.TotalEventsNum, ed.Stats.EnabledPointsSum)

	text += fmt.Sprintf("\nSchemi Attivi (%v):\n", ed.Stats.EnabledSetsNum)
	for _, setName := range ed.Stats.EnabledSets {
		text += fmt.Sprintf(" | %q\n", setName)
	}

	text += fmt.Sprintf("\nEffetti Attivi (%v):\n", ed.Stats.EnabledEffectsNum)
	for effectName, effectNum := range ed.Stats.EnabledEffects {
		text += fmt.Sprintf(" | %q = %v\n", effectName, effectNum)
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
	PinnedResetMessage = EventsResetMessage{
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
