module github.com/bytebone/verbilobot

go 1.22.5

require (
	github.com/go-telegram/bot v1.6.1
	github.com/joho/godotenv v1.5.1
	github.com/jpoz/groq v0.0.0-20240513145022-7a02894105a0
	github.com/u2takey/ffmpeg-go v0.5.0
)

require (
	github.com/aws/aws-sdk-go v1.55.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/u2takey/go-utils v0.3.1 // indirect
)

replace github.com/jpoz/groq v0.0.0-20240513145022-7a02894105a0 => github.com/bytebone/groq v0.0.0-20240909145308-1341f117c3cb
