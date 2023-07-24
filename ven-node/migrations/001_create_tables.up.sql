CREATE DATABASE IF NOT EXISTS app_db;

-- USE app_db;

CREATE TABLE IF NOT EXISTS queue (
  id INT AUTO_INCREMENT PRIMARY KEY,
  workflow_id INT,
  status ENUM('pending', 'processing', 'done') DEFAULT 'pending',
  enqueued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS workflow (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(255) NOT NULL,
  duration INT NOT NULL,
  received_at DATETIME NOT NULL,
  started_execution_at DATETIME NULL,
  completed_at DATETIME NULL
);

CREATE TABLE IF NOT EXISTS queue_size (
  id INT AUTO_INCREMENT PRIMARY KEY,
  size INTEGER NOT NULL
);

