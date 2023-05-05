-- Up
CREATE TABLE IF NOT EXISTS workflow (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(255) NOT NULL,
  duration INT NOT NULL,
  received_at DATETIME NOT NULL,
  started_execution_at DATETIME NULL,
  completed_at DATETIME NULL
);

-- Down
DROP TABLE IF EXISTS workflow;
