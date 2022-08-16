# Flamingo Upload Bot
A Telegram bot (based on [Pyrogram](https://docs.pyrogram.org/)) to mirror and upload files from any direct link or a gdrive link to our [hosted website](https://downloads.e11z.net/flamingo/).

## Setup virtual environment

```bash
python -m venv .venv

source .venv/bin/activate
```

## Install dependencies.

```bash
pip install -r requirements.txt
```

## Run the bot

```bash
python3 -m bot
```

Make sure you have set config.ini properly.
To set it, copy sample.config.ini to config.ini and fill it as per the comments

```bash
cp sample.config.ini config.ini
```
You can also host the bot in heroku. Make sure to set the config var as per config.ini

To get info on how to make various config files for microsoft onedrive, follow [this guide](https://rclone.org/onedrive/#getting-your-own-client-id-and-key) from rclone. To get various details about telegram related configs, follow the official [telegram website](https://core.telegram.org/api/obtaining_api_id#obtaining-api-id) to get them. For mongoDB related config vars. follow [this guide](https://www.mongodb.com/docs/manual/reference/connection-string/) from mongo. Finally for gdrive, you can follow official documentation from [google](https://github.com/googleapis/google-api-python-client/blob/main/docs/README.md#usage-guides)
