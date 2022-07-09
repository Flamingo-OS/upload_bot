# Thanks to sukionotes for parts of this code

import time
from datetime import timedelta
from pyrogram.errors import MessageNotModified
from bot.utils.messaging import edit_message

progress_callback_data = dict()


# https://stackoverflow.com/a/49361727
def format_bytes(size: int) -> str:
    size = int(size)
    power = 1024
    n = 0
    power_labels = {0: '', 1: 'K', 2: 'M', 3: 'G', 4: 'T'}
    while size > power:
        size /= power
        n += 1
    return f"{size:.2f} {power_labels[n]+'B'}"


# https://stackoverflow.com/a/34325723
def return_progress_string(current: int, total: int) -> str:
    filled_length = int(30 * current // total)
    return '[' + '=' * filled_length + ">" + ' ' * (30 - filled_length) + ']'


def calculate_eta(current: int, total: int, start_time: float) -> str:
    if not current:
        return '00:00:00'
    end_time: float = time.time()
    elapsed_time: float = end_time - start_time
    eta_seconds: float = (elapsed_time * (total / current)) - elapsed_time
    eta: str = ''.join(str(
        timedelta(seconds=eta_seconds)).split('.')[:-1]).split(', ')
    eta[-1] = eta[-1].rjust(8, '0')
    return ', '.join(eta)


def calculate_speed(current: int, start_time: float):
    if not current:
        return '0 B'
    end_time: float = time.time()
    elapsed_time = end_time - start_time
    speed = format_bytes(current / elapsed_time)
    return speed


async def progress_callback(
    current: int,
    total: int,
    message,
    text: str,
    upload: bool = False,
):
    """
    current: int, defines the amount of progress completed. can be in bytes
    total: int, defines the total amount of progress/parts to complete. can be in bytes 
    """

    message_identifier = (message.chat.id, message.id)
    last_edit_time, prevtext, start_time = progress_callback_data.get(
        message_identifier, (0, None, time.time()))
    if current == total:
        try:
            progress_callback_data.pop(message_identifier)
        except KeyError:
            pass
    elif (time.time() - last_edit_time) > 1:
        handle = 'Upload' if upload else 'Download'
        if last_edit_time:
            speed = calculate_speed(current, start_time)
        else:
            speed = '0 B'
        text = f'''{text}
<code>{return_progress_string(current, total)}</code>
<b>Total Size:</b> {format_bytes(total)}
<b>{handle}ed Size:</b> {format_bytes(current)}
<b>{handle} Speed:</b> {speed}/s
<b>ETA:</b> {calculate_eta(current, total, start_time)}'''
        if prevtext != text:
            try:
                await edit_message(message, text)
                prevtext = text
                last_edit_time = time.time()
                progress_callback_data[
                    message_identifier] = last_edit_time, prevtext, start_time
            except (MessageNotModified):
                pass
