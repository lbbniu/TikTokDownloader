package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/tjfoc/gmsm/sm3"
)

type ABogus struct {
	chunk       []byte
	size        int
	reg         []uint32
	uaCode      []byte
	browser     string
	browserLen  int
	browserCode []int
}

var (
	filter    = regexp.MustCompile(`%([0-9A-F]{2})`)
	arguments = []int{0, 1, 14}
	uaKey     = "\u0000\u0001\u000e"
	endString = "cus"
	version   = []int{1, 0, 1, 5}
	browser   = "1536|742|1536|864|0|0|0|0|1536|864|1536|864|1536|742|24|24|Win32"
	reg       = []uint32{
		1937774191,
		1226093241,
		388252375,
		3666478592,
		2842636476,
		372324522,
		3817729613,
		2969243214,
	}
	stringsMap = map[string]string{
		"s0": "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=",
		"s1": "Dkdpgh4ZKsQB80/Mfvw36XI1R25+WUAlEi7NLboqYTOPuzmFjJnryx9HVGcaStCe=",
		"s2": "Dkdpgh4ZKsQB80/Mfvw36XI1R25-WUAlEi7NLboqYTOPuzmFjJnryx9HVGcaStCe=",
		"s3": "ckdp1h4ZKsUB80/Mfvw36XIgR25+WQAlEi7NLboqYTOPuzmFjJnryx9HVGDaStCe",
		"s4": "Dkdpgh2ZmsQB80/MfvV36XI1R45-WUAlEixNLwoqYTOPuzKFjJnry79HbGcaStCe",
	}
)

func NewABogus(platform string) *ABogus {
	uaCode := []byte{
		34, 167, 211, 143, 231, 217, 33, 244,
		208, 33, 142, 226, 219, 0, 182, 214,
		50, 32, 197, 93, 75, 3, 223, 172,
		226, 95, 80, 143, 61, 49, 216, 112,
	}
	browser := browser
	if platform != "" {
		browser = generateBrowserInfo(platform)
	}
	browserLen := len(browser)
	browserCode := charCodeAt(browser)

	return &ABogus{
		chunk:       []byte{},
		size:        0,
		reg:         append([]uint32{}, reg...),
		uaCode:      uaCode,
		browser:     browser,
		browserLen:  browserLen,
		browserCode: browserCode,
	}
}

func charCodeAt(s string) []int {
	result := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		result[i] = int(s[i])
	}
	return result
}

func generateBrowserInfo(platform string) string {
	innerWidth := randInt(1280, 1920)
	innerHeight := randInt(720, 1080)
	outerWidth := randInt(innerWidth, 1920)
	outerHeight := randInt(innerHeight, 1080)
	screenX := 0
	screenY := randChoice([]int{0, 30})
	valueList := []int{
		innerWidth,
		innerHeight,
		outerWidth,
		outerHeight,
		screenX,
		screenY,
		0,
		0,
		outerWidth,
		outerHeight,
		outerWidth,
		outerHeight,
		innerWidth,
		innerHeight,
		24,
		24,
	}
	return strings.Join(intSliceToStringSlice(valueList), "|") + "|" + platform
}

func intSliceToStringSlice(ints []int) []string {
	strs := make([]string, len(ints))
	for i, v := range ints {
		strs[i] = fmt.Sprintf("%d", v)
	}
	return strs
}

func randInt(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(n.Int64()) + min
}

func randChoice(choices []int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(choices))))
	return choices[n.Int64()]
}

func (a *ABogus) generateString1(randomNum1, randomNum2, randomNum3 *float64) string {
	return fromCharCode(list1(randomNum1)) + fromCharCode(list2(randomNum2)) + fromCharCode(list3(randomNum3))
}

func fromCharCode(args []int) string {
	var result strings.Builder
	for _, code := range args {
		result.WriteByte(byte(code))
	}
	return result.String()
}

func list1(randomNum *float64) []int {
	return randomList(randomNum, 170, 85, 1, 2, 5, 170&1)
}

func list2(randomNum *float64) []int {
	return randomList(randomNum, 170, 85, 1, 0, 0, 0)
}

func list3(randomNum *float64) []int {
	return randomList(randomNum, 170, 85, 1, 0, 5, 0)
}

func randomList(a *float64, b, c, d, e, f, g int) []int {
	r := 0.0
	if a != nil {
		r = *a
	} else {
		r = randFloat() * 10000
	}
	v := []int{
		int(r),
		int(r) & 255,
		int(r) >> 8,
	}
	s := v[1]&b | d
	v = append(v, s)
	s = v[1]&c | e
	v = append(v, s)
	s = v[2]&b | f
	v = append(v, s)
	s = v[2]&c | g
	v = append(v, s)
	return v[len(v)-4:]
}

func randFloat() float64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return float64(n.Int64())
}

func (a *ABogus) generateString2(urlParams, method string, startTime, endTime int) string {
	paramsArray := a.generateParamsCode(urlParams)
	methodArray := a.generateMethodCode(method)
	endTime = endTimeOrDefault(endTime)
	startTime = startTimeOrDefault(startTime)
	list := list4(
		(endTime>>24)&255,
		paramsArray[21],
		int(a.uaCode[23]),
		(endTime>>16)&255,
		paramsArray[22],
		int(a.uaCode[24]),
		(endTime>>8)&255,
		endTime&255,
		(startTime>>24)&255,
		(startTime>>16)&255,
		(startTime>>8)&255,
		startTime&255,
		methodArray[21],
		methodArray[22],
		int(endTime/256/256/256/256),
		int(startTime/256/256/256/256),
		a.browserLen,
	)
	e := endCheckNum(list)
	list = append(list, a.browserCode...)
	list = append(list, e)
	return rc4Encrypt(fromCharCode(list), "y")
}

func endTimeOrDefault(endTime int) int {
	if endTime == 0 {
		return int(time.Now().UnixNano()/1e6) + randInt(4, 8)
	}
	return endTime
}

func startTimeOrDefault(startTime int) int {
	if startTime == 0 {
		return int(time.Now().UnixNano() / 1e6)
	}
	return startTime
}

func list4(a, b, c, d, e, f, g, h, i, j, k, m, n, o, p, q, r int) []int {
	return []int{
		44, a, 0, 0, 0, 0, 24, b, n, 0, c, d, 0, 0, 0, 1, 0, 239, e, o, f, g, 0, 0, 0, 0, h, 0, 0, 14, i, j, 0, k, m, 3, p, 1, q, 1, r, 0, 0, 0,
	}
}

func endCheckNum(a []int) int {
	r := 0
	for _, i := range a {
		r ^= i
	}
	return r
}

func rc4Encrypt(plaintext, key string) string {
	s := make([]int, 256)
	for i := 0; i < 256; i++ {
		s[i] = i
	}
	j := 0
	for i := 0; i < 256; i++ {
		j = (j + s[i] + int(key[i%len(key)])) % 256
		s[i], s[j] = s[j], s[i]
	}
	i, j := 0, 0
	var cipher bytes.Buffer
	for k := 0; k < len(plaintext); k++ {
		i = (i + 1) % 256
		j = (j + s[i]) % 256
		s[i], s[j] = s[j], s[i]
		t := (s[i] + s[j]) % 256
		cipher.WriteByte(byte(s[t] ^ int(plaintext[k])))
	}
	return cipher.String()
}

func (a *ABogus) generateParamsCode(params string) []int {
	return sm3ToArray(sm3Hash(params + endString))
}

func (a *ABogus) generateMethodCode(method string) []int {
	return sm3ToArray(sm3Hash(method + endString))
}

func sm3ToArray(data string) []int {
	result := make([]int, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		val, _ := hex.DecodeString(data[i : i+2])
		result[i/2] = int(val[0])
	}
	return result
}

func sm3Hash(data string) string {
	hash := sm3.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func (a *ABogus) getValue(urlParams interface{}, method string, startTime, endTime int, randomNum1, randomNum2, randomNum3 *float64) string {
	string1 := a.generateString1(randomNum1, randomNum2, randomNum3)
	urlParamsStr := ""
	switch v := urlParams.(type) {
	case string:
		urlParamsStr = v
	case map[string]string:
		params := make([]string, 0, len(v))
		for key, value := range v {
			params = append(params, fmt.Sprintf("%s=%s", key, value))
		}
		urlParamsStr = strings.Join(params, "&")
	}
	string2 := a.generateString2(urlParamsStr, method, startTime, endTime)
	result := string1 + string2
	return base64.StdEncoding.EncodeToString([]byte(result))
}

func main() {
	abogus := NewABogus("")
	urlParams := "device_platform=webapp&aid=6383&channel=channel_pc_web&update_version_code=170400&pc_client_type=1&version_code=190500&version_name=19.5.0&cookie_enabled=true&screen_width=1536&screen_height=864&browser_language=zh-SG&browser_platform=Win32&browser_name=Chrome&browser_version=126.0.0.0&browser_online=true&engine_name=Blink&engine_version=126.0.0.0&os_name=Windows&os_version=10&cpu_core_num=16&device_memory=8&platform=PC&downlink=10&effective_type=4g&round_trip_time=200&msToken=eHUQHQOZgTUdIyobTzkIBOxmCGDUmm6PTJzDi2PtXcP5XHCEKVrdcCNcfE8DhShYk_1P3llPBA6BYia8HNE7HcSMdpuV_XFOURF9gbEHnwolgwUzy9j12lL1UYekBA%3D%3D&aweme_id=6870423037087436046"
	method := "GET"
	startTime := 0
	endTime := 0
	randomNum1 := 1.0
	randomNum2 := 1.0
	randomNum3 := 1.0

	value := abogus.getValue(urlParams, method, startTime, endTime, &randomNum1, &randomNum2, &randomNum3)
	fmt.Println("Generated Value:", value)
}

