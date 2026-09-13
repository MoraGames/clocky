package structs

import "testing"

func hasEffect(user *User, effectName string) bool {
	for _, effect := range user.Effects {
		if effect.Name == effectName {
			return true
		}
	}
	return false
}

func ensureHasEffects(t *testing.T, user *User, effects ...*Effect) {
	for _, effect := range effects {
		if !hasEffect(user, effect.Name) {
			t.Errorf("User should have effect %q", effect.Name)
		}
	}
}

func ensureNotHasEffects(t *testing.T, user *User, effects ...*Effect) {
	for _, effect := range effects {
		if hasEffect(user, effect.Name) {
			t.Errorf("User should not have effect %q", effect.Name)
		}
	}
}

var testEffect1 = Effect{
	Name: "Test1",
}
var testEffect2 = Effect{
	Name: "Test2",
}
var testEffect3 = Effect{
	Name: "Test3",
}

func Test_RemoveUserEffect(t *testing.T) {
	user := User{
		Effects: []*Effect{
			&testEffect1,
			&testEffect2,
			&testEffect3,
		},
	}

	user.RemoveEffect(&testEffect2)
	ensureHasEffects(t, &user, &testEffect1, &testEffect3)
	ensureNotHasEffects(t, &user, &testEffect2)
}

func Test_RemoveUserEffect_SameEffectMultipleTimes(t *testing.T) {
	user := User{
		Effects: []*Effect{
			&testEffect1,
			&testEffect2,
			&testEffect2,
			&testEffect2,
			&testEffect2,
			&testEffect2,
			&testEffect3,
		},
	}

	user.RemoveEffect(&testEffect2)
	ensureHasEffects(t, &user, &testEffect1, &testEffect3)
	ensureNotHasEffects(t, &user, &testEffect2)
}

func TestUserOverheatingUsesBucketThresholds(t *testing.T) {
	user := &User{}
	user.OverheatingPoints = 90
	if user.IsOverheating() {
		t.Fatal("user should not overheat at 90 points")
	}
	user.OverheatingPoints = 91
	if !user.IsOverheating() {
		t.Fatal("user should overheat above 90 points")
	}
	if !user.RegisterEventParticipation(1) {
		t.Fatal("overheating should not block participation")
	}
	user.OverheatingPoints = 40
	if !user.IsOverheating() {
		t.Fatal("user should remain overheated at 40 points")
	}
	user.OverheatingPoints = 39
	if user.IsOverheating() {
		t.Fatal("user should leave overheating below 40 points")
	}
}

func TestUserOverheatingUpdatesStreaksAndDecay(t *testing.T) {
	user := &User{}
	user.RegisterEventParticipation(1)
	user.RegisterEventWin()
	user.RegisterEventParticipation(4)

	if user.ChampionshipParticipationStreak != 1 || user.ChampionshipAbsenceStreak != 0 || user.ChampionshipWinStreak != 0 {
		t.Fatalf("unexpected streaks: participation=%d absence=%d win=%d", user.ChampionshipParticipationStreak, user.ChampionshipAbsenceStreak, user.ChampionshipWinStreak)
	}
	if user.OverheatingPoints != -1 {
		t.Fatalf("expected two absence decays and one win gain, got %d points", user.OverheatingPoints)
	}

	user.RegisterEventWin()
	if user.ChampionshipWinStreak != 1 || user.OverheatingPoints != 2 {
		t.Fatalf("unexpected win streak or points: win=%d points=%d", user.ChampionshipWinStreak, user.OverheatingPoints)
	}
	user.RegisterEventParticipation(5)
	if user.ChampionshipAbsenceStreak != 0 || user.ChampionshipParticipationStreak != 2 {
		t.Fatalf("unexpected streak reset: participation=%d absence=%d", user.ChampionshipParticipationStreak, user.ChampionshipAbsenceStreak)
	}
}

func TestUserResetChampionshipOverheating(t *testing.T) {
	user := &User{
		ChampionshipParticipationStreak: 3,
		ChampionshipWinStreak:           4,
		ChampionshipAbsenceStreak:       5,
		OverheatingPoints:               91,
		Overheating:                     true,
	}

	user.ResetChampionshipOverheating()

	if user.ChampionshipParticipationStreak != 0 || user.ChampionshipWinStreak != 0 || user.ChampionshipAbsenceStreak != 0 || user.OverheatingPoints != 0 || user.Overheating {
		t.Fatalf("championship overheating statistics were not reset: %+v", user)
	}
}
