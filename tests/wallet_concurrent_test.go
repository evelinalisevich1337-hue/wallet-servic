func TestConcurrentWithdraw(t *testing.T) {
	ctx := context.Background()
	walletID := uuid.New()


	_, err := db.Exec(ctx,
		`INSERT INTO wallets (id, balance) VALUES ($1, $2)`,
		walletID, 100000,
	)
	require.NoError(t, err)

	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = repo.UpdateBalance(ctx, walletID, 1, "WITHDRAW")
		}()
	}

	wg.Wait()

	balance, _ := repo.GetBalance(ctx, walletID)
	require.Equal(t, int64(99000), balance)
}

