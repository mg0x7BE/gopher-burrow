// Inventory shows the C++ elements the theme colors.
#include <algorithm>
#include <format>
#include <iostream>
#include <iterator>
#include <map>
#include <memory>
#include <optional>
#include <ranges>
#include <string>
#include <vector>

namespace inventory {

enum class Status { Draft, Published, Archived };

struct Item {
    long id;
    std::string name;
    double price;
    std::vector<std::string> tags;
    Status status = Status::Draft;
};

class Store {
public:
    virtual ~Store() = default;
    [[nodiscard]] virtual std::optional<Item> get(long id) const = 0;
    virtual bool put(Item item) = 0;
};

class MemStore final : public Store {
public:
    static constexpr std::size_t max_items = 1'000;

    [[nodiscard]] std::optional<Item> get(long id) const override {
        if (auto it = items_.find(id); it != items_.end()) {
            return it->second;
        }
        return std::nullopt;
    }

    bool put(Item item) override {
        if (items_.size() >= max_items) {
            return false;
        }
        items_.emplace(item.id, std::move(item));
        return true;
    }

private:
    std::map<long, Item> items_;
};

template <typename T, typename Pred>
std::vector<T> filter(const std::vector<T>& in, Pred keep) {
    std::vector<T> out;
    std::ranges::copy_if(in, std::back_inserter(out), keep);
    return out;
}

}  // namespace inventory

int main() {
    using namespace inventory;

    auto store = std::make_unique<MemStore>();
    Item item{42, "Gopher plush", 19.99, {"go", "theme"}, Status::Published};
    store->put(item);

    if (const auto found = store->get(7); !found) {
        std::cerr << std::format("missing item {}\n", 7);
    }

    const auto cheap = filter(std::vector{item}, [](const Item& i) { return i.price < 20.0; });
    auto names = cheap | std::views::transform(&Item::name);

    const char* raw = R"(C:\path\with\backslashes)";
    constexpr auto mask = 0xFFu;
    std::cout << std::format("{} {} {} {}\n", raw, mask, 'g', names.front()) << std::boolalpha
              << (item.status == Status::Archived) << '\n';
    return 0;
}
