# GPT-3 AI Chatbot

This is a simple AI chatbot written in Go that utilizes OpenAI's GPT-3 API to generate responses. The chatbot allows users to have conversations with an AI that responds to their input.

## Installation and Setup

To use this chatbot, you will need to have a valid OpenAI API key. If you do not have one, you can sign up for a free trial at https://beta.openai.com/signup.

Clone this repository to your local machine.
Set the OPENAI_API_KEY environment variable to your OpenAI API key.

```bash
export OPENAI_API_KEY=<your API key>
```

Install the necessary dependencies by running go get in the project directory.

## Usage

To use the chatbot, simply run the main package in the command line. You will be prompted to enter a message, and the chatbot will respond with a message generated by GPT-3.

```go
go run main.go
```

## Example Usage

```console
You: Hello, how are you?
AI: I'm doing well, thank you for asking.
You: What is your name?
AI: I am an AI language model created by OpenAI.
You: Can you tell me a joke?
AI: Why did the tomato turn red? Because it saw the salad dressing!
You: exit
```

## Configuration

The chatbot can be configured using the following constants:

apiEndpoint: the URL for the OpenAI API endpoint to use (default: "https://api.openai.com/v1/chat/completions").

temperature: the temperature parameter to use when generating responses with GPT-3 (default: 0.5).

aiModel: the name of the GPT-3 model to use (default: "gpt-3.5-turbo-0301").

## Logging

The chatbot logs all interactions to a file called interaction.log in the project directory. The log file includes timestamps and the user's input and the AI's response for each interaction.

## TODO

Docker container
