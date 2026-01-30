package events

import (
	"time"

	"github.com/MoraGames/clockyuwu/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserTrackersMap map[int64]*UserTracker
type UserTracker struct {
	User            *tgbotapi.User
	DailyActivities []structs.DailyActivity
}

func NewUserTracker(telegramUser *tgbotapi.User) *UserTracker {
	return &UserTracker{
		User:            telegramUser,
		DailyActivities: make([]structs.DailyActivity, 0),
	}
}

func (ut *UserTracker) InitializeNewDailyActivity() {
	eventsParticipation := make([]structs.EventParticipation, 0, len(Events.Keys))
	for _, val := range Events.Keys {
		eventsParticipation = append(eventsParticipation, structs.EventParticipation{Event: val, Participated: false, Won: false})
	}

	ut.DailyActivities = append(ut.DailyActivities, structs.DailyActivity{
		Date:                 time.Now().Format("2006-01-02"),
		Activities:           make([]structs.Activity, 0),
		EventsParticipations: eventsParticipation,
	})
}

func (ut *UserTracker) PushActivity(activity structs.Activity) {
	// If there is no daily activity for today, create it before pushing the activity in
	if len(ut.DailyActivities) == 0 || ut.DailyActivities[len(ut.DailyActivities)-1].Date != time.Now().Format("2006-01-02") {
		ut.InitializeNewDailyActivity()
	}

	// Push the activity in the last daily activity
	ut.DailyActivities[len(ut.DailyActivities)-1].Activities = append(ut.DailyActivities[len(ut.DailyActivities)-1].Activities, activity)

	// If the activity is related to an event, update the event participation accordingly
	if activity.Type == structs.EventParticipationActivity || activity.Type == structs.EventWinActivity {
		for i, eventParticipation := range ut.DailyActivities[len(ut.DailyActivities)-1].EventsParticipations {
			if eventParticipation.Event == activity.TelegramTime.Format("15:04") {
				ut.DailyActivities[len(ut.DailyActivities)-1].EventsParticipations[i].Participated = true
				if activity.Type == structs.EventWinActivity {
					ut.DailyActivities[len(ut.DailyActivities)-1].EventsParticipations[i].Won = true
				}
				break
			}
		}
	}
}
