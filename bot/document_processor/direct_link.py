from bot.document_processor.base import DocumentProccesor
import requests
import os


class DirectLink(DocumentProccesor):

    def download(self, url: str) -> str:

        temp_folder_path: str = "./DumpsterFire/"

        if not os.path.exists(temp_folder_path):
            os.mkdir(temp_folder_path)

        local_filename: str = temp_folder_path + url.split('/')[-1]
        res = requests.get(url, stream=True, allow_redirects=True)
        content_length = int(res.headers['Content-Length'])

        data = b''

        for chunk in res.iter_content(chunk_size=1024 * 1024 * 10):
            if (chunk):
                data += chunk
            print(len(data) / content_length)

        open(local_filename, 'wb').write(data)
        return local_filename