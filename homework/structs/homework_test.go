package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
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
		manaValue := uint32(mana) & 0x3FF
		person.attributesMask = (person.attributesMask &^ 0x3FF) | manaValue
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		healthValue := uint32(health) & 0x3FF
		person.attributesMask = (person.attributesMask &^ (0x3FF << 10)) | (healthValue << 10)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		respectValue := uint16(respect) & 0x0F
		person.statsMask = (person.statsMask &^ 0x0F) | respectValue
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		strengthValue := uint16(strength) & 0x0F
		person.statsMask = (person.statsMask &^ (0x0F << 4)) | (strengthValue << 4)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		expValue := uint16(experience) & 0x0F
		person.statsMask = (person.statsMask &^ (0x0F << 8)) | (expValue << 8)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		levelValue := uint16(level) & 0x0F
		person.statsMask = (person.statsMask &^ (0x0F << 12)) | (levelValue << 12)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributesMask |= 1 << 20
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributesMask |= 1 << 21
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributesMask |= 1 << 22
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributesMask = (person.attributesMask &^ (0x03 << 23)) | (uint32(personType&0x03) << 23)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

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
	return int(p.attributesMask & 0x3FF)

}

func (p *GamePerson) Health() int {
	return int((p.attributesMask >> 10) & 0x3FF)
}

func (p *GamePerson) Respect() int {
	return int(p.statsMask & 0x0F)
}

func (p *GamePerson) Strength() int {
	return int((p.statsMask >> 4) & 0x0F)
}

func (p *GamePerson) Experience() int {
	return int((p.statsMask >> 8) & 0x0F)
}

func (p *GamePerson) Level() int {
	return int((p.statsMask >> 12) & 0x0F)
}

func (p *GamePerson) HasHouse() bool {
	return (p.attributesMask & (1 << 20)) != 0
}

func (p *GamePerson) HasGun() bool {
	return (p.attributesMask & (1 << 21)) != 0
}

func (p *GamePerson) HasFamilty() bool {
	return (p.attributesMask & (1 << 22)) != 0
}

func (p *GamePerson) Type() int {
	return int((p.attributesMask >> 23) & 0x03)
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
