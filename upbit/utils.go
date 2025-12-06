package upbit

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// structToMap 구조체를 map[string]string 으로 변환
func structToMap(item interface{}) map[string]string {
	res := map[string]string{}
	if item == nil {
		return res
	}
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i)
		val := v.Field(i)

		// 값이 비어있으면 스킵
		if val.IsZero() {
			continue
		}

		// json 태그 값을 키로 사용
		tag := field.Tag.Get("json")
		key := strings.Split(tag, ",")[0]
		if key == "" {
			key = field.Name
		}

		// 값 변환
		var strVal string
		switch val.Kind() {
		case reflect.Int, reflect.Int64:
			strVal = strconv.FormatInt(val.Int(), 10)
		case reflect.Float64:
			strVal = strconv.FormatFloat(val.Float(), 'f', -1, 64)
		case reflect.Bool:
			strVal = strconv.FormatBool(val.Bool())
		default:
			strVal = fmt.Sprintf("%v", val.Interface())
		}
		res[key] = strVal
	}
	return res
}

// generateQueryString: 맵을 정렬된 쿼리 스트링으로 변환
func generateQueryString(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Upbit는 파라미터 정렬을 권장함

	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	return strings.Join(parts, "&")
}

// generateToken: JWT 토큰 생성
func (c *Client) generateToken(queryParams map[string]string) (string, error) {
	claims := jwt.MapClaims{
		"access_key": c.AccessKey,
		"nonce":      uuid.New().String(),
	}

	if len(queryParams) > 0 {
		queryString := generateQueryString(queryParams)
		hash := sha512.Sum512([]byte(queryString))
		hashHex := hex.EncodeToString(hash[:])

		claims["query_hash"] = hashHex
		claims["query_hash_alg"] = "SHA512"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.SecretKey))
}

// doRequest: 실제 요청 수행
func (c *Client) doRequest(method, endpoint string, params interface{}, result interface{}) error {
	paramMap := structToMap(params)

	var body io.Reader
	urlString := BaseURL + endpoint

	// GET/DELETE는 쿼리 스트링에 파라미터 추가
	if method == http.MethodGet || method == http.MethodDelete {
		if len(paramMap) > 0 {
			q := url.Values{}
			for k, v := range paramMap {
				q.Add(k, v)
			}
			urlString += "?" + q.Encode()
		}
	} else {
		// POST는 JSON Body 사용
		if params != nil {
			jsonBytes, err := json.Marshal(paramMap)
			if err != nil {
				return err
			}
			body = bytes.NewBuffer(jsonBytes)
		}
	}

	req, err := http.NewRequest(method, urlString, body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	// 인증 토큰 추가 (Auth가 필요한 경우)
	// paramMap은 Hash 생성을 위해 사용됨
	token, err := c.generateToken(paramMap)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Result가 nil이 아니고 포인터일 때만 언마샬링
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("json unmarshal error: %w, body: %s", err, string(respBody))
		}
	}

	return nil
}

// doNonAuthRequest: 인증이 필요 없는 요청 (시세 조회 등)
func (c *Client) doNonAuthRequest(endpoint string, params interface{}, result interface{}) error {
	paramMap := structToMap(params)
	urlString := BaseURL + endpoint

	if len(paramMap) > 0 {
		q := url.Values{}
		for k, v := range paramMap {
			q.Add(k, v)
		}
		urlString += "?" + q.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, urlString, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}
