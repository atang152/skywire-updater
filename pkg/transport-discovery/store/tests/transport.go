package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/watercompany/skywire-node/pkg/transport"
	"github.com/watercompany/skywire-services/pkg/transport-discovery/store"
)

type TransportSuite struct {
	suite.Suite
	Store store.TransportStore
}

func (s *TransportSuite) SetupTest() {
}

func (s *TransportSuite) TestRegister() {
	t := s.T()
	ctx := context.Background()

	pk1, _ := cipher.GenerateKeyPair()
	pk2, _ := cipher.GenerateKeyPair()

	sEntry := &transport.SignedEntry{
		Entry: &transport.Entry{
			ID:     uuid.New(),
			Edges:  [2]string{pk1.Hex(), pk2.Hex()},
			Type:   "messaging",
			Public: true,
		},
		Signatures: [2]string{"foo", "bar"},
	}

	t.Run(".RegisterTransport", func(t *testing.T) {
		require.NoError(t, s.Store.RegisterTransport(ctx, sEntry))
		assert.True(t, sEntry.Registered > 0)
	})

	t.Run(".GetTransportByID", func(t *testing.T) {
		found, err := s.Store.GetTransportByID(ctx, sEntry.Entry.ID)
		require.NoError(t, err)
		assert.Equal(t, sEntry.Entry, found.Entry)
		assert.True(t, found.IsUp)
	})

	t.Run(".GetTransportsByEdge", func(t *testing.T) {
		entries, err := s.Store.GetTransportsByEdge(ctx, pk1)
		require.NoError(t, err)
		require.Len(t, entries, 1)
		assert.Equal(t, sEntry.Entry, entries[0].Entry)
		assert.True(t, entries[0].IsUp)

		entries, err = s.Store.GetTransportsByEdge(ctx, pk2)
		require.NoError(t, err)
		require.Len(t, entries, 1)
		assert.Equal(t, sEntry.Entry, entries[0].Entry)
		assert.True(t, entries[0].IsUp)

		pk, _ := cipher.GenerateKeyPair()
		entries, err = s.Store.GetTransportsByEdge(ctx, pk)
		require.NoError(t, err)
		require.Len(t, entries, 0)
	})

	t.Run(".UpdateStatus", func(t *testing.T) {
		entry, err := s.Store.UpdateStatus(ctx, sEntry.Entry.ID, false)
		require.Error(t, err)
		assert.Equal(t, "invalid auth", err.Error())

		entry, err = s.Store.UpdateStatus(context.WithValue(ctx, "auth-pub-key", pk1), sEntry.Entry.ID, false)
		require.NoError(t, err)
		assert.Equal(t, sEntry.Entry, entry.Entry)
		assert.False(t, entry.IsUp)
	})

	t.Run(".DeregisterTransport", func(t *testing.T) {
		entry, err := s.Store.DeregisterTransport(ctx, sEntry.Entry.ID)
		require.NoError(t, err)
		assert.Equal(t, sEntry.Entry, entry)

		_, err = s.Store.GetTransportByID(ctx, sEntry.Entry.ID)
		require.Error(t, err)
		assert.Equal(t, "Transport not found", err.Error())
	})
}
