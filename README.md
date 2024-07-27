# Verbilobot

![Go version](https://img.shields.io/github/go-mod/go-version/bytebone/verbilobot?style=flat-square)
![GitHub License](https://img.shields.io/github/license/bytebone/verbilobot?style=flat-square)
![GitHub CI Status](https://img.shields.io/github/actions/workflow/status/bytebone/verbilobot/ci.yml?style=flat-square&logo=github&label=CI)
![GitHub Issues](https://img.shields.io/github/issues/bytebone/verbilobot?style=flat-square)

Verbilobot is a Telegram bot written in Go that transcribes voice messages, video notes, and any other media files. It uses the [Groq API](https://groq.com) to transcribe the audio, and ffmpeg to convert any incoming audio to a format that Groq is happiest to convert.

## How to Run

However you plan to run the bot, make sure to rename the `.env.example` file to `.env` and fill in your Telegram Bot Token and Groq API Token.

### Local Go

To build and run the project locally, you will need to have Go installed on your machine.

On Linux:
```bash
git clone https://github.com/bytebone/verbilobot.git
cd verbilobot
go build -v -o verbilobot .
./verbilobot
```
Or on Windows: 
```pwsh
git clone https://github.com/bytebone/verbilobot.git
cd verbilobot
go build -v -o verbilobot.exe .
start verbilobot.exe
```

### With Docker

To build and run the project with Docker, you will need to have Docker installed on your machine.

```bash
git clone https://github.com/bytebone/verbilobot.git
cd verbilobot/docker
docker build -t verbilobot .
docker compose up
```
Thanks to Docker being awesome, this works the same on any platform. 

## Usage

The bot usually takes around 2 seconds to come online. Once the bot is running, you can forward any audio or video files to it to start the transcription process. Thanks to the high speeds at Groq, a minute of incoming audio takes only a few moments to transcribe and return to your chat. The main bottleneck you might notice is the local transcoding, which can take a *noticeable* amount of time to complete.