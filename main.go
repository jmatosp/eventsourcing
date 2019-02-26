package main

import (
	"time"
)

type EventProcessor struct {
	log []Event
}

func (ep EventProcessor) Process(e Event) {
	e.Process()
	ep.log = append(ep.log, e)
}

type Event interface {
	Process()
}

type Cargo struct {
	Name            string
	Port            *Port
	Ship            *Ship
	HasBeenInCanada bool
}

func (c *Cargo) HandleArrival(e ArrivalEvent) {
	if Country("CA") == e.Port.Country {
		c.HasBeenInCanada = true
	}
}

func (c *Cargo) HandleLoad(e LoadEvent) {
	c.Port = Sea
	c.Ship = e.Ship
	c.Ship.HandleLoad(e)
}
func (c *Cargo) HandleUnload(e UnloadEvent) {
	c.Port = e.Ship.Port
	c.Ship = NoShip
	e.Ship.HandleUnload(e)
}

type Ship struct {
	Name  string
	Port  *Port
	Cargo []*Cargo
}

var NoShip = &Ship{}

func NewShip(name string) *Ship {
	return &Ship{
		Name:  name,
		Port:  Sea,
		Cargo: []*Cargo{},
	}
}

func (s *Ship) HandleDeparture(e DepartureEvent) {
	s.Port = Sea
}

func (s *Ship) HandleArrival(e ArrivalEvent) {
	s.Port = e.Port
	for _, c := range s.Cargo {
		c.HandleArrival(e)
	}
}

func (s *Ship) HandleLoad(e LoadEvent) {
	s.Cargo = append(s.Cargo, e.Cargo)
}

func (s *Ship) HandleUnload(e UnloadEvent) {
	var leftCargo []*Cargo
	for _, cargo := range s.Cargo {
		if e.Cargo == cargo {
			// skip we no longer carrying this cargo
			continue
		}
		leftCargo = append(leftCargo, cargo)
	}
	s.Cargo = leftCargo
}

type Country string

type Port struct {
	Name    string
	Country Country
}

var Sea = &Port{"Sea", "NA"}

type ArrivalEvent struct {
	Occurred time.Time
	Recorded time.Time
	Port     *Port
	Ship     *Ship
}

func NewArrivalEvent(occurred time.Time, port *Port, ship *Ship) ArrivalEvent {
	return ArrivalEvent{
		Occurred: occurred,
		Recorded: time.Now(),
		Port:     port,
		Ship:     ship,
	}
}

func (e ArrivalEvent) Process() {
	e.Ship.HandleArrival(e)
}

type DepartureEvent struct {
	Occurred time.Time
	Recorded time.Time
	Port     *Port
	Ship     *Ship
}

func NewDepartureEvent(occurred time.Time, port *Port, ship *Ship) DepartureEvent {
	return DepartureEvent{
		Occurred: occurred,
		Recorded: time.Now(),
		Port:     port,
		Ship:     ship,
	}
}

func (e DepartureEvent) Process() {
	e.Ship.HandleDeparture(e)
}

type LoadEvent struct {
	Occurred time.Time
	Recorded time.Time
	Ship     *Ship
	Cargo    *Cargo
}

func NewLoadEvent(occurred time.Time, cargo *Cargo, ship *Ship) LoadEvent {
	return LoadEvent{
		Occurred: occurred,
		Recorded: time.Now(),
		Ship:     ship,
		Cargo:    cargo,
	}
}

func (e LoadEvent) Process() {
	e.Cargo.HandleLoad(e)
}

type UnloadEvent struct {
	Occurred time.Time
	Recorded time.Time
	Ship     *Ship
	Cargo    *Cargo
}

func NewUnloadEvent(occurred time.Time, cargo *Cargo, ship *Ship) UnloadEvent {
	return UnloadEvent{
		Occurred: occurred,
		Recorded: time.Now(),
		Ship:     ship,
		Cargo:    cargo,
	}
}

func (e UnloadEvent) Process() {
	e.Cargo.HandleUnload(e)
}
