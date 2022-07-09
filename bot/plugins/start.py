from pyrogram import filters, Client
from bot.utils.logging import logger
from bot.utils.messaging import reply_message


@Client.on_message(~filters.sticker & ~filters.via_bot
                   & ~filters.forwarded & filters.command(commands=(["start"]))
                   )
async def start(client, message):
    logger.info("Someone called for me?")
    await reply_message(message, "Hello User, I am alive :D")
