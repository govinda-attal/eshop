INSERT INTO eshop.inventory (sku, name, price, quantity) VALUES 
    ('120P90', 'Google Home', 49.99, 10),
    ('43N23P', 'MacBook Pro', 5399.99, 5),
    ('A304SD', 'Alexa Speaker', 109.50, 10),
    ('234234', 'Raspberry Pi B', 30.00, 2)
    ;

INSERT INTO eshop.promotions (sku, active, promotions) VALUES
    ('43N23P', true, '[{"buy":1,"type":"FREE","item":"234234","units": 1}]'),
    ('120P90', true, '[{"buy":3,"type":"PRICE","units":2}]'),
    ('A304SD', true, '[{"buy":3,"type":"DISCOUNT","rate":10}]')