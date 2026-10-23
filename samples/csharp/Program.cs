using System.Text.Json;
using System.Text.Json.Serialization;

namespace Inventory;

/// <summary>The state of an item.</summary>
public enum Status
{
    Draft,
    Published,
    Archived,
}

public record Item(long Id, string Name, decimal Price, IReadOnlyList<string> Tags)
{
    [JsonIgnore]
    public Status Status { get; init; } = Status.Draft;

    public bool IsCheap => Price < 20m;
}

public interface IStore
{
    Task<Item?> GetAsync(long id, CancellationToken ct = default);
    void Put(Item item);
}

public sealed class MemoryStore : IStore
{
    private const int MaxItems = 1_000;
    private static readonly string[] DefaultTags = ["go", "theme"];

    private readonly Dictionary<long, Item> _items = new();
    private readonly Lock _lock = new();

    public event Action<string>? Changed;

    public async Task<Item?> GetAsync(long id, CancellationToken ct = default)
    {
        await Task.Delay(10, ct);
        lock (_lock)
        {
            return _items.TryGetValue(id, out var item) ? item : null;
        }
    }

    public void Put(Item item)
    {
        lock (_lock)
        {
            if (_items.Count >= MaxItems)
                throw new InvalidOperationException("store is full");
            _items[item.Id] = item with { Tags = [..item.Tags, ..DefaultTags] };
        }
        Changed?.Invoke($"stored {item.Name}");
    }
}

public static class Extensions
{
    public static IEnumerable<T> Cheap<T>(this IEnumerable<T> source, Func<T, decimal> price) =>
        source.Where(x => price(x) < 20m);
}

public static class Program
{
    public static async Task Main()
    {
        var store = new MemoryStore();
        store.Changed += msg => Console.WriteLine($"notify: {msg}");

        store.Put(new Item(42, "Gopher plush", 19.99m, ["toy"]) { Status = Status.Published });

        var missing = await store.GetAsync(7);
        Console.WriteLine(missing is null ? "missing item 7" : missing.Name);

        var items = new List<Item> { new(1, "Sticker", 2.5m, []), new(2, "Mug", 24m, ["kitchen"]) };
        var cheap = items.Cheap(i => i.Price).Select(i => i.Name.ToUpperInvariant()).ToList();

        var label = items[0].Status switch
        {
            Status.Draft => "draft",
            Status.Published or Status.Archived => "visible",
            _ => throw new ArgumentOutOfRangeException(),
        };

        const string path = @"C:\path\with\backslashes";
        // language=json
        var json = """{"name": "gopher", "ok": true}""";
        Console.WriteLine($"{cheap.Count} cheap, {label}, {path}, {json.Length}, {0xFF}, {'g'}, {items[1].IsCheap}");
        Console.WriteLine(JsonSerializer.Serialize(items[1]));
    }
}
