from pyrogram import filters, Client
from bot.utils.logging import logger
from bot.database.maintainer_details import maintainer_details
from bot.utils.messaging import edit_message, reply_message


@Client.on_message(~filters.sticker & ~filters.via_bot
                   & ~filters.forwarded & filters.command(commands=(["add"])))
async def add_maintainer(client, message):
    logger.info("Lets add a new maintainer")
    requester_id: int = message.from_user.id
    if not message.reply_to_message and len(message.command) < 1:
        await reply_message(message,
                            "Reply to a user or specify their userID to add them as a maintainer"
                            )
        return

    if message.reply_to_message:
        maintainer_id: int = message.reply_to_message.from_user.id
        name: str = message.reply_to_message.from_user.first_name
        device: str = message.command[1]

    elif len(message.command) == 4:
        maintainer_id: int = int(message.command[1])
        device: str = message.command[3]
        name: str = message.command[2]

    else:
        await reply_message(message, "Please specify a device, id and name")
        return

    replied_message = await reply_message(message, "Adding a maintainer")
    maintainer_details.add_maintainer(requester_id, maintainer_id, name,
                                      device)
    await edit_message(replied_message,
                       "Successfully added the user as an maintainer")
