package inventory;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Predicate;

/** Inventory shows the Java elements the theme colors. */
public final class Inventory {

    public enum Status { DRAFT, PUBLISHED, ARCHIVED }

    public record Item(long id, String name, double price, List<String> tags, Status status) {
        public boolean isCheap() {
            return price < 20.0;
        }
    }

    public sealed interface Result permits Found, Missing {}

    public record Found(Item item) implements Result {}

    public record Missing(long id) implements Result {}

    public interface Store {
        Result get(long id);

        void put(Item item);

        default int size() {
            return 0;
        }
    }

    static final int MAX_ITEMS = 1_000;
    private static final List<String> DEFAULT_TAGS = List.of("go", "theme");

    static class MemStore implements Store {
        private final Map<Long, Item> items = new ConcurrentHashMap<>();

        @Override
        public Result get(long id) {
            var item = items.get(id);
            return item != null ? new Found(item) : new Missing(id);
        }

        @Override
        public void put(Item item) {
            if (items.size() >= MAX_ITEMS) {
                throw new IllegalStateException("store is full");
            }
            items.put(item.id(), item);
        }

        @Override
        public int size() {
            return items.size();
        }
    }

    static <T> List<T> filter(List<T> in, Predicate<? super T> keep) {
        var out = new ArrayList<T>();
        for (T v : in) {
            if (keep.test(v)) {
                out.add(v);
            }
        }
        return out;
    }

    public static void main(String[] args) throws InterruptedException {
        Store store = new MemStore();
        var item = new Item(42, "Gopher plush", 19.99, DEFAULT_TAGS, Status.PUBLISHED);
        store.put(item);

        String label = switch (store.get(7)) {
            case Found(Item found) -> "found " + found.name();
            case Missing(long id) -> "missing item " + id;
        };
        System.out.println(label);

        var thread = Thread.ofVirtual().start(() -> System.out.println(item.name().toUpperCase()));
        thread.join();

        String json = """
                {"name": "gopher", "ok": true}
                """;
        char letter = 'g';
        int mask = 0xFF;
        Optional<String> first = filter(List.of(item), Item::isCheap).stream().map(Item::name).findFirst();
        System.out.printf("%s %d %c %s %s %d%n", json.strip(), mask, letter, first.orElse("none"), Status.DRAFT, store.size());
    }
}
