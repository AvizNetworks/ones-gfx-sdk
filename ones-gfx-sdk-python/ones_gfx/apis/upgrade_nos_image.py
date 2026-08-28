"""
POST /upgradeNOSImage  ->  /api/fm/upgradeNOSImage

NOS upgrade.
"""
from __future__ import annotations


from .._types import DeviceDetail
from ..client import Client

class ImageUpgradeDetailsItem(DeviceDetail):
    """Helper/ImageUpgradeDetails.java (extends DeviceDetail)"""

    pathToImage: str


def upgrade_nos_image(client: Client, *, items: list[ImageUpgradeDetailsItem]) -> bool:
    path = "upgradeNOSImage"
    return client.call_api("POST", path, json_body=items)
