import asyncio
from bot.utils.logging import logger
from pyrogram.errors import FloodWait

"""
A custom wrapper over all of edit text and reply text to avoid me getting flood waited and having to handle the said cases separately
"""


async def edit_message(message, text):
    try:
        return await message.edit_text(text)
    except FloodWait as e:
        logger.error(f"Floodwait: Sleeping for {e.value} seconds")
        await asyncio.sleep(e.value)
        return await message.edit_text(text)


async def reply_message(message, text):
    try:
        return await message.reply_text(text)
    except FloodWait as e:
        logger.error(f"Floodwait: Sleeping for {e.value} seconds")
        await asyncio.sleep(e.value)
        return await message.reply_text(text)
