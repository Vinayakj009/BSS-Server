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