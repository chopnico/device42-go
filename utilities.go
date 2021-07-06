package device42

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/chopnico/structs"
)

// REVIEW: refactor
// NOTES: i would like to extend the structs packages so that it can build
// 		  dynamic structs based on particular tags. this will be helpful
// 		  accross this library as device42 requires parameters instead of
// 		  a post body
func parameters(i interface{}) url.Values {
	d := url.Values{}
	z := structs.New(i)

	for _, f := range z.Fields() {
		mtags := strings.Split(f.Tag("methods"), ",")
		for _, tag := range mtags {
			switch tag {
			case "post":
				jtags := strings.Split(f.Tag("json"), ",")
				if !f.IsZero() {
					d.Set(jtags[0], fmt.Sprintf("%v", f.Value()))
				}
			}
		}
	}
	return d
}
