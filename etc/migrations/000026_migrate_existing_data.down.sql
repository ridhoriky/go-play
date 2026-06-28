-- Rollback data migration
-- WARNING: This will delete ALL migrated orders, order_items, and system-generated stores/users.
-- Review carefully before running.

-- Delete order_items for migrated orders
DELETE FROM order_items
WHERE order_id IN (
    SELECT id FROM orders WHERE order_number LIKE 'ORD-%'
);

-- Delete migrated orders
DELETE FROM orders WHERE order_number LIKE 'ORD-%';

-- Reset store_id on products that were assigned during migration
UPDATE products SET store_id = NULL
WHERE store_id IN (
    SELECT id FROM stores WHERE description = 'Default store'
       OR description = 'Default store for migrated data'
);

-- Delete system-generated stores
DELETE FROM stores
WHERE description = 'Default store'
   OR description = 'Default store for migrated data';

-- Delete system user
DELETE FROM users
WHERE id = '00000000-0000-0000-0000-000000000001';
