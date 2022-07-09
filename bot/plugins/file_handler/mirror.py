import shutil
from typing import List
from pyrogram import filters, Client

from bot.constants import TEMP_FOLDER_PATH
from bot.document_processor.base import DocumentProccesor
from bot.document_processor.factory import DocumentProcessorFactory
from bot.utils.logging import logger
from bot.utils.messaging import edit_message, reply_message


@Client.on_message(~filters.sticker & ~filters.via_bot & ~filters.forwarded
                   & filters.command(commands=(["Mirror", "mirror"])))
async def mirror(client, message):

    logger.info("God asked me to mirror something")

    try:
        if (len(message.command) == 1):
            logger.info("No download URL were provided")
            await reply_message(message,
                                "No download link was provided.\nPlease provide one")
            return

        download_urls: List[str] = message.command[1:]
        logger.info(f"Found download url as {download_urls}")
        replied_message = await reply_message(message, "Starting the download for you")
        mirrored_url: List[str] = []

        for download_url in download_urls:
            handler: DocumentProccesor = DocumentProcessorFactory.create_document_processor(
                download_url, replied_message)

            file_name = await handler.download(message.from_user.id,
                                               download_url)
            logger.info(f"Downloaded file at {file_name}")

            if file_name is None:
                logger.info("File didn't download")
                await edit_message(replied_message,
                                   "Download failed because you are not a maintainer")
                return

            await edit_message(replied_message,
                               "Downloaded successfully. \nStarting upload now")

            url: str = await handler.upload(message.from_user.id, file_name)
            logger.info(f"Uploaded file at {url}")
            mirrored_url.append(url)

        shutil.rmtree(TEMP_FOLDER_PATH)
        await replied_message.delete()
        reply_text = f"Successfully uploaded file. you can find it at "
        for url in mirrored_url:
            reply_text += f"\n{url}"
        await reply_message(message, reply_text)

    except Exception as e:
        logger.exception(e)
        await replied_message.edit_text(
            "Mirror failed.\nPlease check the link and try again")
