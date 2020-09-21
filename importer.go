package nflrushing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const statsFile = "../../rushing.json"

// ImportStats takes an input file and writes the stats to DB
func ImportStats() error {
	file, err := ioutil.ReadFile(statsFile)
	if err != nil {
		return err
	}
	var ss Stats
	if err = json.Unmarshal([]byte(file), &ss); err != nil {
		return err
	}
	if err = cleanStats(ss); err != nil {
		return err
	}
	for _, stat := range ss {
		if err = InsertRushingStat(stat); err != nil {
			return err
		}
	}
	log.Println("Done importing stats")
	return nil
}

// cleanStats is responsible for parsing and populating type specific fields
func cleanStats(ss Stats) error {
	var err error
	for i, stat := range ss {
		// Yards can be either int or string with comma formatted number
		if ss[i].YardsParsed, err = parseToInt(stat.Yards, ","); err != nil {
			return err
		}
		// Longest can be either int or a string which can be suffixed with "T" sometimes
		// TODO: Populate IsLngTouchdown field based on whether the longest rush was a touchdown
		if ss[i].LongestParsed, err = parseToInt(stat.Longest, "T"); err != nil {
			return err
		}
	}
	return nil
}

// parseToInt is a utility function that converts a variety of formats into go int type
func parseToInt(v interface{}, replaceSet string) (int, error) {
	switch v.(type) {
	case int:
		if val, ok := v.(int); ok {
			return val, nil
		}
		return 0, fmt.Errorf("Failed to parse int")
	case string:
		if val, ok := v.(string); ok {
			return strconv.Atoi(strings.ReplaceAll(val, replaceSet, ""))
		}
		return 0, fmt.Errorf("Failed to parse string")
	default:
		return strconv.Atoi(fmt.Sprintf("%v", v))
	}
	return 0, fmt.Errorf("Unknown type: %v", v)
}
