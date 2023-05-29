CREATE TABLE IF NOT EXISTS "user" (
    "id" UUID PRIMARY KEY,
    "first_name" VARCHAR(30) NOT NULL,
    "last_name" VARCHAR(30) NOT NULL,
    "phone_number" VARCHAR(17) NOT NULL UNIQUE,
    "date_of_birth" DATE,
    "password" VARCHAR(30) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
); 