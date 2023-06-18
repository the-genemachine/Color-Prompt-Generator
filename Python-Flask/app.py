import openai
from dotenv import dotenv_values
from flask import Flask, render_template, request

config = dotenv_values(".env")
openai.api_key = config['OPENAI_API_KEY']

app = Flask(__name__, template_folder='templates')


@app.route('/palette', methods=["POST"])
def prompt_palette():
    if request.method == "POST":
        return render_template('palette.html')
    else:
        return render_template('index.html')


@app.route('/')
def index():
    response = openai.Completion.create(
        model="text-davinci-003",
        prompt="Explain in short the importance of an index.html file in a website.",
        max_tokens=200,
    )
    return response.choices[0].text

    #return render_template('index.html')


if __name__ == "__main__":
    app.run(debug=True)
