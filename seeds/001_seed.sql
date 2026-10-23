
INSERT INTO payments (
    user_id,
    customer_number,
    customer_name,
    amount,
    method,
    status
)
SELECT
    id,
    'CUST-DEMO-001',
    'Budi Santoso',
    150000,
    'transfer',
    'success'
FROM users
WHERE email = 'demo@example.com'
  AND NOT EXISTS (
      SELECT 1
      FROM payments
      WHERE user_id = users.id
        AND customer_number = 'CUST-DEMO-001'
  );

INSERT INTO payments (
    user_id,
    customer_number,
    customer_name,
    amount,
    method,
    status
)
SELECT
    id,
    'CUST-DEMO-002',
    'Siti Aminah',
    75000,
    'ewallet',
    'pending'
FROM users
WHERE email = 'demo@example.com'
  AND NOT EXISTS (
      SELECT 1
      FROM payments
      WHERE user_id = users.id
        AND customer_number = 'CUST-DEMO-002'
  );

INSERT INTO payments (
    user_id,
    customer_number,
    customer_name,
    amount,
    method,
    status
)
SELECT
    id,
    'CUST-DEMO-003',
    'Andi Pratama',
    250000,
    'va',
    'failed'
FROM users
WHERE email = 'demo@example.com'
  AND NOT EXISTS (
      SELECT 1
      FROM payments
      WHERE user_id = users.id
        AND customer_number = 'CUST-DEMO-003'
  );
