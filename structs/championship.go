package structs

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/MoraGames/clockyuwu/pkg/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

var (
	PinnedChampionshipResetMessage ChampionshipResetPinnedMessage
)

type ChampionshipResetPinnedMessage struct {
	Exist     bool
	ChatID    int64
	MessageID int
}
type Championship struct {
	Name         string
	StartDate    time.Time
	Duration     time.Duration
	Status       string
	FinalRanking []Rank
}

func NewChampionship(name string, startDate time.Time, duration time.Duration, status string, finalRanking []Rank) *Championship {
	return &Championship{name, startDate, duration, status, finalRanking}
}

func NewEndedChampionship(name string, startDate time.Time, duration time.Duration, finalRanking []Rank) *Championship {
	return &Championship{name, startDate, duration, "ended", finalRanking}
}

func CreateChampionship(name string, startDate time.Time, duration time.Duration) *Championship {
	curTime := time.Now()
	var status string
	if startDate.After(curTime) {
		status = "upcoming"
	} else if startDate.Before(curTime) && startDate.Add(duration).Before(curTime) {
		status = "ended"
	} else {
		status = "ongoing"
	}
	return &Championship{name, startDate, duration, status, nil}
}

func (c *Championship) End(finalRanking []Rank) {
	c.FinalRanking = finalRanking
	c.Status = "ended"
}

func (c *Championship) Reset(ranking []Rank, writeMsgData *types.WriteMessageData, utils types.Utils) {
	// Read from file (previous championship data)
	prevChampionship := ReadFromFile(utils)

	// Save on file the new data
	c.End(ranking)
	c.SaveOnFile(utils)

	// Write Reset Message
	if writeMsgData != nil {
		WriteChampionshipResetMessage(c, prevChampionship, writeMsgData, utils)
	}
}

func ReadFromFile(utils types.Utils) *Championship {
	file, err := os.ReadFile("files/championship.json")
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while reading Championship data from file")
		return nil
	}
	var championship Championship
	err = json.Unmarshal(file, &championship)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while unmarshalling Championship data from file")
		return nil
	}
	return &championship
}

func (c *Championship) SaveOnFile(utils types.Utils) {
	championshipFile, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while marshalling Championship data")
	}
	err = os.WriteFile("files/championship.json", championshipFile, 0644)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while writing Championship data")
	}
}

func WriteChampionshipResetMessage(curChamp, prevChamp *Championship, writeMsgData *types.WriteMessageData, utils types.Utils) {
	var improvedPlayers, wrosenedPlayers, newPlayers []Rank

	var text string
	if curChamp.FinalRanking == nil || len(curChamp.FinalRanking) == 0 {
		text = "Il campionato è giunto al termine, ma nessun giocatore ha partecipato questa volta.\n\nForza giocatori! Vogliamo vedere più duelli nel campionato che sta per iniziare!"
	} else {
		for curPosition, curRank := range curChamp.FinalRanking {
			var founded bool = false
			if prevChamp == nil {
				newPlayers = append(newPlayers, curRank)
				continue
			}
			for prevPosition, prevRank := range prevChamp.FinalRanking {
				if curRank.UserTelegramID == prevRank.UserTelegramID {
					founded = true
					if curPosition < prevPosition {
						improvedPlayers = append(improvedPlayers, curRank)
					} else if curPosition > prevPosition {
						wrosenedPlayers = append(wrosenedPlayers, curRank)
					}
				}
			}
			if !founded {
				newPlayers = append(newPlayers, curRank)
			}
		}

		// Generate text
		text = "Il campionato è giunto al termine ed un nuovo Clocky Champion è stato incoronato!\n\n"
		for i, rank := range curChamp.FinalRanking[:int(math.Min(3, float64(len(curChamp.FinalRanking))))] {
			text += fmt.Sprintf("%d°: %s con %d punti e %d partecipazioni\n", i+1, rank.Username, rank.Points, rank.Partecipations)
		}
		if len(newPlayers) > 0 {
			text += "\nDiamo inoltre il benvenuto a "
			for _, player := range newPlayers {
				text += fmt.Sprintf("%s, ", player.Username)
			}
			text += "che hanno deciso di scompigliare i piani dei pù esperti.\n"
		}

		if len(improvedPlayers) > 0 && len(wrosenedPlayers) > 0 {
			text += "\nInfine, i migliori giocatori sono stati "
			for _, player := range improvedPlayers {
				text += fmt.Sprintf("%s, ", player.Username)
			}
			text += "che sono riusciti a migliorare la loro posizione in classifica a discapito di"
			for i, player := range wrosenedPlayers {
				if i == len(wrosenedPlayers)-1 {
					text += fmt.Sprintf("%s.\n", player.Username)
				} else {
					text += fmt.Sprintf("%s, ", player.Username)
				}
			}
		} else if len(improvedPlayers) > 0 && len(wrosenedPlayers) == 0 {
			text += "\nInfine, grazie a qualche abbandono, "
			for _, player := range improvedPlayers {
				text += fmt.Sprintf("%s, ", player.Username)
			}
			text += "sono riusciti a migliorare la loro posizione in classifica.\n"
		} else if len(improvedPlayers) == 0 && len(wrosenedPlayers) > 0 {
			text += "\nNon è un caso che siano riusciti a battere "
			for i, player := range wrosenedPlayers {
				if i == len(wrosenedPlayers)-1 {
					text += fmt.Sprintf("%s.\n", player.Username)
				} else {
					text += fmt.Sprintf("%s, ", player.Username)
				}
			}
		}

		text += "\nMa bando alle ciance, i preparativi per il prossimo campionato sono già completati.\nConcorrenti preparatevi, è l'ora di ricominciare a fare punti!"
	}

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
	UpdatePinnedChampionshipMessage(writeMsgData, utils, msg)
}

func UpdatePinnedChampionshipMessage(writeMsgData *types.WriteMessageData, utils types.Utils, msgToPin tgbotapi.Message) {
	// Unpin the old reset message if exists
	if PinnedChampionshipResetMessage.Exist {
		msg, err := writeMsgData.Bot.Send(tgbotapi.UnpinChatMessageConfig{
			ChatID:    PinnedChampionshipResetMessage.ChatID,
			MessageID: PinnedChampionshipResetMessage.MessageID,
		})
		if err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"err": err,
				"msg": msg,
			}).Error("Error while unpinning message")
		}
	}

	// Update the pinned reset message
	PinnedChampionshipResetMessage = ChampionshipResetPinnedMessage{
		true,
		msgToPin.Chat.ID,
		msgToPin.MessageID,
	}

	// Save PinnedResetMessage
	pinnedMessageFile, err := json.MarshalIndent(PinnedChampionshipResetMessage, "", " ")
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while marshalling Championship data")
	}
	err = os.WriteFile("files/pinnedChampionshipMessage.json", pinnedMessageFile, 0644)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while writing Championship data")
	}

	// Pin the new reset message if exists
	if PinnedChampionshipResetMessage.Exist {
		msg, err := writeMsgData.Bot.Send(tgbotapi.PinChatMessageConfig{
			ChatID:              PinnedChampionshipResetMessage.ChatID,
			MessageID:           PinnedChampionshipResetMessage.MessageID,
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
