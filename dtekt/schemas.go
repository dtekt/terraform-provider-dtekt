package dtekt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"strconv"
	"strings"
)

var AlertSchema = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"handlers": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Set:      schema.HashString,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
			},
			"crit": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0,
			},
			"warn": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0,
			},
			"window": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	},
}

func everyValidatorFunc(val interface{}, key string) (warns []string, errs []error) {
	re, _ := regexp.Compile(`^([1-9]{1}[0-9]{0,3})m$`)
	match := re.Find([]byte(val.(string)))
	if match == nil {
		errs = append(errs, fmt.Errorf("%q must be formatted as Xm where X is number of minutes e.g. 15m: %d", key, val.(string)))
		return
	}
	slice := strings.Split(val.(string), "m")
	minutes, _ := strconv.Atoi(slice[0])
	if minutes > 1440 {
		errs = append(errs, fmt.Errorf("%q maximum value is 1440m: %d", key, val.(string)))
	}
	return
}

// Converts human-readable values to float representing number of seconds
// For example "1m" => "60.0" or "5m" => "300.0"
func everyStateFunc(val interface{}) string {
	slice := strings.Split(val.(string), "m")
	minutes, _ := strconv.Atoi(slice[0])
	seconds := float64(minutes * 60)
	return fmt.Sprintf("%.1f", seconds)
}
