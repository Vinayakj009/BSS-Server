\c bss;
-- Plans table
CREATE TABLE IF NOT EXISTS plans (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	code VARCHAR(50) UNIQUE NOT NULL,
	name VARCHAR(255) NOT NULL,
	price_cents BIGINT NOT NULL,
	currency VARCHAR(3) NOT NULL DEFAULT 'USD',
	duration_days INTEGER NOT NULL,
	data_mb BIGINT NOT NULL,
	active BOOLEAN NOT NULL DEFAULT true,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Insert sample plan
INSERT INTO plans (id, code, name, price_cents, currency, duration_days, data_mb, active)
VALUES (
    '11111111-1111-1111-1111-111111111111',
    'BASIC-MONTHLY',
    'Basic Monthly Plan',
    999,
    'USD',
    30,
    5120,
    true
);

-- Subscriptions table
CREATE TABLE IF NOT EXISTS subscriptions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	customer_id UUID NOT NULL,
	plan_id UUID NOT NULL REFERENCES plans(id),
	start_date TIMESTAMP WITH TIME ZONE NOT NULL,
	end_date TIMESTAMP WITH TIME ZONE NOT NULL,
	status VARCHAR(20) NOT NULL CHECK (status IN ('ACTIVE', 'CANCELLED', 'EXPIRED')),
	auto_renew BOOLEAN NOT NULL DEFAULT true,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Insert sample subscription
INSERT INTO subscriptions (id, customer_id, plan_id, start_date, end_date, status, auto_renew)
VALUES (
    gen_random_uuid(),
    '00000000-0000-0000-0000-000000000000',
    '11111111-1111-1111-1111-111111111111',
    NOW() - INTERVAL '1 days',
    NOW() + INTERVAL '30 days',
    'ACTIVE',
    true
),
(
    gen_random_uuid(),
    '00000000-0000-0000-0000-000000000000',
    '11111111-1111-1111-1111-111111111111',
    NOW() - INTERVAL '30 days',
    NOW() - INTERVAL '2 days',
    'EXPIRED',
    true
),
(
    gen_random_uuid(),
    '00000000-0000-0000-0000-000000000001',
    '11111111-1111-1111-1111-111111111111',
    NOW() - INTERVAL '32 days',
    NOW() - INTERVAL '4 days',
    'EXPIRED',
    true
);

-- Events table
CREATE TABLE IF NOT EXISTS events (
	id BIGSERIAL PRIMARY KEY,
	event_type VARCHAR(100) NOT NULL,
	resource_id UUID NOT NULL,
	payload JSONB,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_subscriptions_customer_id ON subscriptions(customer_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_plan_id ON subscriptions(plan_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_events_resource_id ON events(resource_id);
CREATE INDEX IF NOT EXISTS idx_events_event_type ON events(event_type);
CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at);