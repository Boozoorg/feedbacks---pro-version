package handlers

import (
	"encoding/json"
	"errors"
	"feedbacks/models"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

//ApiMiddleware to check headers sent by Web terminal API
func ApiMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		reqToken := c.Query("Authorization")
		log.Println("received: ", reqToken, "expected: ", models.Conf.TokenAuth)
		log.Println("url: ", c.Request.URL.Path)
		log.Println("url: ", c.Request.URL.RawPath, c.Request.URL.RawQuery)
		if reqToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "token is empty 99999"})
			//e.With(errors.New("token is empty 99999")).Write(c)
			c.Abort()
			return
		}
		if reqToken != models.Conf.TokenAuth {
			log.Println("received: ", reqToken, "expected: ", models.Conf.TokenAuth)
			c.JSON(http.StatusBadRequest, gin.H{"message": "secret key is not valid"})
			//e.With(errors.New("secret key is not valid")).Write(c)
			c.Abort()
			return

		}
		userID := (c.Query("userID"))
		log.Println("userID:", userID)
		if userID == "" {
			//e.With(errors.New("empty user_id")).Write(c)
			c.JSON(http.StatusBadRequest, gin.H{"message": "empty user_id"})
			c.Abort()
			return
		}
		log.Println("received: ", reqToken, "expected: ", models.Conf.TokenAuth)
	}
}

func interchangeMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		headers := map[string]string{}
		var resp interface{}
		if err := sendRequestToInterchange(c.Request.Method, c.Request.URL.Path[strings.Index(c.Request.URL.Path, "/cft")+len("/cft"):]+"?"+c.Request.URL.RawQuery, c.Request.Body, &resp, headers); err != nil {
			//e.With(err).Write(c)
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func sendRequestToInterchange(method, uri string, reader io.Reader, respStruct interface{}, headers map[string]string) error {
	req, err := http.NewRequest(method, models.Conf.UrlWtInterchange, reader)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: 15 * time.Second,
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
