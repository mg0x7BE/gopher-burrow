// Inventory shows the Kotlin elements the theme colors.
package inventory

import java.util.concurrent.ConcurrentHashMap

const val MAX_ITEMS = 1_000
val defaultTags = listOf("go", "theme")

/** The state of an item. */
enum class Status { DRAFT, PUBLISHED, ARCHIVED }

data class Item(
    val id: Long,
    val name: String,
    val price: Double,
    val tags: List<String> = defaultTags,
    var status: Status = Status.DRAFT,
) {
    val isCheap: Boolean get() = price < 20.0
}

sealed interface Lookup {
    data class Found(val item: Item) : Lookup
    data class Missing(val id: Long) : Lookup
}

interface Store {
    fun get(id: Long): Lookup
    fun put(item: Item)
}

class MemStore : Store {
    private val items = ConcurrentHashMap<Long, Item>()

    override fun get(id: Long): Lookup =
        items[id]?.let { Lookup.Found(it) } ?: Lookup.Missing(id)

    override fun put(item: Item) {
        check(items.size < MAX_ITEMS) { "store is full" }
        items[item.id] = item
    }

    companion object {
        fun withDefaults(): MemStore = MemStore()
    }
}

fun <T> List<T>.keep(predicate: (T) -> Boolean): List<T> = filter(predicate)

@Suppress("MagicNumber")
fun main() {
    val store = MemStore.withDefaults()
    val item = Item(42, "Gopher plush", 19.99, status = Status.PUBLISHED)
    store.put(item)

    val label = when (val result = store.get(7)) {
        is Lookup.Found -> "found ${result.item.name}"
        is Lookup.Missing -> "missing item ${result.id}"
    }
    println(label)

    var count = 0
    for (name in listOf("alpha", "beta")) {
        count += name.length
    }

    val raw = """C:\path\with\backslashes"""
    val cheap = listOf(item).keep { it.isCheap }.map(Item::name)
    println("$raw ${cheap.size} $count ${0xFF} ${'g'} ${Status.ARCHIVED} ${item.tags.joinToString()}")
}
