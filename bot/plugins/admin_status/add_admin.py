import asyncio

from pyrogram import filters, Client
from bot.utils.logging import logger
from pyrogram.errors import FloodWait
from bot.database.maintainer_details import maintainer_details


@Client.on_message(~filters.sticker & ~filters.via_bot
                   & ~filters.forwarded
                   & filters.command(commands=(["addAdmin"])))
async def add_admin(client, message):
    logger.info("Lets add a new admin")
    try:
        if not message.reply_to_message:
            await message.reply_text("Reply to a user to add them as a admin")
            return
        replied_message = await message.reply_text("Adding a new admin...")
        maintainer_details.add_admin(message.from_user.id,
                                     message.reply_to_message.from_user.id)
        await replied_message.edit_text("Successfully added the user as an admin")

    except FloodWait as e:
        logger.error(f"Floodwait: Sleeping for {e.value} seconds")
        await asyncio.sleep(e.value)
