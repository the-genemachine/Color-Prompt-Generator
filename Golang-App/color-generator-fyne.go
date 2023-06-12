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
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

var w fyne.Window
var colors []string

func main() {
	a := app.New()
	w := a.NewWindow("Color Palette")
	w.SetContent(widget.NewLabel("Fetching color palette..."))
	w.Show()

	apiKey := os.Getenv("OPENAI_API_KEY")
	generatePalette("A beautiful colors rendered by a prism", apiKey)

	showPalette(w, colors)
	w.ShowAndRun()

}

func generatePalette(promptText, apiKey string) {
	prompt := fmt.Sprintf(`
Generate a palette of 5 colors for: '%s' in the format: ["#color1", "#color2", "#color3", "#color4", "#color5"]
`, promptText)

	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": "Bearer " + apiKey},
		JSON: map[string]interface{}{
			"model":       "gpt-3.5-turbo-0301",
			"messages":    []interface{}{map[string]interface{}{"role": "system", "content": "You are a helpful assistant."}, map[string]interface{}{"role": "user", "content": prompt}},
			"max_tokens":  150,
			"temperature": 0.3,
		},
	}

	resp, err := grequests.Post(apiEndpoint, ro)
	if err != nil {
		fmt.Printf("Couldn't make request: %v", err)
		return
	}

	var completion Completion
	err = json.Unmarshal(resp.Bytes(), &completion)
	if err != nil {
		fmt.Printf("Couldn't parse response: %v", err)
		return
	}

	if len(completion.Choices) == 0 {
		fmt.Println("No choices in response")
		return
	}

	// Get the text of the first choice
	text := completion.Choices[0].Message.Content

	// Find the opening and closing brackets of the JSON array in the response
	start := strings.Index(text, "[")
	end := strings.LastIndex(text, "]")

	// Check if brackets were found
	if start == -1 || end == -1 {
		fmt.Println("Couldn't find color array in response")
		return
	}

	// Extract the JSON part
	jsonPart := text[start : end+1]

	// var colors []string
	err = json.Unmarshal([]byte(jsonPart), &colors)
	if err != nil {
		fmt.Printf("Couldn't parse colors: %v", err)
		return
	}

	fmt.Printf("Received colors: %v\n", colors)

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

func showPalette(w fyne.Window, colors []string) {
	// Check if window is nil
	if w == nil {
		fmt.Println("Window is nil")
		return
	}

	box := container.NewVBox()

	for _, colorCode := range colors {
		color, err := parseHexColor(colorCode)
		if err != nil {
			fmt.Printf("Couldn't parse color '%s': %v\n", colorCode, err)
			continue
		}

		colorBox := canvas.NewRectangle(color)
		colorBox.SetMinSize(fyne.NewSize(50, 50)) // adjust size as needed
		box.Add(colorBox)
	}

	w.SetContent(box)
}
