package repository

import (
	"context"
	"testing"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDatabaseWalletRepository(t *testing.T) {
	ctx := context.Background()
	dbname := "repositoryTests"
	dbuser := "user"
	dbpassword := "password"

	ctr, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbname),
		postgres.WithUsername(dbuser),
		postgres.WithPassword(dbpassword),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	err = ctr.Snapshot(ctx)
	require.NoError(t, err)

	dbURL, err := ctr.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := gorm.Open(gormpostgres.Open(dbURL), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Wallet{})
	require.NoError(t, err)

	d := DatabaseWalletRepository{}

	t.Run("create wallet", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")

		wallet := &model.Wallet{
			Address: "0x0000000000000000000000000000000000000000",
			Tokens:  1_000_000,
		}

		err := d.AddWallet(ctx, db, wallet)
		require.NoError(t, err)
	})

	t.Run("query existing wallet", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000000", 1_000_000)

		wallet, err := d.GetWalletByAddress(ctx, db, "0x0000000000000000000000000000000000000000")
		require.NoError(t, err)
		require.Equal(t, "0x0000000000000000000000000000000000000000", wallet.Address)
		require.Equal(t, 1_000_000, wallet.Tokens)
	})

	t.Run("query non-existing wallet", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")

		_, err := d.GetWalletByAddress(ctx, db, "0x000000000000000000000000000000000000")
		require.Error(t, err)
	})

	t.Run("update wallet", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000000", 1_000_000)

		err := d.UpdateWalletTokensByAddress(ctx, db, "0x0000000000000000000000000000000000000000", 150)
		require.NoError(t, err)

		wallet, err := d.GetWalletByAddress(ctx, db, "0x0000000000000000000000000000000000000000")
		require.NoError(t, err)
		require.Equal(t, "0x0000000000000000000000000000000000000000", wallet.Address)
		require.Equal(t, 150, wallet.Tokens)
	})
}
