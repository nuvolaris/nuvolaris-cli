package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNuvScan(t *testing.T) {
	t.Run("should have scan subcmd help", func(t *testing.T) {
		var cli CLI

		app := NewTestApp(t, &cli)
		require.PanicsWithValue(t, true, func() {
			_, err := app.Parse([]string{"scan", "--help"})
			require.NoError(t, err)
		})
	})
}
