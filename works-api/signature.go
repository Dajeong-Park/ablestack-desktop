package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/gofrs/uuid"
	"os"
	"sort"
	"strings"
	"time"
)

type SortMoldParams []MoldParams
type MoldParams map[string]string
type Uuids struct {
	Uuid, UuidUse       string
	CreateDate, Removed time.Time
}

func (s SortMoldParams) Len() int {
	return len(s)
}
func (s SortMoldParams) Less(i, j int) bool {
	for keyi, _ := range s[i] {
		for keyj, _ := range s[j] {
			return keyi < keyj
		}
	}
	return false
}

func (s SortMoldParams) Swap(i, j int) {
	for keyi, valuei := range s[i] {
		for keyj, valuej := range s[j] {
			s[i][keyj] = valuej
			s[j][keyi] = valuei
			delete(s[i], keyi)
			delete(s[j], keyj)
			return
		}
	}
}
func makeStringParams(params []MoldParams) string {
	var result string

	params1 := []MoldParams{
		{"apikey": os.Getenv("MoldApiKey")},
		{"response": "json"},
	}
	params = append(params, params1...)
	sort.Sort(SortMoldParams(params))
	fmt.Println("--------------------------------------------------")
	fmt.Println(params)
	fmt.Println("--------------------------------------------------")
	for _, tuple := range params {
		for key, value := range tuple {
			result = result + key + "=" + value + "&"
		}
	}
	result = strings.TrimRight(result, "&")
	fmt.Println("makeStringParams 결과")
	fmt.Println(result)

	return result
}

func makeSignature(payload string) string {
	secretkey := os.Getenv("MoldSecretKey")
	strurl := strings.Replace(strings.ToLower(payload), "+", "%20", -1)

	secret := []byte(secretkey)
	message := []byte(strurl)
	hash := hmac.New(sha1.New, secret)
	hash.Write(message)
	strHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	return strHash
}

func getUuid(payload string) string {

	db, err := sql.Open(os.Getenv("MsqlType"), os.Getenv("DbInfo"))
	if err != nil {
		fmt.Println("DB connect error")
		fmt.Println(err)
	}
	defer db.Close()
	var uuidValue uuid.UUID
	fmt.Println("DB connect success")
	for i := 0; i < 10; i++ {
		var count int
		uuidValue, err = uuid.NewV4()
		if err != nil {
			log.Error("UUID 생성 오류가 발생하였습니다.")
		}
		err = db.QueryRow("SELECT count(*) FROM uuids where uuid = ? and removed is null", uuidValue).Scan(&count)
		if err != nil {
			log.Error("UUID 중복 확인 쿼리에서 에러가 발생하였습니다.")
			log.Error(err.Error())
		}
		if count == 0 {
			result, err := db.Exec("INSERT INTO uuids(uuid, uuid_use, create_date) VALUES (?, ?, now())", uuidValue, payload)
			if err != nil {
				log.Error("UUID 생성 후 DB Insert 중 오류가 발생하였습니다.")
				log.Error(err)
				break
			}
			n, err := result.RowsAffected()
			if n == 1 {
				log.Info("UUID 값이 정상적으로 생성되었습니다.")
				goto END
			}
		}
	}
END:
	return uuidValue.String()
}
