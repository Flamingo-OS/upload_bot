import io
import os

from googleapiclient.discovery import build
from googleapiclient.http import MediaIoBaseDownload
from google.oauth2 import service_account

from bot import CLIENT_EMAIL, CLIENT_ID_GDRIVE, CLIENT_X509_CERT_URL, PRIVATE_KEY_GDRIVE, PRIVATE_KEY_ID, PROJECT_ID_GDRIVE
from bot.constants import TEMP_FOLDER_PATH
from bot.database.maintainer_details import maintainer_details
from bot.document_processor.base import DocumentProccesor
from bot.utils.parser import find_device

from bot.utils.progress import progress_callback
from bot.utils.logging import logger

creds = {
    "type": "service_account",
    "project_id": PROJECT_ID_GDRIVE,
    "private_key_id": PRIVATE_KEY_ID,
    "private_key": PRIVATE_KEY_GDRIVE.replace("\\n", "\n"),
    "client_email": CLIENT_EMAIL,
    "client_id": CLIENT_ID_GDRIVE,
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url":
    "https://www.googleapis.com/oauth2/v1/certs",
    "client_x509_cert_url": CLIENT_X509_CERT_URL
}


class GDrive(DocumentProccesor):

    async def download(self, user_id: int, url: str) -> str:

        logger.info("Starting download from gdrive")

        file_id = self.__parse_url__(url)
        logger.info(f"File id recieved from url is {file_id}")

        try:

            logger.info("Trying to get access to gdrive details")
            credentials = service_account.Credentials.from_service_account_info(
                creds)
            drive_service = build('drive', 'v3', credentials=credentials)
            drive = drive_service.files()
            drive_file = drive.get(fileId=file_id,
                                   fields='name,size').execute()
            local_filename = drive_file['name']

            device: str = find_device(local_filename)
            official_devices = maintainer_details.get_devices(user_id)
            if not official_devices:
                official_devices = []
            if device not in official_devices and not maintainer_details.is_admin(
                    user_id):
                logger.info("This user is not a maintainer of this device")
                raise Exception("INVALID_DEVICE")

            logger.info(f"The file name to be downloaded is {local_filename}")
            content_size = int(drive_file['size'])
            logger.info(f"File size is {content_size}")

            download_request = drive.get_media(fileId=file_id)

            if not os.path.exists(TEMP_FOLDER_PATH):
                os.mkdir(TEMP_FOLDER_PATH)

            local_file = io.FileIO(TEMP_FOLDER_PATH + local_filename,
                                   'wb')  # this can be used to write to disk
            downloader = MediaIoBaseDownload(fd=local_file,
                                             request=download_request,
                                             chunksize=1024 * 1024 * 10)
            done = False
            logger.info("Starting downloading")
            while done is False:
                status, done = downloader.next_chunk()
                progress = status.progress() * content_size
                await progress_callback(progress, content_size, self.message,
                                        "Starting Download ....")
                logger.debug("Download %d%%." % int(status.progress() * 100))

            return local_filename

        except Exception as e:
            logger.error("Build failed to download as", e)
            await self.message.edit_text("Download failed")
        return None

    def __parse_url__(self, url: str) -> str:
        id: str = url.split("/")[5]
        if id.startswith("uc?id="):
            id = id.replace("uc?id=", "").replace("&export=download", "")
        return id
