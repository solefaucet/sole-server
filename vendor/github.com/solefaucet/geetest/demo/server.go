package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/solefaucet/geetest"
)

const (
	captchaID  = "b46d1900d0a894591916ea94ea91bd2c"
	privateKey = "36fc3fe98530eea08dfc6ce76e3d24c4"
)

func main() {
	g := geetest.New(captchaID, privateKey, false, 5*time.Second, 5*time.Second, 8)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/geetest/validate", func(writer http.ResponseWriter, req *http.Request) {
		challenge := req.PostFormValue("geetest_challenge")
		validate := req.PostFormValue("geetest_validate")
		seccode := req.PostFormValue("geetest_seccode")
		result := struct {
			Status string `json:"status"`
			Info   string `json:"info"`
		}{}
		ok, err := g.Validate(challenge, validate, seccode)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if ok {
			result.Status = "success"
			result.Info = "登陆成功"
		} else {
			result.Status = "fail"
			result.Info = "登陆失败"
		}

		resultBytes, _ := json.Marshal(result)
		writer.Write(resultBytes)
	})

	http.HandleFunc("/geetest/register", func(writer http.ResponseWriter, req *http.Request) {
		challenge, err := g.Register()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		result, _ := json.Marshal(map[string]interface{}{
			"gt":        g.CaptchaID(),
			"challenge": challenge,
			"success":   1,
		})
		writer.Write(result)
	})

	log.Println("listen and serve on 8080...")
	http.ListenAndServe("127.0.0.1:8080", nil)
}
