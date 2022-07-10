from abc import ABC, abstractmethod
import hashlib
from typing import List
from graph_onedrive import OneDrive
from bot import CLIENT_ID_ONEDRIVE, CLIENT_SECRET, TENANT, REFRESH_TOKEN
from bot.constants import BASE_URL, TEMP_FOLDER_PATH
from bot.database import maintainer_details

from bot.utils.parser import find_device
from bot.utils.progress import progress_callback
from bot.utils.logging import logger


class DocumentProccesor(ABC):
    """Base abstract class for downloading and uploading of all the files from the server.
    This class defines the bare amount of functions all the downloaders must have.
    Methods
        download: Downloads the file from the given url and stores it in a temporary location
        upload: Uploads the given file to gdrive. The file is deleted after the upload.
    """

    message: any

    def __init__(self, message):
        self.message = message

    @abstractmethod
    async def download(self, user_id: int, url: str) -> str:
        """Downloads the file from the given url and stores it in a temporary location.
        Args:
            url: The url of the file to be downloaded.
            message: The message to edit to show the progress of the download.
        """
        pass

    async def upload(self, user_id: int, file_name: str) -> str:
        """Uploads the given file to onedrive. The file is deleted after the upload.
        Args:
            file_path: The path of the file to be uploaded.
            folder_id: The id of the folder where the file will be uploaded.
        """
        logger.info("Starting upload to onedrive")
        sha5sum_file: str = await self.__find_sha5_sum__(file_name)

        # Upload path of the final file
        device = find_device(file_name)
        file_upload_path = "flamingo/A12.1/" + device

        try:
            # Use the context manager to manage a session instance
            my_drive = OneDrive(CLIENT_ID_ONEDRIVE, CLIENT_SECRET, TENANT,
                                "http://localhost:8080", REFRESH_TOKEN)

        except:
            logger.error("Failed to login to OneDrive")

        # Get the details of all the items in the root directory

        dir_to_travel: List[str] = file_upload_path.split("/")
        # Search through the root directory to find the file
        parent_folder_id = None
        dest_folder_id = None

        try:
            official_devices = maintainer_details.get_devices(user_id)
            if not official_devices:
                official_devices = []
            if device not in official_devices and not maintainer_details.is_maintainer_or_admin(
                    user_id):
                logger.info("This user is not a maintainer of this device")
                raise Exception("INVALID_DEVICE")
        except:
            logger.error("No devices found")

        try:
            items = my_drive.list_directory()
            for folder in dir_to_travel:
                dest_folder_id = None
                for item in items:
                    if "folder" in item and item.get("name") == folder:
                        dest_folder_id = item["id"]
                        break

                if dest_folder_id is None:
                    dest_folder_id = my_drive.make_folder(
                        folder, parent_folder_id)

                items = my_drive.list_directory(dest_folder_id)
                parent_folder_id = dest_folder_id

            # Upload the file
            await my_drive.upload_file(
                file_path=TEMP_FOLDER_PATH + sha5sum_file,
                parent_folder_id=dest_folder_id,
                if_exists="replace",
                chunk_size=1024 * 1024 * 40,
                verbose=False,
                callback=self.__callback__)

            await my_drive.upload_file(
                file_path=TEMP_FOLDER_PATH + file_name,
                parent_folder_id=dest_folder_id,
                if_exists="replace",
                chunk_size=1024 * 1024 * 40,
                verbose=False,
                callback=self.__callback__)

            url: str = BASE_URL + "d/" + file_upload_path.replace(
                "flamingo/", "") + "/" + file_name
            return url
        except:
            raise Exception("Failed to upload file")

    async def __callback__(self, progress: int, total: int, content_size: int):
        content: str = "Uploading file... " + str(progress / total * 100) + "%"
        logger.info(content)
        progress = progress / total * content_size

        try:
            await progress_callback(progress, content_size, self.message,
                                    "Uploading file...")
        except Exception as e:
            raise Exception("Something went wrong")

    async def __find_sha5_sum__(self, file_name: str) -> str:
        """Finds the sha5 sum of the given file.
        Args:
            file_path: The path of the file to find the sha5 sum of.
        """
        logger.info(f"Finding sha5 sum of file, {file_name}")
        BUF_SIZE = 65536  # lets read stuff in 64kb chunks!

        sha1 = hashlib.sha1()

        with open(TEMP_FOLDER_PATH + file_name, 'rb') as f:
            while True:
                data = f.read(BUF_SIZE)
                if not data:
                    break
                sha1.update(data)

        with open(TEMP_FOLDER_PATH + file_name + ".sha5sum", 'w') as f:
            f.write(sha1.hexdigest())

        return file_name + ".sha5sum"
