package dtekt

import (
	"strings"
)

var everyFloatToString = map[string]string{
	"60.0":  "1m",
	"300.0": "5m",
	"900.0": "15m",
}

// This functions does not validate locations, it assumes user-provided
// locations are validated using terraforms schema validator
func BuildLocationConfig(locations []interface{}) map[string][]string {
	var locationConfig = map[string][]string{}

	for _, loc := range locations {
		split := strings.Split(loc.(string), "-")
		region := split[0]
		location := split[1]

		if _, ok := locationConfig[region]; !ok {
			locationConfig[region] = make([]string, 0)
		}

		locationConfig[region] = append(locationConfig[region], location)
	}

	return locationConfig
}

func FlattenLocationConfig(locationConfig *map[string][]string) []string {
	var locationList = make([]string, 0)

	for region, locations := range *locationConfig {
		for _, locationNumber := range locations {
			locationList = append(locationList, region+"-"+locationNumber)
		}

	}

	return locationList
}

func CovertEvery(every interface{}) interface{} {
	if strings.Contains(every.(string), ".") {
		return everyFloatToString[every.(string)]
	}
	return every
	// switch every.(type) {
	// case float64:
	// 	return everyFloatToString[every.(float64)]
	// default:
	// 	return every
	// }
}

// This function returns true in development env
func DevEnv() bool {
	return true
}
