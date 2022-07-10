import math
import os
import aiohttp

from bot.constants import TEMP_FOLDER_PATH
from bot.database.maintainer_details import maintainer_details
from bot.document_processor.base import DocumentProccesor
from bot.utils.parser import find_device
from bot.utils.progress import progress_callback
from bot.utils.logging import logger


class DirectLink(DocumentProccesor):

    async def download(self, user_id: int, url: str) -> str:

        if not os.path.exists(TEMP_FOLDER_PATH):
            os.mkdir(TEMP_FOLDER_PATH)

        logger.info("Downloading file from direct link")
        local_filename: str = ""
        for part in url.split('/'):
            if "Flamingo" and "Official" in part:
                local_filename = part
                break
        if not local_filename:
            raise Exception("Failed to find the file name")
        device: str = find_device(local_filename)

        try:
            official_devices = maintainer_details.get_devices(user_id)
            if not official_devices:
                official_devices = []
            if device not in official_devices and not maintainer_details.is_maintainer_or_admin(
                    user_id):
                logger.info("This user is not a maintainer of this device")
                return
        except:
            pass

        data = b''

        try:
            logger.info("Starting download................")
            session = aiohttp.ClientSession()
            response = await session.get(url)
            total = int(response.headers["Content-Length"])
            logger.info(f"The recieved file size is {total}")
            progress = 0
            last_update_value = -1
            while True:
                chunk = await response.content.read(1024 * 10)
                if not chunk:
                    break
                progress += len(chunk)
                data += chunk

                if math.floor(progress / total * 100) % 1 == 0 and math.floor(
                        progress / total * 100) != last_update_value:
                    await progress_callback(progress,
                                            total,
                                            message=self.message,
                                            text="Starting Download ... ")
                    last_update_value = math.floor(progress / total * 100)

            logger.info("Successfully downloaded the file")
            session.close

        except Exception as e:
            logger.exception(e)

        try:
            logger.info(
                f"Trying to write the file to {TEMP_FOLDER_PATH}{local_filename}"
            )
            with open(TEMP_FOLDER_PATH + local_filename, "wb") as f:
                f.write(data)
        except:
            return None

        return local_filename
