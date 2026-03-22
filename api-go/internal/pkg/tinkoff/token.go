package tinkoff

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Token — подпись Т-Банк (eacq): сортировка ключей, конкатенация значений + Password, SHA256 hex.
// Параметры как в Nest payment.service (String(value) для каждого поля).
func Token(secret string, params map[string]string) string {
	m := make(map[string]string, len(params)+1)
	for k, v := range params {
		m[k] = v
	}
	m["Password"] = secret
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(m[k])
	}
	sum := sha256.Sum256([]byte(b.String()))
	return fmt.Sprintf("%x", sum)
}

// ParamsFromNotificationMap — все скалярные поля из JSON уведомления кроме Token (как filter в Nest).
func ParamsFromNotificationMap(root map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for k, v := range root {
		if k == "Token" || v == nil {
			continue
		}
		if !isScalarJSON(v) {
			continue
		}
		out[k] = ScalarString(v)
	}
	return out
}

func isScalarJSON(v interface{}) bool {
	switch v.(type) {
	case string, bool, json.Number, float64, int, int32, int64, uint, uint32, uint64:
		return true
	default:
		return false
	}
}

// ScalarString — как String() в JS для полей уведомления Т-Банк.
func ScalarString(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case json.Number:
		return x.String()
	case float64:
		if !math.IsNaN(x) && x == math.Trunc(x) && x >= math.MinInt64 && x <= math.MaxInt64 {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case int:
		return strconv.Itoa(x)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	default:
		return fmt.Sprint(x)
	}
}
