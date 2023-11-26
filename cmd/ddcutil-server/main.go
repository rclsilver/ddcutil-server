package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/rclsilver/ddcutil-server/pkg/ddc"
)

const (
	USB_C  ddc.InputSource = 0x1b
	DP     ddc.InputSource = 0x0f
	HDMI_1 ddc.InputSource = 0x11
	HDMI_2 ddc.InputSource = 0x12
)

const (
	USB_C_NAME  = "USB-C"
	DP_NAME     = "DP"
	HDMI_1_NAME = "HDMI-1"
	HDMI_2_NAME = "HDMI-2"
)

var (
	sourcesByName = map[ddc.InputSourceName]ddc.InputSource{
		USB_C_NAME:  USB_C,
		DP_NAME:     DP,
		HDMI_1_NAME: HDMI_1,
		HDMI_2_NAME: HDMI_2,
	}

	sourcesByID = map[ddc.InputSource]ddc.InputSourceName{
		USB_C:  USB_C_NAME,
		DP:     DP_NAME,
		HDMI_1: HDMI_1_NAME,
		HDMI_2: HDMI_2_NAME,
	}
)

type StatusResponse struct {
	Status string `json:"status"`
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, StatusResponse{Status: "OK"})
}

func getInputSource(c *gin.Context) {
	src, err := ddcClient.GetInputSource(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to get current input source: %s", err)})
		return
	}

	if _, ok := sourcesByID[src]; !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unknown input source %d", src)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"source": sourcesByID[src]})
}

func setInputSource(c *gin.Context) {
	var input struct {
		Source ddc.InputSourceName `json:"source"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if _, ok := sourcesByName[input.Source]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input source"})
		return
	}

	if err := ddcClient.SetInputSource(c, sourcesByName[input.Source]); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unable to change the input source: %s", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Input source set successfully"})
}

var ddcClient = ddc.Client{}

func main() {
	var listenHost string
	var listenPort int

	if envListenHost, ok := os.LookupEnv("LISTEN_HOST"); ok {
		listenHost = envListenHost
	} else {
		listenHost = "localhost"
	}

	if envListenPort, ok := os.LookupEnv("LISTEN_PORT"); ok {
		v, err := strconv.ParseInt(envListenPort, 10, 64)
		if err != nil {
			log.Fatalf("invalid listen port: %q", envListenPort)
		}
		listenPort = int(v)
	} else {
		listenPort = 8080
	}

	flag.StringVar(&listenHost, "host", listenHost, "Host to listen on")
	flag.IntVar(&listenPort, "port", listenPort, "Port to listen on")
	flag.Parse()

	router := gin.Default()

	router.GET("/mon/ping", pingHandler)
	router.GET("/input-source", getInputSource)
	router.POST("/input-source", setInputSource)

	router.Run(fmt.Sprintf("%s:%d", listenHost, listenPort))
}
