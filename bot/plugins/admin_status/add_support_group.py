from pyrogram import filters, Client
from bot.utils.logging import logger
from bot.database.maintainer_details import maintainer_details
from bot.utils.messaging import edit_message, reply_message


@Client.on_message(~filters.sticker & ~filters.via_bot & ~filters.forwarded
                   & filters.command(commands=(["addSupport"])))
async def add_support_group(client, message):
    logger.info("Lets update the support group...")
    if not message.reply_to_message:
        await reply_message(message, "Reply to a user to add the support group")
        return
    if len(message.command) == 1:
        await reply_message(message, "Please specify a support group")
        return
    replied_message = await reply_message(message, "Ading support group")
    maintainer_details.add_support_group(message.from_user.id,
                                         message.reply_to_message.from_user.id,
                                         message.command[1])
    await edit_message(replied_message, "Successfully updated support group")
