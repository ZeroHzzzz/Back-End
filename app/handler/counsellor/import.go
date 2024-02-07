package counsellor

import (
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
)

const savePath = ""

func ImportStudentInformation(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	err := c.Request.ParseMultipartForm(10 << 20) // 最大10MB
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to parse form")
		return
	}

	// 获取上传的文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to get file")
		return
	}
	defer file.Close()
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to open Excel file")
		return
	}
	// 遍历每个工作表和行
	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read rows")
		return
	}
	item := rows[0]                // 表头
	for _, row := range rows[1:] { // 切片，去除表头
		userId := row[0]
		userName := row[1]
		grade := row[2]
		profession := row[3]
		class := row[4]
		academicYear := row[5]
		user := models.Student{
			UserId:     userId,
			UserName:   userName,
			Class:      class,
			Profession: profession,
			Grade:      grade,
		}
		for colIndex, colValue := range row[6:] { //第六列后面就是成绩了
			itemName := item[colIndex]
			value := colValue
			sorce := models.Score{
				UserId:       userId,
				AcademicYear: academicYear,
				ItemName:     itemName,
				Grade:        value,
			}
			_ = service.InsertOne(c, "", "", sorce)
		}

		// 插入记录到MongoDB集合中
		_ = service.InsertOne(c, "", "", user)

	}

	c.String(http.StatusOK, "Data inserted successfully into MongoDB!")
}

type information struct {
	ItemName     string `json:"itemName"`
	AcademicYear string `json:"academicYear"`
	CorrectGrade string `json:"correctGrade"`
}

// 改正成绩的
func CorrectGrade(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var information information
	userId := c.Param("userID")
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	filter := bson.M{
		"_id":          userId,
		"academicYear": information.AcademicYear,
		"itemName":     information.ItemName,
	}
	modified := bson.M{
		"$set": bson.M{
			"grade": information.CorrectGrade,
		},
	}
	_ = service.UpdateOne(c, "", "", filter, modified)
	utils.ResponseSuccess(c, nil)
}
