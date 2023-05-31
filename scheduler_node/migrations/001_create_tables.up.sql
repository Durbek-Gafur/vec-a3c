CREATE DATABASE IF NOT EXISTS app_db;

-- USE app_db;
-- Table for VEN info
CREATE TABLE ven_info (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    ram VARCHAR(50),
    core VARCHAR(50),
    url VARCHAR(255),
    max_queue_size VARCHAR(50),
    current_queue_size VARCHAR(50),
    preference_list VARCHAR(255),
    trust_score VARCHAR(50),
    max_queue_size_last_updated TIMESTAMP,
    current_queue_size_last_updated TIMESTAMP,
    trust_score_last_updated TIMESTAMP
);

-- Table for Workflow info
CREATE TABLE workflow_info (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    type VARCHAR(255),
    ram VARCHAR(50),
    core VARCHAR(50),
    policy VARCHAR(255),
    expected_execution_time VARCHAR(50),
    actual_execution_time VARCHAR(50),
    assigned_vm VARCHAR(255),
    assigned_at DATETIME,
    completed_at DATETIME,
    status VARCHAR(50),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
);
