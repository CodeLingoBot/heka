/***** BEGIN LICENSE BLOCK *****
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.
#
# The Initial Developer of the Original Code is the Mozilla Foundation.
# Portions created by the Initial Developer are Copyright (C) 2012
# the Initial Developer. All Rights Reserved.
#
# Contributor(s):
#   Victor Ng (vng@mozilla.com)
#
# ***** END LICENSE BLOCK *****/

package pipeline

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type JsonPath struct {
	json_data interface{}
	json_text string
}

var json_re = regexp.MustCompile("^([^0-9\\s\\[][^\\s\\[]*)?(\\[[0-9]+\\])?$")

func (j *JsonPath) SetJsonText(json_text string) (err error) {
	j.json_text = json_text
	err = json.Unmarshal([]byte(json_text), &j.json_data)
	return
}

func (j *JsonPath) Find(jp string) (result interface{}, err error) {
	var ok bool

	if jp == "" {
		return result, errors.New("invalid path")
	}

	// Need to grab a pointer to the top of the data structure
	v := j.json_data

	for _, token := range strings.Split(jp, "/") {
		sl := json_re.FindAllStringSubmatch(token, -1)
		if len(sl) == 0 {
			return result, errors.New("invalid path")
		}
		ss := sl[0]
		if ss[1] != "" {
			v, ok = v.(map[string]interface{})[ss[1]]
			if !ok {
				return result, errors.New("invalid path")
			}
		}
		if ss[2] != "" {
			i, err := strconv.Atoi(ss[2][1 : len(ss[2])-1])
			if err != nil {
				return result, errors.New("invalid path")
			}
			v_arr := v.([]interface{})
			if i < 0 || i >= len(v_arr) {
				return result, errors.New(fmt.Sprintf("array out of bounds jsonpath:[%s]", jp))
			}
			v = v_arr[i]
		}
	}

	//rv := reflect.ValueOf(v)
	//rv_kind := rv.Kind()
	result = v
	return result, nil
}
