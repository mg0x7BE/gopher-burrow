-- Inventory schema and reports, PostgreSQL flavor
CREATE TABLE IF NOT EXISTS items (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT           NOT NULL,
    category   VARCHAR(64)    NOT NULL DEFAULT 'misc',
    price      NUMERIC(10, 2) CHECK (price >= 0),
    created_at TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS orders (
    id         BIGSERIAL PRIMARY KEY,
    item_id    BIGINT  NOT NULL REFERENCES items (id),
    quantity   INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_items_category ON items (category);

/* Revenue per category for the last 30 days */
WITH recent AS (
    SELECT i.category, o.quantity * i.price AS revenue
    FROM orders o
    JOIN items i ON i.id = o.item_id
    WHERE o.created_at > now() - INTERVAL '30 days'
)
SELECT category,
       SUM(revenue)                          AS total,
       COUNT(*) FILTER (WHERE revenue > 100) AS big_orders,
       ROUND(AVG(revenue), 2)                AS average
FROM recent
GROUP BY category
HAVING SUM(revenue) > 1000
ORDER BY total DESC
LIMIT 10;

UPDATE items SET price = price * 0.9 WHERE category = 'toys' AND price IS NOT NULL;
