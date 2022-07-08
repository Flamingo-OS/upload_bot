# Parses the given file name to the folder it should go to
# It is assumed that file given through here will always be a KOSP file

from fileinput import filename
from typing import List
from bot.constants import BASE_URL

from bot.utils.logging import logger


def find_device(file_name: str) -> str:
    logger.info(f"Recieved request to find a device for {file_name}")

    parts: List[str] = file_name.split("-")[1:]
    device_name = parts[1]
    logger.info(f"It was parsed as {device_name}")
    return device_name


def parse_post_links(links: List[str]) -> dict:
    logger.info("Parsing links for post generation")

    parsed_links: dict = {
        "boot": [],
        "full": [],
        "fastboot": [],
        "incremental": []
    }

    for link in links:
        logger.info(f"Analysing {link}")
        link_split_number = -1

        if "sourceforge.net" in link:
            link_split_number = -2

        if BASE_URL in link or "sourceforge.net" in link:
            file_name = link.split("/")[link_split_number]
            logger.info(file_name)
            for item in parsed_links:
                if f"-{item}" in file_name:
                    parsed_links[item].append(link)
                    break

    logger.info(f"The links were parsed as {parsed_links}")
    return parsed_links


def find_kosp_ver(url: str) -> str:
    logger.info(f"Recieved request to find a Flamingo version for {url}")

    return url.split("/")[-1].split("-")[1]
