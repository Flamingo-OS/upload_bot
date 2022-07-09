import asyncio
from pyrogram import filters, Client
from pyrogram.errors import FloodWait
from bot.utils.logging import logger
from bot.database.maintainer_details import maintainer_details


@Client.on_message(~filters.sticker & ~filters.via_bot & ~filters.forwarded
                   & filters.command(commands=(["addSupport"])))
async def add_support_group(client, message):
    logger.info("Lets update the support group...")
    try:
        if not message.reply_to_message:
            await message.reply_text("Reply to a user to add the support group")
            return
        if len(message.command) == 1:
            await message.reply_text("Please specify a support group")
            return
        replied_message = await message.reply_text("Ading support group")
        maintainer_details.add_support_group(message.from_user.id,
                                             message.reply_to_message.from_user.id,
                                             message.command[1])
        await replied_message.edit_text("Successfully updated support group")

    except FloodWait as e:
        logger.error(f"Floodwait: Sleeping for {e.value} seconds")
        await asyncio.sleep(e.value)
