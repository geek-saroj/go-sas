// package handlers

// import (
// 	"net/http"
// 	"sas-pro/internal/models"
// 	"sas-pro/internal/repositories"

// 	"github.com/gin-gonic/gin"
// )

// func CreateProduct(c *gin.Context) {
// 	var input models.ProductCreateRequest
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	repo := repositories.NewProductRepository()
// 	product, err := repo.CreateProductWithVariants(input)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, product)
// }

// func GetProducts(c *gin.Context) {
// 	repo := repositories.NewProductRepository()
// 	products, err := repo.GetAllProducts()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, products)
// }

package handlers

import (
	"bytes"
	"log/slog"

	// "fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	// "os"
	"sas-pro/internal/models"
	// "sas-pro/internal/repositories"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var input models.ProductCreateRequest

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	file, err := c.FormFile("variables_file") 
	if err != nil {
		log.Println("No file uploaded or error:", err)
	} else {
	
		if err := c.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		log.Println("File uploaded successfully:", file.Filename)
	}

	log.Println("Received product data:", input)


	if err := sendProductDataToSilkAPI(input, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward product data to Silk API"})
		return
	}


	c.JSON(http.StatusCreated, gin.H{"data": input})
}


func sendProductDataToSilkAPI(input models.ProductCreateRequest, file *multipart.FileHeader) error {

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)


	_ = writer.WriteField("variable_name", input.VariableName)
	_ = writer.WriteField("direction", input.Direction)
	_ = writer.WriteField("default_value", input.DefaultValue)
	_ = writer.WriteField("use_default_value", input.UseDefaultValue)
	_ = writer.WriteField("variable_type", input.VariableType)

	if input.LinkSlug != "" {
		_ = writer.WriteField("link_slug", input.LinkSlug)
	}

	if file != nil {
		filePart, err := writer.CreateFormFile("variables_file", file.Filename)
		if err != nil {
			log.Println("Error creating form file:", err)
			return err
		}
		f, err := file.Open()
		if err != nil {
			log.Println("Error opening file:", err)
			return err
		}
		defer f.Close()
		_, err = io.Copy(filePart, f)
		if err != nil {
			log.Println("Error copying file content:", err)
			return err
		}
	}


	err := writer.Close()
	if err != nil {
		log.Println("Error closing writer:", err)
		return err
	}

	slog.Info("Sending request to Silk API",&requestBody)

	req, err := http.NewRequest("POST", "https://api.goecl.silkinv.com/api/entries/import/macro-economical-variables/", &requestBody)
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	
	
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request to Silk API:", err)
		return err
	}
	defer resp.Body.Close()

	respBody := new(bytes.Buffer)
	respBody.ReadFrom(resp.Body)
	log.Println("Response from Silk API:", resp.Status)
	log.Println("Response body:", respBody.String())

	return nil
}
