-- Create a system user as fallback buyer for migrated transactions
INSERT INTO users (id, email, password_hash, name, role, is_active)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'system@migrated.local',
    '!',
    'System (Migrated)',
    'admin',
    FALSE
) ON CONFLICT (id) DO NOTHING;

-- Create a default store for each existing user
INSERT INTO stores (user_id, store_name, slug, description)
SELECT
    id,
    name || '''s Store',
    LOWER(REGEXP_REPLACE(name, '[^a-z0-9]+', '-', 'g')) || '-' || SUBSTRING(id::text, 1, 8),
    'Default store'
FROM users
WHERE id NOT IN (SELECT user_id FROM stores WHERE user_id IS NOT NULL)
  AND id != '00000000-0000-0000-0000-000000000001';

-- Create a default store for the system user
INSERT INTO stores (user_id, store_name, slug, description)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'System Store',
    'system-store',
    'Default store for migrated data'
) ON CONFLICT (user_id) DO NOTHING;

-- Assign all existing products to the first available store
UPDATE products
SET store_id = (
    SELECT s.id FROM stores s
    WHERE s.user_id IS NOT NULL
    ORDER BY s.created_at
    LIMIT 1
)
WHERE store_id IS NULL;

-- Migrate existing transactions to orders + order_items
DO $$
DECLARE
    v_buyer_id UUID;
    v_store_id UUID;
    v_order_id UUID;
    v_order_number TEXT;
    v_new_status TEXT;
    v_tx RECORD;
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'transactions') THEN
        FOR v_tx IN SELECT * FROM transactions LOOP
            -- Use system user as buyer (original transactions have no user info)
            v_buyer_id := '00000000-0000-0000-0000-000000000001';

            -- Use first available store
            SELECT s.id INTO v_store_id FROM stores s ORDER BY s.created_at LIMIT 1;

            -- Generate unique order number
            v_order_number := 'ORD-' || TO_CHAR(v_tx.created_at, 'YYYYMMDD') || '-' || LPAD(RIGHT(v_tx.id::text, 8), 8, '0');

            -- Map status
            v_new_status := CASE v_tx.status
                WHEN 'completed' THEN 'delivered'
                WHEN 'canceled' THEN 'canceled'
                ELSE 'pending'
            END;

            INSERT INTO orders (buyer_id, store_id, order_number, total_amount, status, shipping_address, created_at)
            VALUES (
                v_buyer_id,
                v_store_id,
                v_order_number,
                v_tx.total_amount,
                v_new_status,
                '{"street": "", "city": "", "state": "", "zip": "", "country": ""}'::jsonb,
                v_tx.created_at
            )
            RETURNING id INTO v_order_id;

            -- Create order_items from transaction_details
            INSERT INTO order_items (order_id, product_id, product_name, quantity, price, subtotal)
            SELECT
                v_order_id,
                td.product_id,
                td.product_name,
                td.quantity,
                td.price,
                td.subtotal
            FROM transaction_details td
            WHERE td.transaction_id = v_tx.id;
        END LOOP;
    END IF;
END $$;
