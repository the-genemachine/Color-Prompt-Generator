package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/widget"
	"github.com/levigross/grequests"
)

type Completion struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	generatePalette("A beautiful sunset", apiKey)
}

func displayColors(colors []string) {
	a := app.New()
	w := a.NewWindow("Colors")

	boxes := []fyne.CanvasObject{}

	for _, color := range colors {
		box := canvas.NewRectangle(parseColor(color))
		box.Resize(fyne.NewSize(200, 150))
		boxes = append(boxes, box)
	}

	content := container.NewVBox(boxes...)

	w.SetContent(content)
	w.ShowAndRun()
}

func generatePalette(promptText, apiKey string) {
	prompt := fmt.Sprintf(`
    You are a color palette generating assistant that responds to text prompts for color palettes.
    Based on the following prompt: '%s', please generate a palette of 5 colors in the format: ["#color1", "#color2", "#color3", "#color4", "#color5"]
    `, promptText)

	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": "Bearer " + apiKey},
		JSON: map[string]interface{}{
			"engine":       "text-davinci-003",
			"prompt":       prompt,
			"max_tokens":   100,
			"temperature":  0.3,
		},
	}

	resp, err := grequests.Post("https://api.openai.com/v1/engines/davinci/completions", ro)
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

	colorsText := strings.Trim(completion.Choices[0].Text, " \n\r")

	var colors []string
	err = json.Unmarshal([]byte(colorsText), &colors)
	if err != nil {
		fmt.Printf("Couldn't parse colors: %v", err)
		return
	}

	displayColors(colors)
}

// parseColor is left as an exercise for the reader.
// It should convert a color string like "#ffffff" to a fyne.Color.
func parseColor(color string) fyne.Color {
	return fyne.Color{} // Replace this with the actual implementation.
}
