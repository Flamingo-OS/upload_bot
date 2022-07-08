from typing import List
from bot.database import init
from bot.utils.logging import logger


class MaintainerDetails:

    def __init__(self):
        logger.info("Initializing Maintainer Details")
        self.maintainer_db = init()["maintainer"]

    # Check if user is a maintainer or not
    def is_maintainer_or_admin(self, user_id: int, is_admin: bool = False) -> bool:
        logger.info(f"Checking if user {user_id} is a maintainer")
        get_user = self.maintainer_db.find_one({"user_id": user_id})

        if not get_user:
            logger.info(f"User {user_id} not found in maintainers table")
            return False

        elif get_user["is_maintainer" if not is_admin else "is_admin"] or get_user["is_admin"]:
            logger.info(f"user {user_id} is a maintainer or an admin")
            return True

        logger.info(f"user {user_id} is a not a maintainer")
        return False

    # Add a maintainer
    def add_maintainer(self, requested_id: int, user_id: int, name: str,
                       device: str) -> bool:
        logger.info(
            f"Adding user {user_id} to maintainer database since {requested_id} requested it"
        )

        if not self.is_maintainer_or_admin(requested_id, True):
            logger.info(f"User {requested_id} is not an admin")
            return False

        if self.is_maintainer_or_admin(user_id):
            logger.info(f"User {user_id} is already a maintainer")
            if device in self.get_devices(user_id):
                logger.info(f"User {user_id} already has {device}")
                return False
            else:
                logger.info(f"Adding {device} to user {user_id}")
                self.maintainer_db.update_one({"user_id": user_id},
                                              {"$push": {
                                                  "device": device
                                              }})
                return True
        else:
            self.maintainer_db.insert_one({
                "user_id": user_id,
                "name": name,
                "is_maintainer": True,
                "is_admin": False,
                "device": [device],
                "support_group": ""
            })
            logger.info("Added the user to a maintainer")
            return True

    # Remove a maintainer
    def remove_maintainer(self, requester_id: int, user_id: int) -> bool:
        logger.info(f"Removing user {user_id} from maintainer database")

        if not self.is_maintainer_or_admin(requester_id):
            logger.info("The user is not an admin")
            return False

        if not self.is_maintainer_or_admin(user_id):
            logger.info(f"User {user_id} is not a maintainer")
            return False
        else:
            self.maintainer_db.delete_one({"user_id": user_id})
            return True

    # Add an admin
    def add_admin(self, requester_id: int, user_id: int) -> bool:
        logger.info(
            f"Adding user {user_id} to admin database since {requester_id} requested it"
        )

        if not self.is_maintainer_or_admin(requester_id, True):
            logger.info(f"User {requester_id} is not an admin")
            return False

        logger.info("Promoting the user to admin status")
        self.maintainer_db.update_one({"user_id": user_id},
                                      {"$set": {
                                          "is_admin": True
                                      }})
        return True

    # Get devicrs
    def get_devices(self, user_id: int) -> List[dict]:
        logger.info(f"Getting devices for user {user_id}")
        get_user = self.maintainer_db.find_one({"user_id": user_id})

        # User is new to the database
        if get_user is None:
            logger.info(f"user {user_id} not found in maintainer database")
            return False
        logger.info(f"Found devices as {get_user['device']}")
        return get_user["device"]

    def get_maintainers(self, device: str):
        logger.info(f"Getting maintainers for device {device}")

        try:
            get_device = self.maintainer_db.find({"device": device})
            if not get_device:
                logger.info("This device isn't officially maintainer")
                return False

            maintainers: List[dict] = []

            for maintainer in get_device:
                logger.info(
                    f"{maintainer['name']} is a maintainer for {device}")
                maintainers.append({
                    "name": maintainer["name"],
                    "user_id": maintainer["user_id"]
                })
            return maintainers
        except:
            logger.info("Something went wrong with the db")
            return

    def get_device_support_group(self, device: str) -> str:
        logger.info(f"Getting support group for device {device}")
        try:
            get_device = self.maintainer_db.find_one({"device": device})

            # User is new to the database
            if not get_device:
                logger.info(
                    f"device {device} not found in maintainer database")
                return False
            return get_device["support_group"]
        except:
            logger.info("Something went wrong with the db")
            return

    def add_support_group(self, requester_id: int, user_id: int,
                          support_group: str) -> bool:
        logger.info(
            f"Adding support group for user {user_id} since {requester_id} requested it"
        )

        if not self.is_maintainer_or_admin(requester_id):
            logger.info(f"User {requester_id} is not a maintainer or admin")
            return False

        logger.info("Adding support group")
        self.maintainer_db.update_one(
            {"user_id": user_id}, {"$set": {
                "support_group": support_group
            }})
        return True


maintainer_details = MaintainerDetails()
