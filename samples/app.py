"""Small inventory helpers used by the reporting job."""
from dataclasses import dataclass, field
from typing import Iterable

MAX_ITEMS = 1_000


@dataclass
class Item:
    id: int
    name: str
    price: float = 0.0
    tags: list[str] = field(default_factory=list)

    @property
    def is_cheap(self) -> bool:
        return self.price < 20


def total(items: Iterable[Item], discount: float = 0.1) -> float:
    """Return the discounted sum of all prices."""
    subtotal = sum(item.price for item in items)
    return round(subtotal * (1 - discount), 2)


if __name__ == "__main__":
    stock = [Item(1, "Gopher plush", 19.99, ["toy"]), Item(2, "Sticker", 2.5)]
    cheap = [i.name for i in stock if i.is_cheap and len(stock) < MAX_ITEMS]
    print(f"{len(cheap)} cheap items, total {total(stock):.2f}", end="\n")
