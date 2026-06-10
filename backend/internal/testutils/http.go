package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
)

func DoRequest(handler http.Handler, method, path string, body any, cookies ...*http.Cookie) *http.Response {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w.Result()
}

func ReadBody(resp *http.Response) []byte {
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	resp.Body.Close()
	return buf.Bytes()
}

func GenAccessToken(j *jwtutils.JWT, idUser int32, roles []string, email, nik string) string {
	t, err := j.EncodeWithTTL(jwtutils.Claim{
		IDUser: idUser,
		Roles:  roles,
		Email:  email,
		NIK:    nik,
	}, 30*time.Minute)
	if err != nil {
		panic(fmt.Sprintf("GenAccessToken: %v", err))
	}
	return t
}

func AccessCookie(j *jwtutils.JWT, idUser int32, roles []string) *http.Cookie {
	return &http.Cookie{Name: "access_token", Value: GenAccessToken(j, idUser, roles, "test@example.com", "1234567890123456"), Path: "/"}
}

func GenRefreshCookie() *http.Cookie {
	id, err := uuid.NewV7()
	if err != nil {
		panic(fmt.Sprintf("GenRefreshCookie: %v", err))
	}
	return &http.Cookie{Name: "refresh_token", Value: id.String(), Path: "/"}
}
