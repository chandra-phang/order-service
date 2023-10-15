CREATE TABLE `cart_products` (
  `id` varchar(100) PRIMARY KEY,
  `user_id` varchar(100) NOT NULL,
  `product_id` varchar(100) NOT NULL,
  `quantity` int,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  INDEX (`user_id`)
);
