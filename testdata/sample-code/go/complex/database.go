package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// ============================================================================
// Configuration
// ============================================================================

// Config holds database configuration
type Config struct {
	Driver          string        `json:"driver" env:"DB_DRIVER" default:"postgres"`
	Host            string        `json:"host" env:"DB_HOST" default:"localhost"`
	Port            int           `json:"port" env:"DB_PORT" default:"5432"`
	User            string        `json:"user" env:"DB_USER" default:"postgres"`
	Password        string        `json:"password" env:"DB_PASSWORD" default:"postgres"`
	Database        string        `json:"database" env:"DB_NAME" default:"myapp"`
	SSLMode         string        `json:"ssl_mode" env:"DB_SSL_MODE" default:"disable"`
	MaxOpenConns    int           `json:"max_open_conns" env:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int           `json:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME" default:"5m"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time" env:"DB_CONN_MAX_IDLE_TIME" default:"5m"`
}

// ConnectionString returns the database connection string based on driver
func (c *Config) ConnectionString() string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
		)
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
			c.User, c.Password, c.Host, c.Port, c.Database,
		)
	case "sqlite3":
		return c.Database
	default:
		return ""
	}
}

// ============================================================================
// Database Interface
// ============================================================================

// DB is the main database interface
type DB interface {
	// Core operations
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row

	// Transaction management
	Begin(ctx context.Context) (Tx, error)
	WithTx(ctx context.Context, fn func(Tx) error) error

	// Health checks
	Ping(ctx context.Context) error
	Stats() sql.DBStats

	// Close
	Close() error
}

// Tx represents a database transaction
type Tx interface {
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	Commit() error
	Rollback() error
}

// ============================================================================
// SQL Implementation
// ============================================================================

type sqlDB struct {
	db  *sql.DB
	cfg *Config
	mu  sync.RWMutex
}

type sqlTx struct {
	tx *sql.Tx
}

// New creates a new database connection
func New(cfg *Config) (DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &sqlDB{
		db:  db,
		cfg: cfg,
	}, nil
}

func (d *sqlDB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.db.ExecContext(ctx, query, args...)
}

func (d *sqlDB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.db.QueryContext(ctx, query, args...)
}

func (d *sqlDB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *sqlDB) Begin(ctx context.Context) (Tx, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &sqlTx{tx: tx}, nil
}

func (d *sqlDB) WithTx(ctx context.Context, fn func(Tx) error) error {
	tx, err := d.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback error: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	return tx.Commit()
}

func (d *sqlDB) Ping(ctx context.Context) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return d.db.PingContext(ctx)
}

func (d *sqlDB) Stats() sql.DBStats {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.db.Stats()
}

func (d *sqlDB) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.db.Close()
}

func (t *sqlTx) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *sqlTx) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *sqlTx) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *sqlTx) Commit() error {
	return t.tx.Commit()
}

func (t *sqlTx) Rollback() error {
	return t.tx.Rollback()
}

// ============================================================================
// Repository Pattern
// ============================================================================

// UserRepository handles user database operations
type UserRepository struct {
	db DB
}

func NewUserRepository(db DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (email, password, name, role, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		user.Email, user.Password, user.Name, user.Role, user.Active,
		now, now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrDuplicateEmail
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	query := `
		SELECT id, email, password, name, role, active, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	user := &User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name,
		&user.Role, &user.Active, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password, name, role, active, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	user := &User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name,
		&user.Role, &user.Active, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, role = $3, active = $4, updated_at = $5
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		user.Name, user.Email, user.Role, user.Active, time.Now(), user.ID,
	).Scan(&user.UpdatedAt)

	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE users SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *UserRepository) List(ctx context.Context, filter UserFilter) ([]*User, error) {
	query := `
		SELECT id, email, name, role, active, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}
	argIdx := 1

	if filter.Role != "" {
		query += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, filter.Role)
		argIdx++
	}

	if filter.Active != nil {
		query += fmt.Sprintf(" AND active = $%d", argIdx)
		args = append(args, *filter.Active)
		argIdx++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (email ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}

	query += " ORDER BY id DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filter.Offset)
		argIdx++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.Name,
			&user.Role, &user.Active, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// UserFilter represents filter options for user listing
type UserFilter struct {
	Role   string
	Active *bool
	Search string
	Limit  int
	Offset int
}

// ============================================================================
// Migrations
// ============================================================================

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrator handles database migrations
type Migrator struct {
	db     *sql.DB
	driver string
}

func NewMigrator(db *sql.DB, driver string) *Migrator {
	return &Migrator{
		db:     db,
		driver: driver,
	}
}

func (m *Migrator) Up() error {
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create source: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", source, m.driver, driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (m *Migrator) Down() error {
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create source: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", source, m.driver, driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if err := migrator.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	return nil
}

func (m *Migrator) Version() (uint, bool, error) {
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create driver: %w", err)
	}

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return 0, false, fmt.Errorf("failed to create source: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", source, m.driver, driver)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrator: %w", err)
	}

	return migrator.Version()
}

// ============================================================================
// Connection Pool Management
// ============================================================================

// PoolManager manages multiple database connections
type PoolManager struct {
	pools map[string]DB
	mu    sync.RWMutex
}

func NewPoolManager() *PoolManager {
	return &PoolManager{
		pools: make(map[string]DB),
	}
}

func (pm *PoolManager) Get(name string) DB {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.pools[name]
}

func (pm *PoolManager) Add(name string, db DB) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.pools[name] = db
}

func (pm *PoolManager) Remove(name string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.pools, name)
}

func (pm *PoolManager) CloseAll() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var errs []error
	for name, db := range pm.pools {
		if err := db.Close(); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing pools: %v", errs)
	}

	return nil
}

// ============================================================================
// Health Checker
// ============================================================================

// HealthChecker performs database health checks
type HealthChecker struct {
	db DB
}

func NewHealthChecker(db DB) *HealthChecker {
	return &HealthChecker{db: db}
}

func (h *HealthChecker) Check(ctx context.Context) map[string]interface{} {
	result := make(map[string]interface{})

	// Ping database
	start := time.Now()
	err := h.db.Ping(ctx)
	pingTime := time.Since(start)

	if err != nil {
		result["status"] = "down"
		result["error"] = err.Error()
	} else {
		result["status"] = "up"
		result["ping_ms"] = pingTime.Milliseconds()
	}

	// Get stats
	stats := h.db.Stats()
	result["stats"] = map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration_ms":     stats.WaitDuration.Milliseconds(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}

	return result
}

// ============================================================================
// Backup and Restore
// ============================================================================

// BackupManager handles database backups
type BackupManager struct {
	db     DB
	config *Config
}

func NewBackupManager(db DB, config *Config) *BackupManager {
	return &BackupManager{
		db:     db,
		config: config,
	}
}

func (b *BackupManager) Backup(ctx context.Context, path string) error {
	// This would use pg_dump or similar
	// Implementation depends on database driver
	return nil
}

func (b *BackupManager) Restore(ctx context.Context, path string) error {
	// This would use pg_restore or similar
	return nil
}

// ============================================================================
// Query Builder
// ============================================================================

// QueryBuilder helps build SQL queries safely
type QueryBuilder struct {
	query strings.Builder
	args  []interface{}
	pos   int
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		args: make([]interface{}, 0),
	}
}

func (qb *QueryBuilder) Write(s string) *QueryBuilder {
	qb.query.WriteString(s)
	return qb
}

func (qb *QueryBuilder) Writef(format string, args ...interface{}) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(format, args...))
	return qb
}

func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	if qb.pos == 0 {
		qb.query.WriteString(" WHERE ")
	} else {
		qb.query.WriteString(" AND ")
	}

	qb.query.WriteString(condition)
	qb.args = append(qb.args, args...)
	qb.pos++

	return qb
}

func (qb *QueryBuilder) In(column string, values []interface{}) *QueryBuilder {
	if qb.pos == 0 {
		qb.query.WriteString(" WHERE ")
	} else {
		qb.query.WriteString(" AND ")
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("$%d", len(qb.args)+i+1)
	}

	qb.query.WriteString(fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", ")))
	qb.args = append(qb.args, values...)
	qb.pos++

	return qb
}

func (qb *QueryBuilder) OrderBy(column string, desc bool) *QueryBuilder {
	qb.query.WriteString(" ORDER BY " + column)
	if desc {
		qb.query.WriteString(" DESC")
	}
	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(" LIMIT $%d", len(qb.args)+1))
	qb.args = append(qb.args, limit)
	return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(" OFFSET $%d", len(qb.args)+1))
	qb.args = append(qb.args, offset)
	return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
	return qb.query.String(), qb.args
}

// ============================================================================
// Custom Errors
// ============================================================================

var (
	ErrNotFound            = errors.New("record not found")
	ErrDuplicateEmail      = errors.New("email already exists")
	ErrDuplicateKey        = errors.New("duplicate key violation")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrCheckViolation      = errors.New("check constraint violation")
)

func isDuplicateKeyError(err error) bool {
	// PostgreSQL error code for unique violation is "23505"
	if pqErr, ok := err.(interface{ Code() string }); ok {
		return pqErr.Code() == "23505"
	}
	return strings.Contains(err.Error(), "duplicate key")
}

// ============================================================================
// Models
// ============================================================================

type User struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Name      string     `json:"name"`
	Role      string     `json:"role"`
	Active    bool       `json:"active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// ============================================================================
// Main Example
// ============================================================================

func main() {
	// Load configuration from environment
	cfg := &Config{
		Driver:          getEnv("DB_DRIVER", "postgres"),
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnvInt("DB_PORT", 5432),
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", "postgres"),
		Database:        getEnv("DB_NAME", "myapp"),
		SSLMode:         getEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
	}

	// Connect to database
	db, err := New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	// Run migrations
	migrator := NewMigrator(db.(*sqlDB).db, cfg.Driver)
	if err := migrator.Up(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create repository
	userRepo := NewUserRepository(db)

	// Example: Create user
	ctx := context.Background()
	user := &User{
		Email:    "john@example.com",
		Password: "hashed_password",
		Name:     "John Doe",
		Role:     "user",
		Active:   true,
	}

	if err := userRepo.Create(ctx, user); err != nil {
		log.Printf("Failed to create user: %v", err)
	}

	// Example: Get user
	user, err = userRepo.GetByEmail(ctx, "john@example.com")
	if err != nil {
		log.Printf("Failed to get user: %v", err)
	}

	// Health check
	health := NewHealthChecker(db)
	status := health.Check(ctx)
	log.Printf("Database status: %v", status)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intVal int
		fmt.Sscanf(value, "%d", &intVal)
		return intVal
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		dur, err := time.ParseDuration(value)
		if err == nil {
			return dur
		}
	}
	return defaultValue
}
