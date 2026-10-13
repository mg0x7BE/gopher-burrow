// Fetches items and renders them into the list
type Status = "draft" | "published";

interface Item {
  id: number;
  name: string;
  price?: number;
  status: Status;
  readonly tags: string[];
}

const API_URL = "/items";

export class ItemList {
  private items: Item[] = [];

  constructor(private readonly root: HTMLElement) {}

  async load(limit = 20): Promise<void> {
    const res = await fetch(`${API_URL}?limit=${limit}`);
    if (!res.ok) {
      throw new Error(`request failed: ${res.status}`);
    }
    this.items = (await res.json()) as Item[];
    this.render();
  }

  private render(): void {
    this.root.innerHTML = this.items
      .filter((item) => item.status === "published" && (item.price ?? 0) > 0)
      .map(({ name, tags }) => `<li>${name} <small>${tags.join(", ")}</small></li>`)
      .join("");
  }
}

const list = new ItemList(document.getElementById("items")!);
list.load().catch(console.error);
