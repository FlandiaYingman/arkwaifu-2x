from dataclasses import dataclass
from typing import Any, List


@dataclass
class Variant:
    artID: str
    variation: str
    contentPresent: bool
    contentWidth: int | None
    contentHeight: int | None

    @staticmethod
    def from_dict(obj: Any) -> 'Variant':
        _artID = str(obj.get("artID"))
        _variation = str(obj.get("variation"))
        _contentPresent = bool(obj.get("contentPresent"))
        _contentWidth = int(obj.get("contentWidth")) if obj.get("contentWidth") else None
        _contentHeight = int(obj.get("contentHeight")) if obj.get("contentHeight") else None
        return Variant(_artID, _variation, _contentPresent, _contentWidth, _contentHeight)


@dataclass
class Art:
    id: str
    category: str
    variants: List[Variant]

    @staticmethod
    def from_dict(obj: Any) -> 'Art':
        _id = str(obj.get("id"))
        _category = str(obj.get("category"))
        _variants = [Variant.from_dict(y) for y in obj.get("variants")]
        return Art(_id, _category, _variants)
