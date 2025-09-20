-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    type ENUM('income', 'expense') NOT NULL,
    color CHAR(7) DEFAULT '#007BFF',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_name_type (name, type)
);

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    type ENUM('income', 'expense') NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    category_id BIGINT NOT NULL,
    transaction_date DATE NOT NULL,
    memo TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_transaction_date (transaction_date),
    INDEX idx_category_id (category_id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Create budgets table
CREATE TABLE IF NOT EXISTS budgets (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    category_id BIGINT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    target_year INT NOT NULL,
    target_month TINYINT NOT NULL CHECK (target_month BETWEEN 1 AND 12),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_budget_period (category_id, target_year, target_month),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Insert default categories
INSERT IGNORE INTO categories (name, type, color) VALUES
-- Income categories
('給与', 'income', '#28a745'),
('副業', 'income', '#17a2b8'),
('その他収入', 'income', '#6f42c1'),
-- Expense categories  
('食費', 'expense', '#dc3545'),
('住居費', 'expense', '#fd7e14'),
('交通費', 'expense', '#20c997'),
('光熱費', 'expense', '#ffc107'),
('通信費', 'expense', '#6610f2'),
('娯楽費', 'expense', '#e83e8c'),
('その他支出', 'expense', '#6c757d');