package main

import (
	"testing"
)

func TestGetStamp(t *testing.T) {
	stamps := "actor(@Chlod) application(Maniphest) author(@Sportzpikachu) herald(H9) herald(H293) herald(H362) monogram(T348212) object-type(TASK) phid(PHID-TASK-l7mhyhysnx7jc6ybjqqq) space(S1) subscriber(@Aklapper) subscriber(@Chlod) subscriber(@Sportzpikachu) subscriber(@Yoshi24517) subtype(feature) tag(#ultraviolet) task-priority(25) task-status(open) task-unassigned() via(web)\n"

	if stamp := GetStamp("actor", stamps, dummyLogger); stamp != "@Chlod" {
		t.Errorf("expected %s for stamp %s, got %s instead", "@Chlod", "actor", stamp)
	}
	if stamp := GetStamp("task-priority", stamps, dummyLogger); stamp != "25" {
		t.Errorf("expected %s for stamp %s, got %s instead", "25", "task-priority", stamp)
	}
	if stamp := GetStamp("tag", stamps, dummyLogger); stamp != "#ultraviolet" {
		t.Errorf("expected %s for stamp %s, got %s instead", "#ultraviolet", "tag", stamp)
	}
	if stamp := GetStamp("herald", stamps, dummyLogger); stamp != "H9" { // match first
		t.Errorf("expected %s for stamp %s, got %s instead", "H9", "herald", stamp)
	}

	if stamp := GetStamp("foo", stamps, dummyLogger); stamp != "" {
		t.Errorf("expected %s to not match, got %s", "foo", stamp)
	}
	if stamp := GetStamp("T348212", stamps, dummyLogger); stamp != "" {
		t.Errorf("expected %s to not match, got %s", "T348212", stamp)
	}
}

func TestMatchesStampFilter(t *testing.T) {
	stampFilters = []string{
		"tag(#foo)",
		"tag(#bar)",
	}
	if stamps := "tag(#foo)"; !MatchesStampFilter(stamps) {
		t.Errorf("expected %s to match", stamps)
	}

	if stamps := "tag(#bar)"; !MatchesStampFilter(stamps) {
		t.Errorf("expected %s to match", stamps)
	}

	if stamps := "tag(#baz) monogram(T123) actor(#foo) tag(#foo)"; !MatchesStampFilter(stamps) {
		t.Errorf("expected %s to match", stamps)
	}

	if stamps := ""; MatchesStampFilter(stamps) {
		t.Errorf("expected %s to not match", stamps)
	}

	if stamps := "lorem ipsum"; MatchesStampFilter(stamps) {
		t.Errorf("expected %s to not match", stamps)
	}

	if stamps := "tag(#foo#bar)"; MatchesStampFilter(stamps) {
		t.Errorf("expected %s to not match", stamps)
	}
}
