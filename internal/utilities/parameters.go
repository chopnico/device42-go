package utilities

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/chopnico/structs"
)

// PostParameters will build a struct from methods:post tags
// REVIEW: refactor
// NOTES: i would like to extend the structs packages so that it can build
//
//	dynamic structs based on particular tags. this will be helpful
//	accross this library as device42 requires parameters instead of
//	a post body
func PostParameters(i interface{}) url.Values {
	d := url.Values{}
	z := structs.New(i)

	for _, f := range z.Fields() {
		mtags := strings.Split(f.Tag("methods"), ",")
		for _, tag := range mtags {
			switch tag {
			case "post":
				jtags := strings.Split(f.Tag("json"), ",")
				if !f.IsZero() {
					switch f.Kind() {
					case reflect.Slice:
						if _, ok := f.Value().([]string); ok {
							for _, i := range f.Value().([]string) {
								d.Set(jtags[0], fmt.Sprintf("%v", i))
							}
						}
					default:
						d.Set(jtags[0], fmt.Sprintf("%v", f.Value()))
					}
				}
			}
		}
	}
	return d
}
