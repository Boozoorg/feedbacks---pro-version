package handlers

import (
	"encoding/json"
	"errors"
	"feedbacks/db"
	"feedbacks/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	SUPERVISOR = 1
	SUPPORT    = 2
	TERMINAL   = 3
	DEALER     = 4
	MONITORING = 6
	MANAGER    = 7
	AFFILIATE  = 8
	DEVELOPER  = 9
)

func sendRequest(method, uri string, reader io.Reader, respStruct interface{}, headers map[string]string) error {
	req, err := http.NewRequest(method, uri, reader)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: 90 * time.Second,
	}
	if headers != nil {
		for s, v := range headers {
			req.Header.Set(s, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		if respStruct != nil {
			err = json.Unmarshal(body, &respStruct)
			if err != nil {
				return err
			}
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body))
	}
	return nil
}

func getUser(c *gin.Context) (user models.TUser, err error) {
	userID := c.Query("userID")
	err = db.GetPGSQL().Raw("select * from t_users where id = ? limit 1", userID).Scan(&user).Error
	if err != nil {
		log.Println("getUser err:", err.Error())
		return user, err
	}
	return user, nil
}

func checkAccess(user *models.TUser) bool {
	if user.RoleID == SUPERVISOR || user.RoleID == DEVELOPER {
		return true
	}
	if user.RoleID != int64(TERMINAL) || (user.RoleID == DEALER && user.TerminalID == 0) {
		return false
	}
	var access bool
	err := db.GetPGSQL().Raw("SELECT orzu_cash_out from terminals where id = ?", user.TerminalID).Scan(&access).Error
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if !access {
		return false
	}
	return true
}
