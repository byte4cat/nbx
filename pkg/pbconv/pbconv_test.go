package pbconv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func TestTimeToPbTimestamp(t *testing.T) {
	t.Run("[SUCCESS] TimeToPbTimestamp in struct", func(t *testing.T) {
		pbObj := testPb{
			Fruit: "Apple",
		}

		ct := time.Now()
		fromObj := testFrom{
			Fruit:     "Apple",
			CreatedAt: ct.AddDate(0, 0, -1),
			UpdatedAt: &ct,
			DeletedAt: nil,
		}

		err := StructTimeToPbTimestamp(&pbObj, &fromObj)

		assert.Nil(t, err)
		assert.Equal(t, fromObj.Fruit, pbObj.Fruit)
		assert.Equal(t, fromObj.CreatedAt.UTC().Format(time.RFC3339), pbObj.CreatedAt.AsTime().UTC().Format(time.RFC3339))
		assert.Equal(t, fromObj.UpdatedAt.UTC().Format(time.RFC3339), pbObj.UpdatedAt.AsTime().UTC().Format(time.RFC3339))

		assert.Nil(t, pbObj.DeletedAt)
		assert.Nil(t, fromObj.DeletedAt)
	})

	t.Run("[SUCCESS] TimeToPbTimestamp in slice of struct", func(t *testing.T) {
		pbObjSlice := []testPb{
			{
				Fruit: "Apple",
			},
			{
				Fruit: "Banana",
			},
		}

		ct := time.Now()
		fromObjSlice := []testFrom{
			{
				Fruit:     "Apple",
				CreatedAt: ct.AddDate(0, 0, -1),
				UpdatedAt: &ct,
				DeletedAt: nil,
			},
			{
				Fruit:     "Banana",
				CreatedAt: ct.AddDate(0, 0, -1),
				UpdatedAt: &ct,
				DeletedAt: nil,
			},
		}

		err := SliceStructTimeToPbTimestamp(&pbObjSlice, &fromObjSlice)

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
}
