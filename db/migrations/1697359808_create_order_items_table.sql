CREATE TABLE `order_items` (
  `id` varchar(100) PRIMARY KEY,
  `order_id` varchar(100) NOT NULL,
  `product_id` varchar(100) NOT NULL,
  `quantity` int,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`order_id`) REFERENCES orders(`id`),
  INDEX (`order_id`)
);
