package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
	"github.com/levigross/grequests"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

type Completion struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

var w fyne.Window

func main() {
	a := app.New()
	w := a.NewWindow("Color Palette")
	w.SetContent(widget.NewLabel("Fetching color palette..."))
	w.Show()

	apiKey := os.Getenv("OPENAI_API_KEY")
	generatePalette("A beautiful sunset", apiKey)

	w.ShowAndRun()
}

func generatePalette(promptText, apiKey string) {
	prompt := fmt.Sprintf("You are a color palette generating assistant that responds to text prompts for color palettes. Based on the following prompt: '%s', please generate a palette of 5 colors in the format: [\"#color1\", \"#color2\", \"#color3\", \"#color4\", \"#color5\"]", promptText)

	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": "Bearer " + apiKey},
		JSON: map[string]interface{}{
			"model":    "gpt-3.5-turbo",
			"messages": []map[string]string{{"role": "system", "content": prompt}},
		},
	}

	resp, err := grequests.Post(apiEndpoint, ro)
	if err != nil {
		fmt.Printf("Couldn't make request: %v", err)
		return
	}

	// Print the response status code and body for debugging
	fmt.Printf("Response status: %v\n", resp.StatusCode)
	fmt.Printf("Response body: %v\n", resp.String())

	// ...

	var completion struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	fmt.Println(resp.String())

	// Added for debugging purposes
	err = json.Unmarshal(resp.Bytes(), &completion)
	if err != nil {
		fmt.Printf("Couldn't parse response: %v", err)
		return
	}

	/* if err = json.Unmarshal(resp.Bytes(), &completion); err != nil {
		fmt.Printf("Couldn't parse response: %v", err)
		return
	} */

	// Check if there are any choices in the response
	if len(completion.Choices) == 0 {
		fmt.Println("No choices in response")
		return
	}

	colorResponse := completion.Choices[0].Message.Content

	// Search the response string for the color array
	colorStart := strings.Index(colorResponse, "[")
	colorEnd := strings.Index(colorResponse, "]")
	if colorStart == -1 || colorEnd == -1 {
		fmt.Println("Couldn't find color array in response")
		return
	}

	// colorsText := colorResponse[colorStart : colorEnd+1]

	/* 	var colors []string
	   	err = json.Unmarshal([]byte(colorsText), &colors)
	   	if err != nil {
	   		fmt.Printf("Couldn't parse colors: %v", err)
	   		return
	   	} */

	// fmt.Printf("Received colors: %v\n", colors)

	// Check if there are any choices in the response
	if len(completion.Choices) == 0 {
		fmt.Println("No choices in response")
		return
	}

	/* colorItems := make([]fyne.CanvasObject, len(colors))

	for i, colorCode := range colors {
		parsedColor, err := parseHexColor(colorCode)
		if err != nil {
			fmt.Printf("Error parsing color: %v\n", err)
			continue
		}

		colorBox := canvas.NewRectangle(parsedColor)
		colorBox.SetMinSize(fyne.NewSize(50, 50))
		colorItems[i] = colorBox
	}

	colorRow := container.NewHBox(colorItems...)
	w.SetContent(colorRow) */

	colorsText := strings.Trim(completion.Choices[0].Message.Content, " \n\r")

	var colors []string
	err = json.Unmarshal([]byte(colorsText), &colors)
	if err != nil {
		fmt.Printf("Couldn't parse colors: %v", err)
		return
	}

	fmt.Printf("Received colors: %v\n", colors)

	objects := []fyne.CanvasObject{}

	for _, color := range colors {
		rgba, err := parseHexColor(color)
		if err != nil {
			log.Fatalf("Could not parse color %s: %v", color, err)
		}
		objects = append(objects, canvas.NewRectangle(rgba))
	}

	/* 	objects := make([]fyne.CanvasObject, len(colors))

	   	for i, color := range colors {
	   		objects[i] = canvas.NewRectangle(parseHexColor(color))
	   	} */

	box := container.NewVBox(objects...)
	/* 	box := container.NewVBox()
	   	for _, hex := range colors {
	   		parsedColor, _ := parseHexColor(hex)
	   		box.Append(widget.NewLabelWithStyle(hex, fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}))
	   		box.Append(canvas.NewRectangle(parsedColor))
	   	} */
	w.SetContent(box)
	w.Show()
}

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s[0] != '#' {
		return c, fmt.Errorf("Invalid format")
	}

	hexColor := strings.TrimPrefix(s, "#")

	// Convert hex color to RGB
	rgb, err := strconv.ParseUint(hexColor, 16, 32)
	if err != nil {
		return c, err
	}

	c.R = uint8(rgb >> 16)
	c.G = uint8(rgb >> 8)
	c.B = uint8(rgb)

	return c, nil
}
