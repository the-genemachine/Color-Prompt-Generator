import openai

from dotenv import load_dotenv

config = dotenv_values(".env")

openai.api_key = config["OPENAI_API_KEY"]

