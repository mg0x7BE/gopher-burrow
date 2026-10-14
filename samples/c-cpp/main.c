/* Inventory shows the C elements the theme colors. */
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX_ITEMS 1000
#define ARRAY_LEN(a) (sizeof(a) / sizeof((a)[0]))

typedef enum { DRAFT, PUBLISHED, ARCHIVED } status_t;

typedef struct {
    long id;
    char name[32];
    double price;
    status_t status;
} item_t;

typedef struct {
    item_t items[MAX_ITEMS];
    size_t count;
} store_t;

static const char *default_tags[] = {"go", "theme"};

static bool store_put(store_t *s, const item_t *item) {
    if (s->count >= MAX_ITEMS) {
        return false;
    }
    s->items[s->count++] = *item;
    return true;
}

static item_t *store_get(store_t *s, long id) {
    for (size_t i = 0; i < s->count; i++) {
        if (s->items[i].id == id) {
            return &s->items[i];
        }
    }
    return NULL;
}

int main(void) {
    store_t *store = calloc(1, sizeof *store);
    if (store == NULL) {
        return EXIT_FAILURE;
    }

    item_t item = {.id = 42, .price = 19.99, .status = PUBLISHED};
    strncpy(item.name, "Gopher plush", sizeof item.name - 1);
    store_put(store, &item);

    const item_t *found = store_get(store, 7);
    if (!found) {
        fprintf(stderr, "missing item %d\n", 7);
    }

    for (size_t i = 0; i < ARRAY_LEN(default_tags); i++) {
        printf("tag %zu: %s\n", i, default_tags[i]);
    }

    unsigned mask = 0xFFu;
    char letter = 'g';
    printf("%s %.2f %u %c %d %d\n", item.name, item.price, mask, letter, item.status == ARCHIVED, DRAFT);

    free(store);
    return 0;
}
