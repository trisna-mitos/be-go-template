INSERT INTO products (id, name, description, price, stock, created_at, updated_at)
VALUES 
    ('550e8400-e29b-41d4-a716-446655440000', 'Test Product 1', 'Description 1', 29.99, 100, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('550e8400-e29b-41d4-a716-446655440001', 'Test Product 2', 'Description 2', 39.99, 50, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO dipan_types (nama_type)
VALUES
    ('Dipan Minimalis'),
    ('Dipan Modern'),
    ('Dipan Kayu Jati');