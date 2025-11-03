package database

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

func createDbForPlanTests(t *testing.T) (context.Context, *DB) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	ctx := context.Background()

	// Create database connection
	db, err := NewDb(ctx)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	return ctx, db
}

// TestPostgresIntegration tests the PostgreSQL connection and basic operations
func TestPostgresIntegration(t *testing.T) {
	ctx, db := createDbForPlanTests(t)

	// Test ping
	if err := db.Ping(ctx); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	t.Log("Successfully connected to PostgreSQL")
}

func TestInitSchema(t *testing.T) {
	ctx, db := createDbForPlanTests(t)

	t.Log("Successfully initialized database schema")

	t.Run("InsertSubscription", func(t *testing.T) {
		// First create a plan
		var planID uuid.UUID
		err := db.Pool.QueryRow(ctx, `
			INSERT INTO plans (code, name, price_cents, currency, duration_days, data_mb)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, "TEST-SUB-PLAN", "Test Subscription Plan", 1999, "USD", 30, 2048).Scan(&planID)
		if err != nil {
			t.Fatalf("Failed to insert plan: %v", err)
		}

		// Insert subscription
		customerID := uuid.New()
		startDate := time.Now()
		endDate := startDate.Add(30 * 24 * time.Hour)

		var subID uuid.UUID
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO subscriptions (customer_id, plan_id, start_date, end_date, status, auto_renew)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, customerID, planID, startDate, endDate, "ACTIVE", true).Scan(&subID)
		if err != nil {
			t.Fatalf("Failed to insert subscription: %v", err)
		}
		t.Logf("Successfully inserted subscription with ID: %s", subID)

		// Clean up
		_, err = db.Pool.Exec(ctx, "DELETE FROM subscriptions WHERE id = $1", subID)
		if err != nil {
			t.Errorf("Failed to clean up test subscription: %v", err)
		}
		_, err = db.Pool.Exec(ctx, "DELETE FROM plans WHERE id = $1", planID)
		if err != nil {
			t.Errorf("Failed to clean up test plan: %v", err)
		}
	})

	t.Run("InsertEvent", func(t *testing.T) {
		resourceID := uuid.New()
		payload := []byte(`{"action": "test", "data": "sample"}`)

		var eventID int64
		err := db.Pool.QueryRow(ctx, `
			INSERT INTO events (event_type, resource_id, payload)
			VALUES ($1, $2, $3)
			RETURNING id
		`, "subscription.created", resourceID, payload).Scan(&eventID)
		if err != nil {
			t.Fatalf("Failed to insert event: %v", err)
		}
		t.Logf("Successfully inserted event with ID: %d", eventID)

		// Clean up
		_, err = db.Pool.Exec(ctx, "DELETE FROM events WHERE id = $1", eventID)
		if err != nil {
			t.Errorf("Failed to clean up test event: %v", err)
		}
	})
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Set test environment variables
	testCases := []struct {
		name     string
		envVars  map[string]string
		expected *Config
	}{
		{
			name: "WithAllEnvVars",
			envVars: map[string]string{
				"POSTGRES_HOST":     "testhost",
				"POSTGRES_PORT":     "5433",
				"POSTGRES_USER":     "testuser",
				"POSTGRES_PASSWORD": "testpass",
				"POSTGRES_DB":       "testdb",
				"POSTGRES_SSLMODE":  "require",
			},
			expected: &Config{
				Host:     "testhost",
				Port:     "5433",
				User:     "testuser",
				Password: "testpass",
				Database: "testdb",
				SSLMode:  "require",
			},
		},
		{
			name:    "WithDefaults",
			envVars: map[string]string{},
			expected: &Config{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "postgres",
				Database: "bss",
				SSLMode:  "disable",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear all postgres env vars
			os.Unsetenv("POSTGRES_HOST")
			os.Unsetenv("POSTGRES_PORT")
			os.Unsetenv("POSTGRES_USER")
			os.Unsetenv("POSTGRES_PASSWORD")
			os.Unsetenv("POSTGRES_DB")
			os.Unsetenv("POSTGRES_SSLMODE")

			// Set test env vars
			for k, v := range tc.envVars {
				os.Setenv(k, v)
			}

			config := loadConfigFromEnv()

			if config.Host != tc.expected.Host {
				t.Errorf("Expected Host=%s, got %s", tc.expected.Host, config.Host)
			}
			if config.Port != tc.expected.Port {
				t.Errorf("Expected Port=%s, got %s", tc.expected.Port, config.Port)
			}
			if config.User != tc.expected.User {
				t.Errorf("Expected User=%s, got %s", tc.expected.User, config.User)
			}
			if config.Password != tc.expected.Password {
				t.Errorf("Expected Password=%s, got %s", tc.expected.Password, config.Password)
			}
			if config.Database != tc.expected.Database {
				t.Errorf("Expected Database=%s, got %s", tc.expected.Database, config.Database)
			}
			if config.SSLMode != tc.expected.SSLMode {
				t.Errorf("Expected SSLMode=%s, got %s", tc.expected.SSLMode, config.SSLMode)
			}
		})
	}
}
