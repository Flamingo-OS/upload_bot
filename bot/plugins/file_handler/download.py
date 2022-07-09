from pyrogram import filters, Client
from bot.document_processor.base import DocumentProccesor
from bot.document_processor.factory import DocumentProcessorFactory
from bot.utils.logging import logger
from bot.utils.messaging import edit_message, reply_message


@Client.on_message(
    ~filters.sticker & ~filters.via_bot & ~filters.forwarded
    & filters.command(
        commands=(["Download", "download", "Downloads", "downloads"])))
async def download(client, message):

    logger.info("God asked me to download something")

    try:
        if (len(message.command) == 1):

            logger.info("No download URL were provided")
            await reply_message(message,
                                "No download link was provided.\nPlease provide one")
            return

        download_url: str = message.command[1]
        logger.info("Found download url as {download_url}")
        replied_message = await reply_message(message, "Starting the download for you")

        handler: DocumentProccesor = DocumentProcessorFactory.create_document_processor(
            download_url, replied_message)

        file_name: str = await handler.download(message.from_user.id,
                                                download_url)
        if not file_name:
            raise Exception("File name is empty")
        logger.info(f"Downloaded file at {file_name}")
        await edit_message(replied_message, "Downloaded file at " + file_name)

    except Exception as e:
        logger.exception(e)
        await edit_message(replied_message,
                           "Download failed.\nPlease check the link and try again")
