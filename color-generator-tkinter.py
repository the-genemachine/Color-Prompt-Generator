import openai
from dotenv import dotenv_values
import tkinter as tk
import ast

config = dotenv_values(".env")

openai.api_key = config["OPENAI_API_KEY"]

def generate_palette(prompt_text):
    prompt = f"""
    You are a color palette generating assistant that responds to text prompts for color palettes
    You should generate color palettes that are aesthetically pleasing and match the text prompt
    You should generate 5 colors per palette
    Desired Format: a JSON array of hex color codes
    Text : {prompt_text}
    Result : ["#FFC300", "#FF5733", "#C70039", "#900C3F", "#581845"]
    """

    response = openai.Completion.create(
        model="text-davinci-003",
        prompt=prompt,
        max_tokens=200
    )
    
    colors = ast.literal_eval(response.choices[0].text.strip())
    display_colors(colors)

def display_colors(colors):
    root = tk.Tk()

    for color in colors:
        frame = tk.Frame(root, bg=color, height=50, width=100)
        frame.pack(fill='both')

    root.mainloop()

generate_palette("A beautiful sunset")

