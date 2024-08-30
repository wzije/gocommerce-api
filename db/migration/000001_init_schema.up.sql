CREATE TABLE IF NOT EXISTS users
(
    id         BIGSERIAL PRIMARY KEY,
    username   VARCHAR(100) UNIQUE NOT NULL,
    email      VARCHAR(100) UNIQUE NOT NULL,
    password   VARCHAR(100)        NOT NULL,
    role       VARCHAR(20) DEFAULT 'CUSTOMER', -- 1. OWNER, 2. ADMIN, 3. CUSTOMER
    created_at timestamptz DEFAULT (now()),
    updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS profiles
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT       NOT NULL,
    name        VARCHAR(255) NOT NULL,
    phone       VARCHAR(50),
    address     TEXT,
    city        VARCHAR(255),
    state       VARCHAR(255),
    postal_code VARCHAR(20),
    country     VARCHAR(255),
    created_at  timestamptz DEFAULT (now()),
    updated_at  timestamptz,
    CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE shops
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT       NOT NULL, --owner / pemilik toko
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    address     TEXT,
    phone       VARCHAR(50),
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- 2. Create Products Table
CREATE TABLE products
(
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255)        NOT NULL,
    description TEXT,
    price       DECIMAL(10, 2)      NOT NULL,
    sku         VARCHAR(100) UNIQUE NOT NULL,
    image_url   TEXT,
    shop_id     BIGINT              NOT NULL,
    created_at  timestamptz DEFAULT (now()),
    updated_at  timestamptz,
    FOREIGN KEY (shop_id) REFERENCES shops (ID) ON DELETE CASCADE
);


-- 4. Create Orders Table
CREATE TABLE orders
(
    id               BIGSERIAL PRIMARY KEY,
    date             TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status           VARCHAR(50)    NOT NULL,
    amount           DECIMAL(10, 2) NOT NULL,
    shipping_date    TIMESTAMP,
    shipping_address TEXT           NOT NULL,
    user_id          BIGINT         NOT NULL,
    shop_id          BIGINT         NOT NULL,
    created_at       timestamptz             DEFAULT (now()),
    updated_at       timestamptz,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (shop_id) REFERENCES shops (id) ON DELETE CASCADE
);

CREATE TABLE payments
(
    id             BIGSERIAL PRIMARY KEY,
    order_id       BIGINT         NOT NULL,
    payment_method VARCHAR(100)   NOT NULL,
    amount         DECIMAL(10, 2) NOT NULL,
    status         VARCHAR(50)    NOT NULL, -- Contoh status: 'Pending', 'Completed', 'Failed'
    created_at     TIMESTAMPTZ DEFAULT now(),
    updated_at     TIMESTAMPTZ,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);

-- 5. Create OrderDetails Table
CREATE TABLE order_details
(
    id             BIGSERIAL PRIMARY KEY,
    order_id       BIGINT         NOT NULL,
    product_id     BIGINT         NOT NULL,
    quantity       INTEGER        NOT NULL CHECK (quantity > 0),
    price_per_unit DECIMAL(10, 2) NOT NULL,
    total_price    DECIMAL(10, 2) GENERATED ALWAYS AS (quantity * price_per_unit) STORED,
    created_at     timestamptz DEFAULT (now()),
    updated_at     timestamptz,
    FOREIGN KEY (order_id) REFERENCES Orders (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES Products (id) ON DELETE CASCADE
);

-- 6. Create Warehouses Table
CREATE TABLE warehouses
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    location   TEXT,
    is_active  BOOLEAN     DEFAULT true,
    shop_id    BIGINT       NOT NULL,
    user_id    BIGINT       NOT NULL, -- pemilik gudang
    created_at timestamptz DEFAULT (now()),
    updated_at timestamptz
);

-- 7. Create WarehouseInventory Table ALIAS STOCK
CREATE TABLE warehouse_inventories
(
    id           BIGSERIAL PRIMARY KEY,
    product_id   BIGINT NOT NULL,
    warehouse_id BIGINT NOT NULL,
    quantity     BIGINT NOT NULL CHECK (quantity >= 0),
    created_at   TIMESTAMPTZ DEFAULT now(),
    updated_at   TIMESTAMPTZ,
    FOREIGN KEY (product_id) REFERENCES Products (id) ON DELETE CASCADE,
    FOREIGN KEY (warehouse_id) REFERENCES Warehouses (id) ON DELETE CASCADE,
    UNIQUE (product_id, warehouse_id) -- Ensure unique product-warehouse combinations
);

CREATE TABLE product_transfer_warehouses
(
    id                       BIGSERIAL PRIMARY KEY,
    source_warehouse_id      BIGINT      NOT NULL,
    destination_warehouse_id BIGINT      NOT NULL,
    product_id               BIGINT      NOT NULL,
    quantity                 BIGINT      NOT NULL,
    status                   VARCHAR(20) NOT NULL CHECK (status IN ('PENDING', 'COMPLETED', 'FAILED')),
    created_at               TIMESTAMPTZ DEFAULT now(),
    updated_at               TIMESTAMPTZ,
    FOREIGN KEY (product_id) REFERENCES Products (id) ON DELETE CASCADE,
    FOREIGN KEY (source_warehouse_id) REFERENCES Warehouses (id) ON DELETE CASCADE,
    FOREIGN KEY (destination_warehouse_id) REFERENCES Warehouses (id) ON DELETE CASCADE
);

CREATE TABLE stock_locks
(
    id           BIGSERIAL PRIMARY KEY,
    order_id     BIGINT      NOT NULL,
    product_id   BIGINT      NOT NULL,
    warehouse_id BIGINT      NOT NULL,
    quantity     INTEGER     NOT NULL CHECK (quantity > 0),
    locked_at    TIMESTAMPTZ DEFAULT now(),
    released_at  TIMESTAMPTZ,
    status       VARCHAR(20) NOT NULL CHECK (status IN ('LOCKED', 'RELEASED')),
    created_at   TIMESTAMPTZ DEFAULT now(),
    updated_at   TIMESTAMPTZ,
    FOREIGN KEY (order_id) REFERENCES orders (id),
    FOREIGN KEY (product_id) REFERENCES products (id),
    FOREIGN KEY (warehouse_id) REFERENCES warehouses (id)
);

CREATE TABLE order_warehouse_allocations
(
    id           BIGSERIAL PRIMARY KEY,
    order_id     BIGINT  NOT NULL,
    warehouse_id BIGINT  NOT NULL,
    product_id   BIGINT  NOT NULL,
    quantity     INTEGER NOT NULL CHECK (quantity > 0),
    created_at   TIMESTAMPTZ DEFAULT now(),
    updated_at   TIMESTAMPTZ,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    FOREIGN KEY (warehouse_id) REFERENCES warehouses (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

-- Optional: Create Indexes for better performance on frequently queried columns

CREATE INDEX idx_profile_user ON profiles (user_id);

CREATE INDEX idx_shop_user ON shops (user_id);

CREATE INDEX idx_order_customer ON orders (user_id);

CREATE INDEX idx_order_details_order ON order_details (order_id);
CREATE INDEX idx_order_details_product ON order_details (product_id);

CREATE INDEX idx_warehouse_shop ON warehouses (shop_id);
CREATE INDEX idx_warehouse_user ON warehouses (user_id);

CREATE INDEX idx_warehouse_inventories_product ON warehouse_inventories (product_id);
CREATE INDEX idx_warehouse_inventories_warehouse ON warehouse_inventories (warehouse_id);

CREATE INDEX idx_stock_lock_warehouse ON stock_locks (warehouse_id);
CREATE INDEX idx_stock_lock_order ON stock_locks (order_id);
CREATE INDEX idx_stock_lock_product ON stock_locks (product_id);

CREATE INDEX idx_order_warehouse_allocation_warehouse ON order_warehouse_allocations (warehouse_id);
CREATE INDEX idx_order_warehouse_allocation_order ON order_warehouse_allocations (order_id);
CREATE INDEX idx_order_warehouse_allocation_product ON order_warehouse_allocations (product_id);

CREATE INDEX idx_product_transfer_warehouses_product ON product_transfer_warehouses (product_id);
CREATE INDEX idx_product_transfer_warehouses_source_warehouse ON product_transfer_warehouses (source_warehouse_id);
CREATE INDEX idx_product_transfer_warehouses_destination_warehouse ON product_transfer_warehouses (destination_warehouse_id);