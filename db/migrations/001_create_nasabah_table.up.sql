-- db/migrations/001_create_nasabah_table.sql
CREATE TABLE nasabah (
    id SERIAL PRIMARY KEY,
    nik VARCHAR(16) UNIQUE NOT NULL,
    no_hp VARCHAR(15) UNIQUE NOT NULL,
    no_rekening VARCHAR(20) UNIQUE NOT NULL,
    saldo DECIMAL(15,2) DEFAULT 0 NOT NULL

);

CREATE TABLE tabungan (
    id SERIAL PRIMARY KEY,
    nasabah_id INT NOT NULL,
    nominal DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (nasabah_id) REFERENCES nasabah(id) ON DELETE CASCADE
);
