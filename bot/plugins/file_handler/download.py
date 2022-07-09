import asyncio
from pyrogram import filters, Client
from pyrogram.errors import FloodWait
from bot.document_processor.base import DocumentProccesor
from bot.document_processor.factory import DocumentProcessorFactory
from bot.utils.logging import logger


@Client.on_message(
    ~filters.sticker & ~filters.via_bot & ~filters.forwarded
    & filters.command(
        commands=(["Download", "download", "Downloads", "downloads"])))
async def download(client, message):

    logger.info("God asked me to download something")

    try:
        if (len(message.command) == 1):

            logger.info("No download URL were provided")
            await message.reply_text(
                "No download link was provided.\nPlease provide one")
            return

        download_url: str = message.command[1]
        logger.info("Found download url as {download_url}")
        replied_message = await message.reply_text("Starting the download for you")

        handler: DocumentProccesor = DocumentProcessorFactory.create_document_processor(
            download_url, replied_message)

        file_name: str = await handler.download(message.from_user.id,
                                                download_url)
        if not file_name:
            raise Exception("File name is empty")
        logger.info(f"Downloaded file at {file_name}")
        await replied_message.edit_text("Downloaded file at " + file_name)

    except FloodWait as e:
        logger.error(f"Floodwait: Sleeping for {e.value} seconds")
        await asyncio.sleep(e.value)

    except Exception as e:
        logger.exception(e)
        await replied_message.edit_text(
            "Download failed.\nPlease check the link and try again")
