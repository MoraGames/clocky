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

func TestUserOverheatingStartsAfterFifteenConsecutiveParticipations(t *testing.T) {
	user := &User{}

	for sequence := int64(1); sequence <= 15; sequence++ {
		if !user.RegisterEventParticipation(sequence) {
			t.Fatalf("participation %d should be allowed", sequence)
		}
	}
	if user.Overheating {
		t.Fatal("user should not overheat after exactly fifteen participations")
	}
	if !user.RegisterEventParticipation(16) {
		t.Fatal("the participation activating overheating should be allowed")
	}
	if !user.Overheating {
		t.Fatal("user should overheat after the sixteenth consecutive participation")
	}
	if user.RegisterEventParticipation(17) {
		t.Fatal("user should not participate while overheated")
	}
}

func TestUserOverheatingEndsAfterFiveMissedEvents(t *testing.T) {
	user := &User{}
	for sequence := int64(1); sequence <= 16; sequence++ {
		user.RegisterEventParticipation(sequence)
	}

	if user.RegisterEventParticipation(22) == false {
		t.Fatal("user should be allowed to participate after missing five events")
	}
	if user.Overheating {
		t.Fatal("user should no longer be overheated")
	}
	if user.OverheatingParticipationStreak != 1 {
		t.Fatalf("expected a new participation streak, got %d", user.OverheatingParticipationStreak)
	}
}
