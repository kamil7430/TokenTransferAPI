package service

import (
	"context"
	"sync"
	"testing"

	"github.com/kamil7430/TokenTransferAPI/graph/model"
	"github.com/kamil7430/TokenTransferAPI/repository"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestWalletService(t *testing.T) {
	ctx := context.Background()
	dbname := "serviceTests"
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

	d := WalletService{
		WalletRepository: &repository.DatabaseWalletRepository{},
		Database:         db,
	}

	t.Run("get wallet", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000001", 100)

		wallet, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000001")
		require.NoError(t, err)
		require.Equal(t, 100, wallet.Tokens)
		require.Equal(t, "0x0000000000000000000000000000000000000001", wallet.Address)
	})

	t.Run("transfer", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000001", 100)
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000002", 200)

		amount, err := d.Transfer(ctx, "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002", 60)
		require.NoError(t, err)
		require.Equal(t, 40, amount)

		fromWallet, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000001")
		require.NoError(t, err)
		require.Equal(t, "0x0000000000000000000000000000000000000001", fromWallet.Address)
		require.Equal(t, 40, fromWallet.Tokens)

		toWallet, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000002")
		require.NoError(t, err)
		require.Equal(t, "0x0000000000000000000000000000000000000002", toWallet.Address)
		require.Equal(t, 260, toWallet.Tokens)
	})

	t.Run("transfer negative token amount", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000001", 100)
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000002", 200)

		_, err := d.Transfer(ctx, "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002", -60)
		require.Error(t, err)
	})

	t.Run("transfer amount higher than wallet balance", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000001", 100)
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000002", 0)

		_, err := d.Transfer(ctx, "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002", 260)
		require.Error(t, err)
	})

	t.Run("transfer from non-existing wallet", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000002", 100)

		_, err := d.Transfer(ctx, "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002", 60)
		require.Error(t, err)
	})

	t.Run("transfer to non-existing wallet", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000001", 100)

		amount, err := d.Transfer(ctx, "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002", 60)
		require.NoError(t, err)
		require.Equal(t, 40, amount)

		fromWallet, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000001")
		require.NoError(t, err)
		require.Equal(t, "0x0000000000000000000000000000000000000001", fromWallet.Address)
		require.Equal(t, 40, fromWallet.Tokens)

		toWallet, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000002")
		require.NoError(t, err)
		require.Equal(t, "0x0000000000000000000000000000000000000002", toWallet.Address)
		require.Equal(t, 60, toWallet.Tokens)
	})

	t.Run("transfer to own wallet", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000001", 100)

		_, err := d.Transfer(ctx, "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000001", 60)
		require.Error(t, err)
	})

	t.Run("parallel transfers example from task", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000001", 10)
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000002", 10)

		const concurrentRoutines = 3
		barrier := make(chan struct{})

		var workWG sync.WaitGroup
		workWG.Add(concurrentRoutines)

		var barrierWG sync.WaitGroup
		barrierWG.Add(concurrentRoutines)

		// 7 tokens from 1 to 2
		go func() {
			barrierWG.Done() // report readiness to start
			<-barrier        // wait on barrier
			_, _ = d.Transfer(ctx, "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002", 7)
			// no error checking because it can either succeed or fail
			workWG.Done()
		}()

		// 4 tokens from 1 to 2
		go func() {
			barrierWG.Done()
			<-barrier
			_, _ = d.Transfer(ctx, "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002", 4)
			workWG.Done()
		}()

		// 1 token from 2 to 1
		go func() {
			barrierWG.Done()
			<-barrier
			_, _ = d.Transfer(ctx, "0x0000000000000000000000000000000000000002", "0x0000000000000000000000000000000000000001", 1)
			workWG.Done()
		}()

		barrierWG.Wait() // wait for all go routines to get ready
		close(barrier)   // unblock all the go routines waiting on the barrier
		workWG.Wait()    // wait for all go routines to finish

		wallet1, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000001")
		require.NoError(t, err)
		wallet2, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000002")
		require.NoError(t, err)

		require.Condition(t, func() bool {
			return (wallet1.Tokens == 7 && wallet2.Tokens == 13) ||
				(wallet1.Tokens == 4 && wallet2.Tokens == 16) ||
				(wallet1.Tokens == 0 && wallet2.Tokens == 20)
		})
	})

	t.Run("cross transfer", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000001", 15)
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000002", 10)

		const concurrentRoutines = 2
		barrier := make(chan struct{})

		var workWG sync.WaitGroup
		workWG.Add(concurrentRoutines)

		var barrierWG sync.WaitGroup
		barrierWG.Add(concurrentRoutines)

		// 15 tokens from 1 to 2
		go func() {
			barrierWG.Done()
			<-barrier
			_, _ = d.Transfer(ctx, "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002", 10)
			workWG.Done()
		}()

		// 15 tokens from 2 to 1
		go func() {
			barrierWG.Done()
			<-barrier
			_, _ = d.Transfer(ctx, "0x0000000000000000000000000000000000000002", "0x0000000000000000000000000000000000000001", 10)
			workWG.Done()
		}()

		barrierWG.Wait()
		close(barrier)
		workWG.Wait()

		wallet1, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000001")
		require.NoError(t, err)
		wallet2, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000002")
		require.NoError(t, err)

		require.Equal(t, 15, wallet1.Tokens)
		require.Equal(t, 10, wallet2.Tokens)
	})

	t.Run("parallel transfers to non-existing wallet", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000001", 15)

		const concurrentRoutines = 2
		barrier := make(chan struct{})

		var workWG sync.WaitGroup
		workWG.Add(concurrentRoutines)

		var barrierWG sync.WaitGroup
		barrierWG.Add(concurrentRoutines)

		// two times: 5 tokens from 1 to 2
		for i := 0; i < concurrentRoutines; i++ {
			go func() {
				barrierWG.Done()
				<-barrier
				_, _ = d.Transfer(ctx, "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002", 5)
				workWG.Done()
			}()
		}

		barrierWG.Wait()
		close(barrier)
		workWG.Wait()

		wallet1, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000001")
		require.NoError(t, err)
		wallet2, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000002")
		require.NoError(t, err)

		require.Equal(t, 5, wallet1.Tokens)
		require.Equal(t, 10, wallet2.Tokens)
	})

	t.Run("massive parallel transfers between three wallets", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE Wallets")
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000001", 1000)
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000002", 2000)
		db.Exec("INSERT INTO Wallets(Address, Tokens) VALUES ($1, $2)", "0x0000000000000000000000000000000000000003", 500)

		const concurrentRoutines = 30
		barrier := make(chan struct{})

		var workWG sync.WaitGroup
		workWG.Add(concurrentRoutines)

		var barrierWG sync.WaitGroup
		barrierWG.Add(concurrentRoutines)

		// 10 times: 10 tokens from 1 to 2
		for i := 0; i < concurrentRoutines/3; i++ {
			go func() {
				barrierWG.Done()
				<-barrier
				_, _ = d.Transfer(ctx, "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002", 10)
				workWG.Done()
			}()
		}

		// 10 times: 10 tokens from 2 to 3
		for i := 0; i < concurrentRoutines/3; i++ {
			go func() {
				barrierWG.Done()
				<-barrier
				_, _ = d.Transfer(ctx, "0x0000000000000000000000000000000000000002", "0x0000000000000000000000000000000000000003", 10)
				workWG.Done()
			}()
		}

		// 10 times: 5 tokens from 3 to 1
		for i := 0; i < concurrentRoutines/3; i++ {
			go func() {
				barrierWG.Done()
				<-barrier
				_, _ = d.Transfer(ctx, "0x0000000000000000000000000000000000000003", "0x0000000000000000000000000000000000000001", 5)
				workWG.Done()
			}()
		}

		barrierWG.Wait()
		close(barrier)
		workWG.Wait()

		wallet1, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000001")
		require.NoError(t, err)
		wallet2, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000002")
		require.NoError(t, err)
		wallet3, err := d.GetWallet(ctx, "0x0000000000000000000000000000000000000003")
		require.NoError(t, err)

		require.Equal(t, 950, wallet1.Tokens)
		require.Equal(t, 2000, wallet2.Tokens)
		require.Equal(t, 550, wallet3.Tokens)
	})
}
