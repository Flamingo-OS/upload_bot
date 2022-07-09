import asyncio
import random
import re
from datetime import datetime
from typing import List
from pyrogram import filters, Client
from pyrogram.errors import FloodWait
from bot import CHANNEL_ID
from bot.database.maintainer_details import maintainer_details
from bot.utils.logging import logger
from bot.utils.messaging import reply_message
from bot.utils.parser import find_device, find_kosp_ver, parse_post_links

banner_photos = [
    {
        "follow":
        '',
        "banner":
        '',
        "support":
        ''
    },
]


@Client.on_message(~filters.sticker & ~filters.via_bot & ~filters.forwarded
                   & filters.command(commands=(["release"])))
async def create_post(client, message):
    logger.info("Creating a new post")

    try:
        if len(message.command) < 2:
            await reply_message(message, "Feed me the links senpai")
            return

        if not maintainer_details.is_maintainer_or_admin(message.from_user.id):
            await reply_message(message,
                                "You are not a maintainer, you can't create a post")
            return

        parsed_links: dict = parse_post_links(message.command[1:])

        download_link_str: str = "Download:\n"

        possible_links = ["full", "fastboot", "incremental", "boot"]

        for link_type in possible_links:

            if len(parsed_links[link_type]) > 0:
                download_link_str += f"""
        *    [{link_type.capitalize()}]({parsed_links[link_type][0]}) """
                for i in range(1, len(parsed_links[link_type])):
                    download_link_str += f"[Mirror {i}]({parsed_links[link_type][i]}) "

        device = find_device(message.command[1])

        maintainers: List[dict] = maintainer_details.get_maintainers(device)
        maintainer_str: str = ""
        for maintainer in maintainers:
            maintainer_str += f"[{maintainer['name']}](tg://user?id={maintainer['user_id']}) "

        device_support_group = maintainer_details.get_device_support_group(
            device)

        flamingo_version = find_kosp_ver(message.command[1])

        caption = f"""
        Flamingo OS {flamingo_version} | Android 12.1 | OFFICIAL

        Maintainers: {maintainer_str}

        Device: {device.capitalize()}
        Date: {datetime.today().strftime('%d-%m-%y')}

        {download_link_str}

        [Changelog](https://raw.githubusercontent.com/FlamingoOS-Devices/ota/main/{device}/changelog_{datetime.today().strftime('%Y_%m_%d')})

        Support group
        
        *   [Official](https://t.me/flamingo_common)"""

        if device_support_group:
            caption += f"""
        *   [Device]({device_support_group})"""

        # We don't have any banners yet
        # random_int = random.randint(0, 5)
        # random_follow = random.randint(0, 1)
        # logger.info(f"Downloading banner {random_int}")

        # await client.send_photo(chat_id=message.chat.id,
        #                         photo=banner_photos[random_int]["banner"],
        #                         parse_mode="md",
        #                         caption=caption)

        # await client.send_photo(chat_id=CHANNEL_ID,
        #                         photo=banner_photos[random_int]["banner"],
        #                         parse_mode="md",
        #                         caption=caption)
        # if random_follow == 1:
        #     await client.send_photo(chat_id=CHANNEL_ID,
        #                             photo=banner_photos[random_int]["follow"])
        # else:
        #     await client.send_photo(chat_id=CHANNEL_ID,
        #                             photo=banner_photos[random_int]["support"])

        await client.send_message(chat_id=CHANNEL_ID, text=caption),
        await client.send_message(chat_id=message.chat.id, text=caption),

    except FloodWait as e:
        logger.error(f"Floodwait: Sleeping for {e.value} seconds")
        await asyncio.sleep(e.value)
        await client.send_message(chat_id=CHANNEL_ID, text=caption),
        await client.send_message(chat_id=message.chat.id, text=caption),

    except:
        await reply_message(message, "Something went wrong")
        return
