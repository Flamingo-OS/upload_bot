# Flamingo Upload Bot
A Telegram bot (based on [GoTG](https://github.com/PaulSonOfLars/gotgbot)) to mirror and upload files from any direct link or a gdrive link to our [hosted website](https://downloads.e11z.net/flamingo/). It will also generate and push OTA while at it

## Setup config files

```bash
cp config.ini.json config.json
```

Now fill out the json file completely

## Build and run the bot

```bash
go build . -o bot
./bot
```

To get info on how to make various config files for microsoft onedrive, follow [this guide](https://rclone.org/onedrive/#getting-your-own-client-id-and-key) from rclone. To get various details about telegram related configs, follow the official [telegram website](https://core.telegram.org/api/obtaining_api_id#obtaining-api-id) to get them. For mongoDB related config vars. follow [this guide](https://www.mongodb.com/docs/manual/reference/connection-string/) from mongo. Finally for gdrive, you can follow official documentation from [google](https://github.com/googleapis/google-api-python-client/blob/main/docs/README.md#usage-guides)
