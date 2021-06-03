package orm

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Condition 查询条件
type Condition struct {
	Field     string      `json:"field"`
	Value     interface{} `json:"value"`
	MatchType OprType     `json:"match_type"`
}

// OprType 操作类型
type OprType string

const (
	// Accurate 精确匹配
	Accurate OprType = "accurate"
	// Fuzzy 模糊匹配
	Fuzzy OprType = "fuzzy"
	// Not 不匹配
	Not OprType = "not"
	// Bwt 位于两者之间（含）
	Bwt OprType = "between"
	// In 在某个范围内
	In OprType = "in"
	// NotIn 不在某个范围内
	NotIn OprType = "notin"
	// Larger 大于
	Larger OprType = "larger"
	// LargerOrEqual 大于或等于
	LargerOrEqual OprType = "le"
	// Smaller 小于
	Smaller OprType = "smaller"
	// SmallerOrEqual 小于或等于
	SmallerOrEqual OprType = "se"
)

// BuildCondition 根据condition生成where条件
func BuildCondition(cons []Condition) (string, []interface{}, error) {
	sql := ""
	pattern := "^[a-z0-9A-Z_\\.]+$"
	var params []interface{}

	flag := "`"

	// 添加字段条件
	for _, con := range cons {
		if ok, _ := regexp.MatchString(pattern, con.Field); !ok {
			return "", params, errors.New("列名错误:" + con.Field)
		}
		if con.MatchType == Fuzzy {
			sql += fmt.Sprintf(" and %s%s%s like ? ", flag,con.Field,flag)
			params = append(params, "%"+fmt.Sprintf("%v", con.Value)+"%")
		} else if con.MatchType == Accurate || con.MatchType == "" {
			sql += fmt.Sprintf(" and %s%s%s = ? ", flag,con.Field,flag)
			params = append(params, con.Value)
		} else if con.MatchType == Not {
			sql += fmt.Sprintf(" and %s%s%s != ? ", flag,con.Field,flag)
			params = append(params, con.Value)
		} else if con.MatchType == In {
			sql += fmt.Sprintf(" and %s%s%s in (?) ", flag,con.Field,flag)
			arr, ok := con.Value.([]interface{})
			if !ok {
				panic("In类型值转换错误,必须为slice")
			}
			params = append(params, arr)
		} else if con.MatchType == NotIn {
			sql += fmt.Sprintf(" and %s%s%s not in (?) ", flag,con.Field,flag)
			arr, ok := con.Value.([]interface{})
			if !ok {
				panic("In类型值转换错误,必须为slice")
			}
			params = append(params, arr)
		} else if con.MatchType == Larger {
			sql += fmt.Sprintf(" and %s%s%s > (?) ", flag,con.Field,flag)
			params = append(params, con.Value)
		} else if con.MatchType == LargerOrEqual {
			sql += fmt.Sprintf(" and %s%s%s >= (?) ", flag,con.Field,flag)
			params = append(params, con.Value)
		} else if con.MatchType == Smaller {
			sql += fmt.Sprintf(" and %s%s%s < (?) ", flag,con.Field,flag)
			params = append(params, con.Value)
		} else if con.MatchType == SmallerOrEqual {
			sql += fmt.Sprintf(" and %s%s%s <= (?) ", flag,con.Field,flag)
			params = append(params, con.Value)
		} else if con.MatchType == Bwt {
			sql += fmt.Sprintf(" and ( %s%s%s between ? and ? ) ", flag,con.Field,flag)
			condition := fmt.Sprintf("%v", con.Value)

			values := strings.Split(condition, ";")
			var start, end string
			if len(values) >= 2 {
				start = values[0]
				end = values[1]
			} else {
				start = fmt.Sprintf("%v", con.Value)
				end = ""
			}
			params = append(params, start)
			params = append(params, end)
		}
	}

	sql = strings.TrimLeft(sql, " and ")
	return sql, params, nil
}