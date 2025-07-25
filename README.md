# Discord-music-bot

A discord music bot that uses `yt-dlp` to stream audio, allowing for commands such as `play`, `pause`, `skip`.

## Installation

### Native

You must have `yt-dlp` and `ffmpeg` installed to path as well as a go compiler to build the bot.

Ensure yt-dlp is working by running a command such as `yt-dlp https://www.youtube.com/watch?v=dQw4w9WgXcQ`

You will also need the relevant opus libraries on your system (the `Dockerfile` maybe be a good reference for this for your system).

```bash
git clone https://github.com/H-Edward/discord-music-bot.git
cd discord-music-bot
make
```

then add your token to `.env` (change the filename of the sample)

then

```bash
./music-bot
```

### Docker

There is also a Dockerfile provided for containerisation

First, fill the `.env` with your token and other configuration options. You can copy the `.env.sample` file to `.env` and edit it as needed.

Then you must create the network for the bot to use:

```bash
make docker-network-create
```

After creating the network, you can build and run the bot using Docker:

```bash
make docker-build
make docker-run # has extra options for security and resource limits
```

#### Read the dockerfile for more options

## Usage

Within a Discord server, invite the bot using the link provided in the Discord developer portal.
You can then test the bot by messaging `!ping` in a channel the bot has access to.

To test functionality, join a voice channel and message `!play <url>` where `<url>` is a link to a video on YouTube or a link to an audio file. The bot will join the voice channel and start playing the audio from the video. You can then use `!pause`, `!resume`, `!skip`, and `!stop` to control playback.

## License

MIT
