package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

// attributesMask
const (
	ManaMask    = 0x3FF
	HealthMask  = 0x3FF
	HealthShift = 10
	HouseFlag   = 1 << 20
	GunFlag     = 1 << 21
	FamilyFlag  = 1 << 22
	TypeMask    = 0x03
	TypeShift   = 23
)

// statsMask
const (
	RespectMask     = 0x0F
	StrengthMask    = 0x0F
	StrengthShift   = 4
	ExperienceMask  = 0x0F
	ExperienceShift = 8
	LevelMask       = 0x0F
	LevelShift      = 12
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		manaValue := uint32(mana) & ManaMask
		person.attributesMask = (person.attributesMask &^ ManaMask) | manaValue
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		healthValue := uint32(health) & HealthMask
		person.attributesMask = (person.attributesMask &^ (HealthMask << HealthShift)) | (healthValue << HealthShift)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		respectValue := uint16(respect) & RespectMask
		person.statsMask = (person.statsMask &^ RespectMask) | respectValue
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		strengthValue := uint16(strength) & StrengthMask
		person.statsMask = (person.statsMask &^ (StrengthMask << StrengthShift)) | (strengthValue << StrengthShift)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		expValue := uint16(experience) & ExperienceMask
		person.statsMask = (person.statsMask &^ (ExperienceMask << ExperienceShift)) | (expValue << ExperienceShift)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		levelValue := uint16(level) & LevelMask
		person.statsMask = (person.statsMask &^ (LevelMask << LevelShift)) | (levelValue << LevelShift)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributesMask |= HouseFlag
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributesMask |= GunFlag
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributesMask |= FamilyFlag
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributesMask = (person.attributesMask &^ (TypeMask << TypeShift)) | (uint32(personType&TypeMask) << TypeShift)
	}
}

type GamePerson struct {
	x, y, z        int32
	gold           uint32
	attributesMask uint32 // 0-9 бит мана, 10-19 здоровье, 20 дом, 21 оружие, 22 семья, 23-24 тип персонажа
	statsMask      uint16 // 0-3 бит уважение, 4-7 сила, 8-11 опыт, 12-15 уровень
	name           [42]byte
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}

	for _, option := range options {
		option(&person)
	}

	return person
}

func (p *GamePerson) Name() string {
	n := 0
	for n < len(p.name) && p.name[n] != 0 {
		n++
	}

	return string(p.name[:n])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.attributesMask & ManaMask)
}

func (p *GamePerson) Health() int {
	return int((p.attributesMask >> HealthShift) & HealthMask)
}

func (p *GamePerson) Respect() int {
	return int(p.statsMask & RespectMask)
}

func (p *GamePerson) Strength() int {
	return int((p.statsMask >> StrengthShift) & StrengthMask)
}

func (p *GamePerson) Experience() int {
	return int((p.statsMask >> ExperienceShift) & ExperienceMask)
}

func (p *GamePerson) Level() int {
	return int((p.statsMask >> LevelShift) & LevelMask)
}

func (p *GamePerson) HasHouse() bool {
	return (p.attributesMask & HouseFlag) != 0
}

func (p *GamePerson) HasGun() bool {
	return (p.attributesMask & GunFlag) != 0
}

func (p *GamePerson) HasFamilty() bool {
	return (p.attributesMask & FamilyFlag) != 0
}

func (p *GamePerson) Type() int {
	return int((p.attributesMask >> TypeShift) & TypeMask)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
