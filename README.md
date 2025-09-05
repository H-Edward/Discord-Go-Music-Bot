# Discord-Go-Music-Bot

This project is a simple-to-use Discord bot you can deploy almost anywhere to play music and audio within Discord servers.

There is a small amount of configuration needed to get it up and running, and once set up, there is little maintenance required to keep it operational.

## Table of Contents

- [Discord-Go-Music-Bot](#discord-go-music-bot)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
    - [Native vs Docker](#native-vs-docker)
    - [Native](#native)
      - [Updating with Native](#updating-with-native)
    - [Docker](#docker)
      - [Updating with Docker](#updating-with-docker)
  - [Environment Variables](#environment-variables)
  - [Usage within Discord](#usage-within-discord)
  - [Contributing](#contributing)
  - [License](#license)

## Features

The bot supports the following features:

- Playing audio from YouTube where you can:
  - Add videos via a URL
  - Search for videos
  - Queue more videos
  - Pause and resume
  - Skip to the next video
  - Stop all playback
  - Change the volume
  - See the current volume
  - Show the current queue

- Nuking messages (if the user has message management permissions)

- Getting the bot's uptime

- Getting the bot's version

- Pinging the bot to make sure it's active

- Getting a help message

- Slash commands and prefix commands \(prefix is `!`\)

_more features to come as they are requested_

## Installation

### Native vs Docker

You may choose to run the bot natively on your system or in a Docker container. The native installation is simpler and easier to follow, as well as maintain, but the Docker installation provides a potentially easier-to-maintain system in the long run, especially if you are familiar with Docker and its commands.

### Native

You must have the following programs/dependencies on your system:

- `yt-dlp` (a fork of youtube-dl)
- `ffmpeg` (for audio processing)
- `make` (for building the bot)
- `go` (the Go programming language, version 1.23.5 or later)
- `git` (for cloning the repository)
- `libopus0`, `libopus-dev`, `libopusenc0`, `libopusfile-dev`, `opus-tools` (for opus audio processing)

On a Debian/Ubuntu-based system, you can install the required dependencies with:

```bash
sudo apt update
sudo apt install ffmpeg libopus0 libopus-dev libopusenc0 libopusfile-dev opus-tools golang make yt-dlp git

# You may need to install `yt-dlp` manually since YouTube sometimes interferes with yt-dlp's video download process.
# You may also opt for installing Go manually if your version from apt is too old.
# See github.com/yt-dlp/yt-dlp/wiki/Installation for yt-dlp
# and go.dev/doc/install for Go
```

After installing these, it's a good idea to check yt-dlp is working by testing it, e.g.

```bash
yt-dlp https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

Then you need to clone the code repository and compile the bot. You can do this with the following commands:

```bash
git clone https://github.com/H-Edward/Discord-Go-Music-Bot
cd Discord-Go-Music-Bot
make
```

You should now see a new file called `music-bot`; this is the executable file.

Now you must make sure you have your bot token from Discord's developer portal,
which can be found at https://discord.com/developers/applications. Then click on your application, then click on `Bot` in the sidebar, then copy the token from this page (you may need to click `Reset Token`).

The token should be written to the file `.env` in the same directory as the `music-bot` executable.

This can be achieved by doing the following:

```bash
mv .env.sample .env 
nano .env
```

You will now see

`DISCORD_BOT_TOKEN=`

Paste in your token after the `=` sign and save the file.

You now have the bot ready and can run it using the following command:

```bash
./music-bot

## you may want to run it in the background using something like the screen command so the bot doesn't stop when you close the terminal
# sudo apt install screen
# screen -S music-bot
# ./music-bot
# then to detach from the screen session, press Ctrl+A then D

# to reattach to the screen session, run
# screen -r music-bot
```

To kill the bot safely, you can press `Ctrl+C` in the terminal where the bot is running.

#### Updating with Native

If you want to update the bot to the latest version, you can do so by running the following commands:

```bash
git pull origin main
make
```

Then you may run the bot again with:

```bash
./music-bot
```


### Docker

First, you must have Docker installed on your system. You can find instructions for installing Docker at https://docs.docker.com/get-docker/.

Then you must clone the repository and change into the directory:

```bash
git clone https://github.com/H-Edward/Discord-Go-Music-Bot
cd Discord-Go-Music-Bot
```

Now you must make sure you have your bot token from Discord's developer portal,
which can be found at https://discord.com/developers/applications.

Then click on your application, then click on `Bot` in the sidebar, then copy the token from this page (you may need to click `Reset Token`).

```bash
cp .env.sample .env 
nano .env
```

You will now see

`DISCORD_BOT_TOKEN=`

Paste in your token after the `=` sign and save the file.

Then to build the bot in Docker:

```bash
make docker-network-create
make docker-build
```

and to deploy the bot, you can run:

```bash
make docker-run
## This command has additional set options for security and to prevent resource hogging
```

_If you would like to use Docker but are unfamiliar, the `Makefile` has some additional commands to help manage the bot._

#### Updating with Docker

If you want to update the bot to the latest version, you can do so by running the following commands:

```bash
make docker-stop
make docker-rm 
git pull origin main
make docker-build
```

Then you may run the bot again with:

```bash
make docker-run
```

Alternatively,

```bash
git pull origin main
make docker-refresh-build
```

## Environment Variables

**You can set the following environment variables in the `.env` file to configure the bot:**

DISCORD_BOT_TOKEN

Your Discord bot token (required)

UNKNOWN_COMMANDS

Ignore - If set to `ignore`, the bot will treat the command as if its just another message

Error - If set to `error`, the bot will respond with an error message if an unknown command is used

Help - If set to `help`, the bot will respond with a help message if an unknown command is used (same as `!help` command)

## Usage within Discord

First, you can test the bot is working by messaging:

```txt
!ping
```

This should return a message saying `Pong` if the bot is running correctly.

Here are some other command examples for your reference:

```txt
Add a video to the queue by URL
!play https://www.youtube.com/watch?v=dQw4w9WgXcQ

Search for a video and add it to the queue
/search Rick Astley Never Gonna Give You Up

/pause

/resume

/volume 50

/volume 200

/volume 100

/skip

.stop

Show the current queue of videos
/queue

Delete the last 50 messages in the channel
/nuke 50
```

All the other commands can be seen by messaging:

```txt
/help
```

This will return a message with all the commands and their usage.

## Contributing

If you would like to contribute there are many ways you can help:

- **Report bugs**: If you find any bugs or issues, please report them on the [GitHub Issues page](https://github.com/H-Edward/Discord-Go-Music-Bot/issues)
- **Suggest features**: If you have any ideas for new features or improvements, please let me know on the [GitHub Issues page](https://github.com/H-Edward/Discord-Go-Music-Bot/issues)
- **Contribute code**: If you would like to contribute code, please fork the repository and submit a pull request and I will review it ASAP.

## License

This project is licensed under the GNU General Public License v3.0 (GPL-3.0). You can find the full license text in the `LICENSE` file in the root of the repository.
