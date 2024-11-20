package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"maxcoba/config"
	"maxcoba/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	data            []models.Data
	modelReferences []models.ModelReference
	techReferences  []models.TechReference
)

func getAllData(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

func getDataByCode(ctx *gin.Context) {
	code := ctx.Param("code")
	for _, elemen := range data {
		if elemen.Code == code {
			translatedTech := translateTech(elemen.Tech)
			ctx.JSON(http.StatusOK, gin.H{
				"code":        elemen.Code,
				"name":        elemen.Name,
				"description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua",
				"model":       translateModel(elemen.Model),
				"tech":        strings.Join(translatedTech, ", "),
				"status":      strings.Title(elemen.Status),
			})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{
		"message": "Data not found",
	})
}

func createData(ctx *gin.Context) {
	var newData models.Data
	if err := ctx.ShouldBindJSON(&newData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data", 
			"error": err.Error(),
		})
		return
	}
	data = append(data, newData)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Data created successfully",
		"data":    newData,
	})
}

func updateData(ctx *gin.Context) {
	code := ctx.Param("code")
	for index, elemen := range data {
		if elemen.Code == code {
			if err := ctx.ShouldBindJSON(&data[index]); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Invalid data",
					"error":   err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Data updated successfully",
				"data":    data[index],
			})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{
		"message": "Data not found",
	})
}

func deleteData(ctx *gin.Context) {
	code := ctx.Param("code")
	for i, d := range data {
		if d.Code == code {
			data = append(data[:i], data[i+1:]...)
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Data deleted successfully",
			})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{
		"message": "Data not found",
	})
}

func getModelReferences(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"models": modelReferences,
	})
}

func getTechReferences(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"techs": techReferences,
	})
}

func filterData(ctx *gin.Context) {
	model := ctx.Query("model")
	tech := ctx.QueryArray("tech")

	var filteredData []models.Data

	for _, elemen := range data {
		if (model == "" || elemen.Model == model) &&
			(len(tech) == 0 || containsAll(strings.Split(elemen.Tech, ", "), tech)) {
			filteredData = append(filteredData, elemen)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": filteredData,
	})
}

func containsAll(source, filters []string) bool {
	for _, elemenFilter := range filters {
		found := false
		for _, elemenSource := range source {
			if strings.TrimSpace(elemenSource) == elemenFilter {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func translateModel(model string) string {
	for _, ref := range modelReferences {
		if ref.Key == model {
			return ref.Value
		}
	}
	return model
}

func translateTech(tech string) []string {
	var translate []string
	for _, elemenTech := range strings.Split(tech, ", ") {
		for _, ref := range techReferences {
			if ref.Key == elemenTech {
				translate = append(translate, ref.Value)
				break
			}
		}
	}
	return translate
}

func initData() {
	file, err := ioutil.ReadFile("data.txt")
	if err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	var jsonData struct {
		Data            []models.Data           `json:"data"`
		ModelReferences []models.ModelReference `json:"model_references"`
		TechReferences  []models.TechReference  `json:"tech_references"`
	}

	if err := json.Unmarshal(file, &jsonData); err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}

	data = jsonData.Data
	modelReferences = jsonData.ModelReferences
	techReferences = jsonData.TechReferences
}

func main() {
	// INIT DATA
	initData()

	// ENV
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading.env file")
	}

	// GIN ENGINE
	route := gin.Default()

	// ROUTES
	route.GET("/data", getAllData)
	route.GET("/data/:code", getDataByCode)
	route.POST("/data", createData)
	route.PUT("/data/:code", updateData)
	route.DELETE("/data/:code", deleteData)

	route.GET("/references/models", getModelReferences)
	route.GET("/references/techs", getTechReferences)

	// SERVER
	route.Run(config.PORT)
}
