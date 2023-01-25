package handlers

import (
	"feedbacks/db"
	"feedbacks/models"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary get FAQ
// @ID      get-faq
// @Produce json
// @Tags    Faq
// @Param   filter body     models.FaqFilters true "search question to faq"
// @Success 200    {object} string
// @Failure 404    {object} string
// @Router  /getFaqs [get]
func GetFaqs(c *gin.Context) {
	var filter models.FaqFilters
	err := c.Bind(&filter)
	if err != nil {
		log.Println("bind err: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	query := db.GetPGSQL().Table("faqs f")
	if filter.ID != nil {
		query = query.Where("f.id = ?", filter.ID)
	}
	if filter.Question != nil {
		query = query.Where("f.question = ?", filter.Question)
	}
	if filter.Key != nil {
		query = query.Where("f.key = ?", filter.Key)
	}
	if filter.CategoryID != nil {
		query = query.Where("f.category_id = ?", filter.CategoryID)
	}
	if filter.Product != nil {
		query = query.Where("c.product = ?", filter.Product)
	}
	if filter.Project != nil {
		query = query.Where("c.project = ?", filter.Project)
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageLimit <= 0 {
		filter.PageLimit = 15
	}

	var count int64
	err = query.Count(&count).Error
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		c.Abort()
		return
	}
	pages := math.Ceil(float64(count) / float64(filter.PageLimit))
	var resp []models.Faq
	err = query.Select("f.*").
		Joins("LEFT JOIN faq_categories c").
		Order("oc.id desc").
		Limit(int(filter.PageLimit)).Offset(int((filter.Page - 1) * filter.PageLimit)).Scan(&resp).Error
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		c.Abort()
		return
	}
	if resp == nil {
		resp = []models.Faq{}
	}
	c.JSON(http.StatusOK, gin.H{
		"resp":        resp,
		"total_pages": pages,
		"page":        filter.Page,
		"total_rows":  count,
	})
}

// @Summary Update FAQ
// @ID      update-faq
// @Produce json
// @Tags    Faq
// @Param   FAQ_data body     models.Faq true "Change the FAQ from db"
// @Success 200      {object} string
// @Failure 404      {object} string
// @Router  /getFaqs [Put]
func UpdateFaqs(c *gin.Context) {
	var newFaq, existingFaq models.Faq
	err := c.ShouldBindJSON(&newFaq)
	if err != nil {
		log.Println("bind err: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		c.Abort()
		return
	}
	err = db.GetPGSQL().Find(&existingFaq, newFaq.ID).Error
	if err != nil {
		log.Println("sql err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		c.Abort()
		return
	}
	if existingFaq.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": err})
		c.Abort()
		return
	} else {
		err = db.GetPGSQL().Updates(&newFaq).Error
		if err != nil {
			log.Println("sql err: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			c.Abort()
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "успех"})
}

// @Summary Create FAQ
// @ID      create-faq
// @Produce json
// @Tags    Faq
// @Param   FAQ_data body     models.Faq true "Create FAQ and save it in db"
// @Success 200      {object} string
// @Failure 404      {object} string
// @Router  /getFaqs [post]
func CreateFaqs(c *gin.Context) {
	var newFaq models.Faq
	err := c.ShouldBindJSON(&newFaq)
	if err != nil {
		log.Println("bind err: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		c.Abort()
		return
	}

	err = db.GetPGSQL().Create(&newFaq).Error
	if err != nil {
		log.Println("sql err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "успех"})
}

// @Summary Delete FAQ
// @ID      delete-faq
// @Produce json
// @Tags    Faq
// @Param   id  path     string true "id of which faq we must delete" Format(id)
// @Success 200 {object} string
// @Failure 404 {object} string
// @Router  /getFaqs [delete]
func DeleteFaqs(c *gin.Context) {
	var faq models.Faq
	id, _ := strconv.Atoi(c.Query("id"))
	faq.ID = int64(id)
	if faq.ID == 0 {
		log.Println("empty id")
		c.JSON(http.StatusBadRequest, gin.H{"message": "empty id query"})
		c.Abort()
		return
	}
	err := db.GetPGSQL().Delete(&faq, faq.ID).Error
	if err != nil {
		log.Println("sql err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "успех"})
}
