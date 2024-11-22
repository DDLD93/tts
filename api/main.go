package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
)

type Voice struct {
	Name       string `json:"Name"`
	ShortName  string `json:"ShortName"`
	Gender     string `json:"Gender"`
	Locale     string `json:"Locale"`
}

func main() {
	r := gin.Default()

	// Endpoint to fetch available voices
	r.GET("/voices", func(c *gin.Context) {
		cmd := exec.Command("edge-tts", "--list-voices")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch voices", "details": err.Error()})
			return
		}

		var voices []Voice
		if err := json.Unmarshal(out.Bytes(), &voices); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse voices", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, voices)
	})

	// Endpoint to synthesize text
	r.POST("/synthesize", func(c *gin.Context) {
		type SynthesizeRequest struct {
			Text       string  `json:"text" binding:"required"`
			Voice      string  `json:"voice" binding:"required"`
			Rate       string  `json:"rate"`
			Volume     string  `json:"volume"`
			Pitch      string  `json:"pitch"`
			ReturnFile bool    `json:"return_file"` // Controls file vs stream
		}

		var req SynthesizeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
			return
		}

		outputFile := "output.mp3"
		cmdArgs := []string{"--text", req.Text, "--voice", req.Voice, "--write-media", outputFile}

		// Optional parameters
		if req.Rate != "" {
			cmdArgs = append(cmdArgs, "--rate", req.Rate)
		}
		if req.Volume != "" {
			cmdArgs = append(cmdArgs, "--volume", req.Volume)
		}
		if req.Pitch != "" {
			cmdArgs = append(cmdArgs, "--pitch", req.Pitch)
		}

		cmd := exec.Command("edge-tts", cmdArgs...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Synthesis failed", "details": stderr.String()})
			return
		}
		defer os.Remove(outputFile)

		if req.ReturnFile {
			c.FileAttachment(outputFile, "synthesized_audio.mp3")
		} else {
			c.Writer.Header().Set("Content-Type", "audio/mpeg")
			c.File(outputFile)
		}
	})

	// Start the server
	r.Run(":8080")
}
