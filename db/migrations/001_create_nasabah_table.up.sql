-- db/migrations/001_create_nasabah_table.sql
CREATE TABLE nasabah (
    id SERIAL PRIMARY KEY,
    nik VARCHAR(16) UNIQUE NOT NULL,
    no_hp VARCHAR(15) UNIQUE NOT NULL,
    no_rekening VARCHAR(20) UNIQUE NOT NULL
);
