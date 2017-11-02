package main

import (
	"github.com/martinlindhe/unit"
	"strings"
)

type Unit uint

const (
	NO_UNITS = iota
	UNIT_CELSIUS
	UNIT_FAHRENHEIT
	UNIT_KELVIN

	UNIT_WATT
	UNIT_KILOWATT

	UNIT_WATTHOUR
	UNIT_KILOWATTHOUR

	UNIT_LUX

	UNIT_PERCENTAGE
	UNIT_RELATIVE_HUMIDITY
)

var UnitLookup = map[string]Unit{
	"w":        UNIT_WATT,
	"watt":     UNIT_WATT,
	"kw":       UNIT_KILOWATT,
	"kilowatt": UNIT_KILOWATT,

	"wh":           UNIT_WATTHOUR,
	"watthour":     UNIT_WATTHOUR,
	"kwh":          UNIT_KILOWATTHOUR,
	"kilowatthour": UNIT_KILOWATTHOUR,

	"c":          UNIT_CELSIUS,
	"celsius":    UNIT_CELSIUS,
	"f":          UNIT_FAHRENHEIT,
	"fahrenheit": UNIT_FAHRENHEIT,
	"k":          UNIT_KELVIN,
	"kelvin":     UNIT_KELVIN,

	"lux": UNIT_LUX,
}

func ParseUnit(s string) Unit {
	if unit, found := UnitLookup[strings.ToLower(s)]; found {
		return unit
	} else {
		return NO_UNITS
	}
}

func ConvertFrom(value float64, fromUnit, toUnit Unit) float64 {
	if fromUnit == toUnit {
		return value
	}
	switch fromUnit {

	case UNIT_WATT:
		switch toUnit {
		case UNIT_KILOWATT:
			return (unit.Watt * unit.Power(value)).Kilowatts()
		}
	case UNIT_KILOWATT:
		switch toUnit {
		case UNIT_WATT:
			return (unit.Kilowatt * unit.Power(value)).Watts()
		}

	case UNIT_WATTHOUR:
		switch toUnit {
		case UNIT_KILOWATTHOUR:
			return (unit.WattHour * unit.Energy(value)).KilowattHours()
		}
	case UNIT_KILOWATTHOUR:
		switch toUnit {
		case UNIT_WATTHOUR:
			return (unit.KilowattHour * unit.Energy(value)).WattHours()
		}

	case UNIT_CELSIUS:
		switch toUnit {
		case UNIT_FAHRENHEIT:
			return unit.FromCelsius(value).Fahrenheit()
		case UNIT_KELVIN:
			return unit.FromCelsius(value).Kelvin()
		}
	case UNIT_FAHRENHEIT:
		switch toUnit {
		case UNIT_CELSIUS:
			return unit.FromFahrenheit(value).Celsius()
		case UNIT_KELVIN:
			return unit.FromFahrenheit(value).Kelvin()
		}
	case UNIT_KELVIN:
		switch toUnit {
		case UNIT_FAHRENHEIT:
			return unit.FromKelvin(value).Fahrenheit()
		case UNIT_CELSIUS:
			return unit.FromKelvin(value).Celsius()
		}
	}
	return value
}
