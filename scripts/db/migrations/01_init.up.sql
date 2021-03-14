
CREATE DATABASE eshop;

CREATE TABLE eshop.inventory (
    sku STRING PRIMARY KEY,
    name STRING NOT NULL, 
    price DECIMAL (10, 2),
    quantity INT NOT NULL DEFAULT 0
);

CREATE TABLE eshop.promotions (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    sku STRING NOT NULL,
    active BOOL NOT NULL,
    promotions JSONB,
    CONSTRAINT "primary" PRIMARY KEY (sku ASC, id ASC)
);

CREATE TABLE eshop.carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart JSONB NOT NULL,
    settled JSONB
);