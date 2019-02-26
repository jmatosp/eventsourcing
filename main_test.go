package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	eProc        EventProcessor
	refact       *Cargo
	kr           *Ship
	sfo, la, yyv *Port
)

func init() {
	eProc = EventProcessor{}
	refact = &Cargo{Name: "Refactoring"}
	kr = NewShip("King Roy")
	sfo = &Port{Name: "San Francisco", Country: "US"}
	la = &Port{Name: "Los Angeles", Country: "US"}
	yyv = &Port{Name: "Vancouver", Country: "CA"}
}

func TestArrivalSetsShipsLocation(t *testing.T) {
	ev := NewArrivalEvent(date(2005, 11, 2), sfo, kr)

	eProc.Process(ev)

	assert.Equal(t, sfo, kr.Port)
}

func TestDeparturePutsShipOutToSea(t *testing.T) {
	eProc.Process(NewArrivalEvent(date(2005, 10, 1), la, kr))
	eProc.Process(NewArrivalEvent(date(2005, 11, 1), sfo, kr))
	eProc.Process(NewDepartureEvent(date(2005, 11, 1), sfo, kr))

	assert.Equal(t, Sea, kr.Port)
}

func TestVisitingCanadaMarksCargo(t *testing.T) {
	eProc.Process(NewLoadEvent(date(2005, 10, 1), refact, kr))
	eProc.Process(NewArrivalEvent(date(2005, 11, 2), yyv, kr))
	eProc.Process(NewDepartureEvent(date(2005, 11, 3), yyv, kr))
	eProc.Process(NewArrivalEvent(date(2005, 11, 4), sfo, kr))
	eProc.Process(NewUnloadEvent(date(2005, 11, 5), refact, kr))

	assert.True(t, refact.HasBeenInCanada)
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
