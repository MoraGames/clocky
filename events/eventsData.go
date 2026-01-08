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
	"github.com/MoraGames/clockyuwu/pkg/utils"
	"github.com/MoraGames/clockyuwu/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type (
	EventsData struct {
		Map   EventsMap
		Keys  EventsKeys
		Stats EventsStats
		Curr  EventsCurr
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
	EventsCurr struct {
		EnabledSets        map[string]int
		RemainedSetsNum    int
		RemainedEventsNum  int
		EnabledEffects     map[string]int
		RemainedEffectsNum int
		RemainedPointsSum  int
		LastUpdate         time.Time
	}

	EventsResetPinnedMessage struct {
		Exist     bool
		ChatID    int64
		MessageID int
		Text      string
		Entities  []tgbotapi.MessageEntity
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
		Events.ResetEventsCurr()
	}
)

func (ed *EventsData) ResetEventsCurr() {
	ed.Curr = EventsCurr{make(map[string]int), 0, 0, make(map[string]int), 0, 0, time.Now()}

	for _, setName := range ed.Stats.EnabledSets {
		setEventsNum := len(EventsOf(Sets.GetByName(setName).Verify))
		ed.Curr.EnabledSets[setName] = setEventsNum
		ed.Curr.RemainedEventsNum += setEventsNum
	}
	ed.Curr.RemainedSetsNum = len(ed.Curr.EnabledSets)
	ed.Curr.EnabledEffects = ed.Stats.EnabledEffects
	ed.Curr.RemainedEffectsNum = ed.Stats.EnabledEffectsNum
	ed.Curr.RemainedPointsSum = ed.Stats.EnabledPointsSum
}

func NewEventsData(newEffects bool, utils types.Utils) *EventsData {
	ed := &EventsData{
		make(EventsMap),
		make(EventsKeys, 0),
		EventsStats{0, 0, nil, 0, 0, 0, 0, make(map[string]int)},
		EventsCurr{make(map[string]int), 0, 0, make(map[string]int), 0, 0, time.Now()},
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
			structs.EffectPresence{Effect: structs.QuintupleNegativePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}}, // "Mul-5" ->  10% of (95E: 01-02 | 218E: 02-04)
			structs.EffectPresence{Effect: structs.QuadrupleNegativePoints, Possible: 0.40, Amount: types.Interval{Min: 0.02, Max: 0.03}}, // "Mul-4" ->  40% of (95E: 02-03 | 218E: 04-07)
			structs.EffectPresence{Effect: structs.TripleNegativePoints, Possible: 0.70, Amount: types.Interval{Min: 0.03, Max: 0.05}},    // "Mul-3" ->  70% of (95E: 03-05 | 218E: 07-11)
			structs.EffectPresence{Effect: structs.DoubleNegativePoints, Possible: 1.00, Amount: types.Interval{Min: 0.05, Max: 0.10}},    // "Mul-2" -> 100% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.DoublePositivePoints, Possible: 1.00, Amount: types.Interval{Min: 0.08, Max: 0.15}},    // "Mul+2" -> 100% of (95E: 08-14 | 218E: 17-33)
			structs.EffectPresence{Effect: structs.TriplePositivePoints, Possible: 0.75, Amount: types.Interval{Min: 0.05, Max: 0.10}},    // "Mul+3" ->  75% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.QuadruplePositivePoints, Possible: 0.50, Amount: types.Interval{Min: 0.03, Max: 0.05}}, // "Mul+4" ->  50% of (95E: 03-05 | 218E: 07-11)
			structs.EffectPresence{Effect: structs.QuintuplePositivePoints, Possible: 0.25, Amount: types.Interval{Min: 0.02, Max: 0.03}}, // "Mul+5" ->  25% of (95E: 02-03 | 218E: 04-07)
			structs.EffectPresence{Effect: structs.SixtuplePositivePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}},  // "Mul+6" ->  10% of (95E: 01-02 | 218E: 02-04)
			structs.EffectPresence{Effect: structs.SubFourPoints, Possible: 0.25, Amount: types.Interval{Min: 0.02, Max: 0.05}},           // "Sub 4" ->  25% of (95E: 02-05 | 218E: 04-11)
			structs.EffectPresence{Effect: structs.SubThreePoints, Possible: 0.50, Amount: types.Interval{Min: 0.05, Max: 0.10}},          // "Sub 3" ->  50% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.SubTwoPoints, Possible: 1.00, Amount: types.Interval{Min: 0.10, Max: 0.20}},            // "Sub 2" -> 100% of (95E: 10-19 | 218E: 22-44)
			structs.EffectPresence{Effect: structs.AddTwoPoints, Possible: 1.00, Amount: types.Interval{Min: 0.10, Max: 0.20}},            // "Add 2" -> 100% of (95E: 10-19 | 218E: 22-44)
			structs.EffectPresence{Effect: structs.AddThreePoints, Possible: 0.50, Amount: types.Interval{Min: 0.05, Max: 0.10}},          // "Add 3" ->  50% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.AddFourPoints, Possible: 0.25, Amount: types.Interval{Min: 0.02, Max: 0.05}},           // "Add 4" ->  25% of (95E: 02-05 | 218E: 04-11)
			structs.EffectPresence{Effect: structs.AddFivePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}},           // "Add 5" ->  10% of (95E: 01-02 | 218E: 02-04)
		)
	}

	return ed
}

func (ed *EventsData) Reset(newEffects bool, writeMsgData *types.WriteMessageData, utils types.Utils) {
	ed.Stats = EventsStats{0, 0, nil, 0, 0, 0, 0, make(map[string]int)}
	ed.Curr = EventsCurr{make(map[string]int), 0, 0, make(map[string]int), 0, 0, time.Now()}
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
			structs.EffectPresence{Effect: structs.QuintupleNegativePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}}, // "Mul-5" ->  10% of (95E: 01-02 | 218E: 02-04)
			structs.EffectPresence{Effect: structs.QuadrupleNegativePoints, Possible: 0.40, Amount: types.Interval{Min: 0.02, Max: 0.03}}, // "Mul-4" ->  40% of (95E: 02-03 | 218E: 04-07)
			structs.EffectPresence{Effect: structs.TripleNegativePoints, Possible: 0.70, Amount: types.Interval{Min: 0.03, Max: 0.05}},    // "Mul-3" ->  70% of (95E: 03-05 | 218E: 07-11)
			structs.EffectPresence{Effect: structs.DoubleNegativePoints, Possible: 1.00, Amount: types.Interval{Min: 0.05, Max: 0.10}},    // "Mul-2" -> 100% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.DoublePositivePoints, Possible: 1.00, Amount: types.Interval{Min: 0.08, Max: 0.15}},    // "Mul+2" -> 100% of (95E: 08-14 | 218E: 17-33)
			structs.EffectPresence{Effect: structs.TriplePositivePoints, Possible: 0.75, Amount: types.Interval{Min: 0.05, Max: 0.10}},    // "Mul+3" ->  75% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.QuadruplePositivePoints, Possible: 0.50, Amount: types.Interval{Min: 0.03, Max: 0.05}}, // "Mul+4" ->  50% of (95E: 03-05 | 218E: 07-11)
			structs.EffectPresence{Effect: structs.QuintuplePositivePoints, Possible: 0.25, Amount: types.Interval{Min: 0.02, Max: 0.03}}, // "Mul+5" ->  25% of (95E: 02-03 | 218E: 04-07)
			structs.EffectPresence{Effect: structs.SixtuplePositivePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}},  // "Mul+6" ->  10% of (95E: 01-02 | 218E: 02-04)
			structs.EffectPresence{Effect: structs.SubFourPoints, Possible: 0.25, Amount: types.Interval{Min: 0.02, Max: 0.05}},           // "Sub 4" ->  25% of (95E: 02-05 | 218E: 04-11)
			structs.EffectPresence{Effect: structs.SubThreePoints, Possible: 0.50, Amount: types.Interval{Min: 0.05, Max: 0.10}},          // "Sub 3" ->  50% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.SubTwoPoints, Possible: 1.00, Amount: types.Interval{Min: 0.10, Max: 0.20}},            // "Sub 2" -> 100% of (95E: 10-19 | 218E: 22-44)
			structs.EffectPresence{Effect: structs.AddTwoPoints, Possible: 1.00, Amount: types.Interval{Min: 0.10, Max: 0.20}},            // "Add 2" -> 100% of (95E: 10-19 | 218E: 22-44)
			structs.EffectPresence{Effect: structs.AddThreePoints, Possible: 0.50, Amount: types.Interval{Min: 0.05, Max: 0.10}},          // "Add 3" ->  50% of (95E: 05-10 | 218E: 11-22)
			structs.EffectPresence{Effect: structs.AddFourPoints, Possible: 0.25, Amount: types.Interval{Min: 0.02, Max: 0.05}},           // "Add 4" ->  25% of (95E: 02-05 | 218E: 04-11)
			structs.EffectPresence{Effect: structs.AddFivePoints, Possible: 0.10, Amount: types.Interval{Min: 0.01, Max: 0.02}},           // "Add 5" ->  10% of (95E: 01-02 | 218E: 02-04)
		)
	}

	// Update Curr Stats
	ed.ResetEventsCurr()

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

func (ed *EventsData) GenerateResetMessageContent() (string, []tgbotapi.MessageEntity) {
	// Sort the data contained by Stats.EnabledSets and Stats.EnabledEffects
	sortedActiveSets := make([]string, len(ed.Stats.EnabledSets))
	copy(sortedActiveSets, ed.Stats.EnabledSets)
	sort.Slice(sortedActiveSets, func(i, j int) bool {
		return sortedActiveSets[i] < sortedActiveSets[j]
	})

	sortedEnabledEffects := make([]EffectPresence, 0, len(ed.Curr.EnabledEffects))
	for effectName, effectNum := range ed.Curr.EnabledEffects {
		sortedEnabledEffects = append(sortedEnabledEffects, EffectPresence{
			Name:   effectName,
			Amount: effectNum,
		})
	}
	sort.Slice(sortedEnabledEffects, func(i, j int) bool {
		return sortedEnabledEffects[i].Name < sortedEnabledEffects[j].Name
	})

	// Generate text
	rawText := "__**Gli eventi son stati resettati.**__\nEcco alcune info sulla giornata restante:\n\n"
	rawText += fmt.Sprintf("Set: %v/%v\nEventi: %v/%v\nEffetti: %v/%v\nPunti: %v/%v\n", ed.Curr.RemainedSetsNum, ed.Stats.EnabledSetsNum, ed.Curr.RemainedEventsNum, ed.Stats.EnabledEventsNum, ed.Curr.RemainedEffectsNum, ed.Stats.EnabledEffectsNum, ed.Curr.RemainedPointsSum, ed.Stats.EnabledPointsSum)

	rawText += "\n**Set e Eventi attivi:**\n"
	for _, setName := range sortedActiveSets {
		if ed.Curr.EnabledSets[setName] == 0 {
			rawText += fmt.Sprintf("| ~~%s -> %v~~\n", setName, ed.Curr.EnabledSets[setName])
		} else {
			rawText += fmt.Sprintf("| %s -> %v\n", setName, ed.Curr.EnabledSets[setName])
		}
	}

	rawText += "\n**Effetti attivi:**\n"
	for _, effect := range sortedEnabledEffects {
		if effect.Amount == 0 {
			rawText += fmt.Sprintf("| ~~%s = %v~~\n", effect.Name, effect.Amount)
		} else {
			rawText += fmt.Sprintf("| %s = %v\n", effect.Name, effect.Amount)
		}
	}

	rawText += "\nBuona fortuna!"

	entities, text := utils.ParseToEntities(rawText, nil)
	return text, entities
}

func (ed *EventsData) WriteResetMessage(writeMsgData *types.WriteMessageData, utils types.Utils) {
	// Generate message
	text, entities := ed.GenerateResetMessageContent()
	message := tgbotapi.NewMessage(writeMsgData.ChatID, text)
	if entities != nil {
		message.Entities = entities
	}

	// Send message
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

func (ed *EventsData) OverwriteResetMessage(msgId int, writeMsgData *types.WriteMessageData, utils types.Utils) {
	// Edit message
	text, entities := ed.GenerateResetMessageContent()
	message := tgbotapi.NewEditMessageText(writeMsgData.ChatID, msgId, text)
	if entities != nil {
		message.Entities = entities
	}

	// Update message
	msg, err := writeMsgData.Bot.Send(message)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": msg,
		}).Error("Error while editing message")
	}
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
		msgToPin.Text,
		msgToPin.Entities,
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

func FastforwardUpdateDailyCounters(utilsVar types.Utils) {
	t, now := Events.Curr.LastUpdate.AddDate(0, 0, 1), time.Now()
	if t.IsZero() {
		utilsVar.Logger.Warn("LastUpdate is zeroed, FastforwardUpdateDailyCounters will not be executed")
		return
	}

	for ; t.Before(now) || now.Second() == 59; t, now = t.Add(time.Minute), time.Now() {
		// Check if the current time is a valid enabled event time (and force skip at 23:59)
		if t.Hour() == 23 && t.Minute() == 59 {
			return
		}
		event, exists := Events.Map[fmt.Sprintf("%d%d:%d%d", t.Hour()/10, t.Hour()%10, t.Minute()/10, t.Minute()%10)]
		if !exists {
			return
		}
		if !event.Enabled {
			return
		}

		// Update the events structures
		enablingSets := CalculateEnablingSets(t)
		for _, setName := range enablingSets {
			Events.Curr.EnabledSets[setName]--
			Events.Curr.RemainedEventsNum--
			if Events.Curr.EnabledSets[setName] == 0 {
				Events.Curr.RemainedSetsNum--
			}
		}
		for _, effect := range event.Effects {
			Events.Curr.EnabledEffects[effect.Name]--
			Events.Curr.RemainedEffectsNum--
		}
		Events.Curr.RemainedPointsSum -= event.CalculateTotalPoints()
		Events.Curr.LastUpdate = t

		// Update the message data
		// This is commented since in this function we have bot/chat data. Since this fastforward is used only at startup, the first cronjob iteration will update the message with the correct data.
		//UpdateEventsDataMessage(&types.WriteMessageData{Bot: utilsVar.Bot, ChatID: utilsVar.ChatID, ReplyMessageID: -1}, utilsVar)

		// Save the Events data
		Events.SaveOnFile(utilsVar)
	}
}

func UpdateEventsDataMessage(writeMsgData *types.WriteMessageData, utilsVar types.Utils) {
	// Read PinnedResetMessage
	pinnedMessageFile, err := os.ReadFile("files/pinnedMessage.json")
	if err != nil {
		utilsVar.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while reading Events data")
		return
	}
	var PinnedResetMessage EventsResetPinnedMessage
	err = json.Unmarshal(pinnedMessageFile, &PinnedResetMessage)
	if err != nil {
		utilsVar.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while unmarshalling Events data")
		return
	}

	// Update the message
	Events.OverwriteResetMessage(PinnedResetMessage.MessageID, writeMsgData, utilsVar)
}
