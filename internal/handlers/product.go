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

//		c.JSON(http.StatusOK, products)
//	}
package handlers

import (
	"bytes"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	// "sas-pro/internal/models"
)

func CreateProduct(c *gin.Context) {
	
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB max size for file uploads
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}
	formData := c.Request.MultipartForm

	slog.Info("Received form data:", formData)

	// Check if formData is nil
	if formData == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No form data found"})
		return
	}
	serviceType, exists := formData.Value["service_type"]
	if !exists || len(serviceType) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'service_type' is required"})
		return
	}

 slog.Info("service Type:", serviceType)
	// Process the form fields dynamically
	if formData.Value != nil {
		for key, values := range formData.Value {
			if key == "service_type" {
				continue
			}
			// Log or handle non-file fields
			for _, value := range values {
				log.Printf("Field %s: %s\n", key, value)
			}
		}
	}

	// Iterate over file fields (for handling files in form-data)
	if formData.File != nil {
		for _, files := range formData.File {
			for _, file := range files {
				// Save each file to the server
				if err := c.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
					return
				}
				log.Println("File uploaded successfully:", file.Filename)
			}
		}
	} else {
		log.Println("No files received in the request")
	}

	
	slog.Info("final Received form data:", formData)

	// Send the received form data (including files) to Silk API
	if err := sendProductDataToSilkAPI(formData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward data to Silk API"})
		return
	}

	// Respond back with the received form data (for testing)
	c.JSON(http.StatusCreated, gin.H{"formData": formData})
}

func sendProductDataToSilkAPI(formData *multipart.Form) error {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add all form fields (non-file)
	if formData.Value != nil {
		for key, values := range formData.Value {
			for _, value := range values {
				_ = writer.WriteField(key, value)
			}
		}
	}

	// Add all files to the request body
	if formData.File != nil {
		for key, files := range formData.File {
			for _, file := range files {
				filePart, err := writer.CreateFormFile(key, file.Filename)
				if err != nil {
					log.Println("Error creating form file:", err)
					return err
				}

				// Open the file and copy its content to the multipart form data
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
		}
	}

	// Close the writer
	err := writer.Close()
	if err != nil {
		log.Println("Error closing writer:", err)
		return err
	}

	// Prepare and send the request to Silk API
	req, err := http.NewRequest("POST", "https://api.goecl.silkinv.com/api/entries/forecast/macro-economical-scenarios/", &requestBody)
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	// Set content type for the multipart request
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request using HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request to Silk API:", err)
		return err
	}
	defer resp.Body.Close()

	// Read and log the response from Silk API
	respBody := new(bytes.Buffer)
	respBody.ReadFrom(resp.Body)
	log.Println("Response from Silk API:", resp.Status)
	log.Println("Response body:", respBody.String())

	return nil
}

