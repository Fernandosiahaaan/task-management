CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY NOT NULL,
    title VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    status VARCHAR(50),
    due_date TIMESTAMP,
    assigned_to UUID NOT NULL,
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

