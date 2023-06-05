import openai

from dotenv import dotenv_values

config = dotenv_values(".env")

openai.api_key = config["OPENAI_API_KEY"]

response = openai.Completion.create(
    model="text-davinci-003",
    prompt="Give me a tagline for a color prompt generator",
    max_tokens=50
)

print(response)