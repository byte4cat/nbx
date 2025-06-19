package pbconv

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTimeToPbTimestamp(t *testing.T) {
	t.Run("[SUCCESS] TimeToPbTimestamp in struct", func(t *testing.T) {
		type testPb struct {
			Fruit     string
			CreatedAt *timestamppb.Timestamp
			UpdatedAt *timestamppb.Timestamp
			DeletedAt *timestamppb.Timestamp
		}

		type testFrom struct {
			Fruit     string
			CreatedAt time.Time
			UpdatedAt *time.Time
			DeletedAt *time.Time
		}

		pbObj := testPb{
			Fruit: "Apple",
		}

		now := time.Now()
		fromObj := testFrom{
			Fruit:     "Apple",
			CreatedAt: now.AddDate(0, 0, -1),
			UpdatedAt: &now,
			DeletedAt: nil,
		}

		err := StructTimeToPbTimestamp(&pbObj, &fromObj, nil)

		assert.Nil(t, err)
		assert.Equal(t, fromObj.Fruit, pbObj.Fruit)
		assert.Equal(t, fromObj.CreatedAt.UTC().Format(time.RFC3339), pbObj.CreatedAt.AsTime().UTC().Format(time.RFC3339))
		assert.Equal(t, fromObj.UpdatedAt.UTC().Format(time.RFC3339), pbObj.UpdatedAt.AsTime().UTC().Format(time.RFC3339))

		assert.Nil(t, pbObj.DeletedAt)
		assert.Nil(t, fromObj.DeletedAt)
	})

	t.Run("[SUCCESS] TimeToPbTimestamp in slice of struct", func(t *testing.T) {
		type testPb struct {
			Fruit     string
			CreatedAt *timestamppb.Timestamp
			UpdatedAt *timestamppb.Timestamp
			DeletedAt *timestamppb.Timestamp
		}

		type testFrom struct {
			Fruit     string
			CreatedAt time.Time
			UpdatedAt *time.Time
			DeletedAt *time.Time
		}

		pbObjSlice := []testPb{
			{
				Fruit: "Apple",
			},
			{
				Fruit: "Banana",
			},
		}

		now := time.Now()
		fromObjSlice := []testFrom{
			{
				Fruit:     "Apple",
				CreatedAt: now.AddDate(0, 0, -1),
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			{
				Fruit:     "Banana",
				CreatedAt: now.AddDate(0, 0, -1),
				UpdatedAt: &now,
				DeletedAt: nil,
			},
		}

		err := SliceStructTimeToPbTimestamp(&pbObjSlice, &fromObjSlice, nil)

		assert.Nil(t, err)

		assert.Equal(t, fromObjSlice[0].Fruit, pbObjSlice[0].Fruit)
		assert.Equal(t, fromObjSlice[0].CreatedAt.UTC().Format(time.RFC3339), pbObjSlice[0].CreatedAt.AsTime().UTC().Format(time.RFC3339))
		assert.Equal(t, fromObjSlice[0].UpdatedAt.UTC().Format(time.RFC3339), pbObjSlice[0].UpdatedAt.AsTime().UTC().Format(time.RFC3339))

		assert.Equal(t, fromObjSlice[1].Fruit, pbObjSlice[1].Fruit)
		assert.Equal(t, fromObjSlice[1].CreatedAt.UTC().Format(time.RFC3339), pbObjSlice[1].CreatedAt.AsTime().UTC().Format(time.RFC3339))
		assert.Equal(t, fromObjSlice[1].UpdatedAt.UTC().Format(time.RFC3339), pbObjSlice[1].UpdatedAt.AsTime().UTC().Format(time.RFC3339))

		assert.Nil(t, pbObjSlice[0].DeletedAt)
		assert.Nil(t, fromObjSlice[0].DeletedAt)
		assert.Nil(t, pbObjSlice[1].DeletedAt)
		assert.Nil(t, fromObjSlice[1].DeletedAt)
	})

	t.Run("[SUCCESS] Slice in struct with time", func(t *testing.T) {
		type authorizer struct {
			AnsweredAt time.Time
			CreatedAt  time.Time
			UpdatedAt  time.Time
		}

		type unlockDeviceSessionView struct {
			Authorizers []authorizer
			AnsweredAt  time.Time
			CreatedAt   time.Time
			UpdatedAt   time.Time
		}

		type pbAuthorizer struct {
			AnsweredAt *timestamppb.Timestamp
			CreatedAt  *timestamppb.Timestamp
			UpdatedAt  *timestamppb.Timestamp
		}

		type pbUnlockDeviceSession struct {
			Authorizers []*pbAuthorizer
			AnsweredAt  *timestamppb.Timestamp
			CreatedAt   *timestamppb.Timestamp
			UpdatedAt   *timestamppb.Timestamp
		}

		// Prepare a sample UnlockDeviceSessionView with Authorizers and time fields
		now := time.Now().Truncate(time.Second)
		domain := unlockDeviceSessionView{
			AnsweredAt: now,
			CreatedAt:  now.Add(-2 * time.Hour),
			UpdatedAt:  now.Add(-1 * time.Hour),
			Authorizers: []authorizer{
				{
					AnsweredAt: now.Add(-10 * time.Minute),
					CreatedAt:  now.Add(-30 * time.Minute),
					UpdatedAt:  now.Add(-20 * time.Minute),
				},
				{
					AnsweredAt: now.Add(-40 * time.Minute),
					CreatedAt:  now.Add(-60 * time.Minute),
					UpdatedAt:  now.Add(-50 * time.Minute),
				},
			},
		}

		// Simulate marshaling and unmarshaling via JSON payload
		payloadBytes, err := json.Marshal(domain)
		require.NoError(t, err)

		var us unlockDeviceSessionView
		err = json.Unmarshal(payloadBytes, &us)
		require.NoError(t, err)

		// Copy to protobuf struct
		var pbMsg pbUnlockDeviceSession
		err = copier.Copy(&pbMsg, &us)
		require.NoError(t, err)

		// Call StructTimeToPbTimestamp for main struct and nested Authorizers
		fields := []string{"AnsweredAt", "CreatedAt", "UpdatedAt"}
		err = StructTimeToPbTimestamp(&pbMsg, &us, &fields, "Authorizers")
		require.NoError(t, err)

		// Main struct time fields
		require.Equal(t, now.Unix(), pbMsg.AnsweredAt.AsTime().Unix())
		require.Equal(t, now.Add(-2*time.Hour).Unix(), pbMsg.CreatedAt.AsTime().Unix())
		require.Equal(t, now.Add(-1*time.Hour).Unix(), pbMsg.UpdatedAt.AsTime().Unix())

		// Nested Authorizers
		require.Len(t, pbMsg.Authorizers, 2)
		require.Equal(t, now.Add(-10*time.Minute).Unix(), pbMsg.Authorizers[0].AnsweredAt.AsTime().Unix())
		require.Equal(t, now.Add(-30*time.Minute).Unix(), pbMsg.Authorizers[0].CreatedAt.AsTime().Unix())
		require.Equal(t, now.Add(-20*time.Minute).Unix(), pbMsg.Authorizers[0].UpdatedAt.AsTime().Unix())
		require.Equal(t, now.Add(-40*time.Minute).Unix(), pbMsg.Authorizers[1].AnsweredAt.AsTime().Unix())
		require.Equal(t, now.Add(-60*time.Minute).Unix(), pbMsg.Authorizers[1].CreatedAt.AsTime().Unix())
		require.Equal(t, now.Add(-50*time.Minute).Unix(), pbMsg.Authorizers[1].UpdatedAt.AsTime().Unix())
	})

	t.Run("[SUCCESS] Override default fields", func(t *testing.T) {
		type testPb struct {
			Fruit     string
			CreatedAt *timestamppb.Timestamp
			UpdatedAt *timestamppb.Timestamp
			DeletedAt *timestamppb.Timestamp
		}

		type testFrom struct {
			Fruit     string
			CreatedAt time.Time
			UpdatedAt *time.Time
			DeletedAt *time.Time
		}

		pbObj := testPb{
			Fruit: "Apple",
		}

		now := time.Now()
		fromObj := testFrom{
			Fruit:     "Apple",
			CreatedAt: now.AddDate(0, 0, -1),
			UpdatedAt: &now,
			DeletedAt: nil,
		}

		OverrideDefaultFields("CreatedAt", "UpdatedAt")

		err := StructTimeToPbTimestamp(&pbObj, &fromObj, nil)

		assert.Nil(t, err)
		assert.Equal(t, fromObj.Fruit, pbObj.Fruit)
		assert.Equal(t, fromObj.CreatedAt.UTC().Format(time.RFC3339), pbObj.CreatedAt.AsTime().UTC().Format(time.RFC3339))
		assert.Equal(t, fromObj.UpdatedAt.UTC().Format(time.RFC3339), pbObj.UpdatedAt.AsTime().UTC().Format(time.RFC3339))

		assert.Nil(t, pbObj.DeletedAt)
		assert.Nil(t, fromObj.DeletedAt)
	})

	t.Run("[SUCCESS] Convert sql.NullTime to *timestamppb.Timestamp", func(t *testing.T) {
		type testPb struct {
			Fruit     string
			CreatedAt *timestamppb.Timestamp
			UpdatedAt *timestamppb.Timestamp
			DeletedAt *timestamppb.Timestamp
		}

		type testFrom struct {
			Fruit     string
			CreatedAt sql.NullTime
			UpdatedAt sql.NullTime
			DeletedAt sql.NullTime
		}

		pbObj := testPb{
			Fruit: "Apple",
		}

		now := time.Now()
		fromObj := testFrom{
			Fruit:     "Apple",
			CreatedAt: sql.NullTime{Time: now.AddDate(0, 0, -1), Valid: true},
			UpdatedAt: sql.NullTime{Time: now, Valid: true},
			DeletedAt: sql.NullTime{Valid: false}, // Simulating a null value
		}

		err := StructTimeToPbTimestamp(&pbObj, &fromObj, nil)

		assert.Nil(t, err)
		assert.Equal(t, fromObj.Fruit, pbObj.Fruit)
		assert.False(t, fromObj.DeletedAt.Valid)
		assert.Equal(t, fromObj.CreatedAt.Time.UTC().Format(time.RFC3339), pbObj.CreatedAt.AsTime().UTC().Format(time.RFC3339))
		assert.Equal(t, fromObj.UpdatedAt.Time.UTC().Format(time.RFC3339), pbObj.UpdatedAt.AsTime().UTC().Format(time.RFC3339))

		assert.Nil(t, pbObj.DeletedAt)
		assert.Equal(t, pbObj.CreatedAt.AsTime().UTC().Format(time.RFC3339), fromObj.CreatedAt.Time.UTC().Format(time.RFC3339))
		assert.Equal(t, pbObj.UpdatedAt.AsTime().UTC().Format(time.RFC3339), fromObj.UpdatedAt.Time.UTC().Format(time.RFC3339))
	})

	t.Run("[SUCCESS] Convert *sql.NullTime to *timestamppb.Timestamp", func(t *testing.T) {
		type testPb struct {
			Fruit     string
			CreatedAt *timestamppb.Timestamp
			UpdatedAt *timestamppb.Timestamp
			DeletedAt *timestamppb.Timestamp
		}

		type testFrom struct {
			Fruit     string
			CreatedAt *sql.NullTime
			UpdatedAt *sql.NullTime
			DeletedAt *sql.NullTime
		}

		pbObj := testPb{
			Fruit: "Apple",
		}

		now := time.Now()
		fromObj := testFrom{
			Fruit:     "Apple",
			CreatedAt: &sql.NullTime{Time: now.AddDate(0, 0, -1), Valid: true},
			UpdatedAt: &sql.NullTime{Time: now, Valid: true},
			DeletedAt: &sql.NullTime{Valid: false}, // Simulating a null value
		}

		err := StructTimeToPbTimestamp(&pbObj, &fromObj, nil)

		assert.Nil(t, err)
		assert.Equal(t, fromObj.Fruit, pbObj.Fruit)
		assert.False(t, fromObj.DeletedAt.Valid)
		assert.Equal(t, fromObj.CreatedAt.Time.UTC().Format(time.RFC3339), pbObj.CreatedAt.AsTime().UTC().Format(time.RFC3339))
		assert.Equal(t, fromObj.UpdatedAt.Time.UTC().Format(time.RFC3339), pbObj.UpdatedAt.AsTime().UTC().Format(time.RFC3339))

		assert.Nil(t, pbObj.DeletedAt)
		assert.Equal(t, pbObj.CreatedAt.AsTime().UTC().Format(time.RFC3339), fromObj.CreatedAt.Time.UTC().Format(time.RFC3339))
		assert.Equal(t, pbObj.UpdatedAt.AsTime().UTC().Format(time.RFC3339), fromObj.UpdatedAt.Time.UTC().Format(time.RFC3339))
	})
}
