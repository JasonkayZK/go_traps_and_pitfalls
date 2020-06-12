package main

import (
	"fmt"
	"time"
)

const (
	timeLayout = `2006-01-02 15:04:00`
)

// timezone timestamp offset
var timezoneConst = map[string]time.Duration{
	"UTC-12 IDLW":    -12 * 3600,
	"UTC-11 MIT":     -11 * 3600,
	"UTC-10 HST":     -10 * 3600,
	"UTC-9:30 MSIT":  -9.5 * 3600,
	"UTC-9 AKST":     -9 * 3600,
	"UTC-8 PST":      -8 * 3600,
	"UTC-7 MST":      -7 * 3600,
	"UTC-6 CST":      -6 * 3600,
	"UTC-5 EST":      -5 * 3600,
	"UTC-4 AST":      -4 * 3600,
	"UTC-3:30 NST":   -3.5 * 3600,
	"UTC-3 SAT":      -3 * 3600,
	"UTC-2":          -2 * 3600,
	"UTC-1 CVT":      -1 * 3600,
	"UTC":            0,
	"UTC+1 CET":      3600,
	"UTC+2 EET":      2 * 3600,
	"UTC+3 MSK":      3 * 3600,
	"UTC+3:30 IRT":   3.5 * 3600,
	"UTC+4 META":     4 * 3600,
	"UTC+4:30 AFT":   4.5 * 3600,
	"UTC+5 METB":     5 * 3600,
	"UTC+5:30 IDT":   5.5 * 3600,
	"UTC+6 BHT":      6 * 3600,
	"UTC+6:30 MRT":   6.5 * 3600,
	"UTC+7 IST":      7 * 3600,
	"UTC+8 EAT":      8 * 3600,
	"UTC+9 FET":      9 * 3600,
	"UTC+9:30 ACST":  9.5 * 3600,
	"UTC+10 AEST":    10 * 3600,
	"UTC+10:30 FAST": 10.5 * 3600,
	"UTC+11 VIT":     11 * 3600,
	"UTC+11:30 NFT":  11.5 * 3600,
	"UTC+12 PSTB":    12 * 3600,
	"UTC+12:45 CIT":  12.75 * 3600,
	"UTC+13 PSTC":    13 * 3600,
	"UTC+14 PSTD":    14 * 3600,
}

func main() {
	fmt.Println(time.Now().UTC().Add(timezoneConst["UTC+14 PSTD"] * time.Second))
	fmt.Println(time.Now().UTC().Add(timezoneConst["UTC+8 EAT"] * time.Second))
}