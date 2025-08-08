package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sas-pro/internal/repositories"

	"github.com/gin-gonic/gin"
)


func ProxyHandler(repo *repositories.ServiceRepository) gin.HandlerFunc {

	return func(c *gin.Context) {
		// Extract service type from request
		// var serviceType string
		var requestData map[string]interface{}
		contentType := c.ContentType()

		switch contentType {
		case "application/json":
			if err := c.ShouldBindJSON(&requestData); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
				return
			}
		case "application/x-www-form-urlencoded", "multipart/form-data":
			if err := c.Request.ParseForm(); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
				return
			}
			requestData = make(map[string]interface{})
			for k, v := range c.Request.PostForm {
				requestData[k] = v[0]
			}
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unsupported content type"})
			return
		}

		// Validate service type
		// st, ok := requestData["service_type"].(string)
		// if !ok || st == "" {
		// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "service_type is required"})
		// 	return
		// }
		// serviceType = st
		delete(requestData, "service_type")

		// Get target URL from database
		// targetURL, err := repo.GetURL(serviceType)
		// if err != nil {
		// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }

		targetURL := "http://localhost:5000/upload"

		// Prepare forwarded request
		var body io.Reader
		switch contentType {
		case "application/json":
			jsonData, _ := json.Marshal(requestData)
			body = bytes.NewBuffer(jsonData)
		default:
			form := url.Values{}
			for k, v := range requestData {
				form.Add(k, fmt.Sprintf("%v", v))
			}
			body = bytes.NewBufferString(form.Encode())
		}

		// Create and send request
		client := &http.Client{}
		req, err := http.NewRequest(c.Request.Method, targetURL, body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		// Copy headers
		req.Header = c.Request.Header.Clone()
		req.Header.Set("Content-Type", contentType)
		req.Header.Del("Content-Length")

		// Forward request
		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Failed to reach backend service"})
			return
		}
		defer resp.Body.Close()

		// Copy response to client
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}