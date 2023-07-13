from __future__ import annotations

import os
import urllib.parse
from typing import Final
from typing import TYPE_CHECKING

import requests

from arkwaifu_2x.arkwaifu.models import Art

if TYPE_CHECKING:
    from _typeshed import StrOrBytesPath

API_URL: Final[str] = "https://arkwaifu.cc/api/v1"
USER_TOKEN: Final[str] = os.environ["USER_TOKEN"]


def get_arts_by_absent_variation(absent_variation: str) -> list[Art]:
    response = requests.get(
        f"{API_URL}/arts",
        params={"absent-variation": absent_variation},
    )
    return [Art.from_dict(x) for x in response.json()]


def put_variant(id: str, variation: str):
    requests.put(
        f"{API_URL}/arts/{urllib.parse.quote(id)}/variants/{urllib.parse.quote(variation)}",
        params={'user': USER_TOKEN},
        data={'artID': id, 'variation': variation}
    ).raise_for_status()


def get_content(id: str, variation: str) -> bytes:
    response = requests.get(
        f"{API_URL}/arts/{urllib.parse.quote(id)}/variants/{urllib.parse.quote(variation)}/content"
    )
    response.raise_for_status()
    return response.content


def put_content(id: str, variation: str, content_path: StrOrBytesPath):
    with open(content_path, 'rb') as content_file:
        requests.put(
            f"{API_URL}/arts/{urllib.parse.quote(id)}/variants/{urllib.parse.quote(variation)}/content",
            params={'user': USER_TOKEN},
            data=content_file,
        ).raise_for_status()
