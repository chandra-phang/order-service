CREATE TABLE `orders` (
  `id` varchar(100) PRIMARY KEY,
  `user_id` varchar(100) NOT NULL,
  `status` varchar(50) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  INDEX (`user_id`, `status`)
);
