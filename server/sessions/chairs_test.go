package sessions_test

import (
	"encoding/json"
	"testing"

	"github.com/DaanV2/f1-game-dashboards/server/sessions"
	"github.com/stretchr/testify/require"
)

func Test_Chairs_Json(t *testing.T) {
	// Test if the chair can be marshalled to json
	chair := sessions.NewChair("test", 1234, true)

	data, err := json.Marshal(chair)
	require.NoError(t, err)

	// Test if the chair can be unmarshalled from json
	var chair2 sessions.Chair
	err = json.Unmarshal(data, &chair2)
	require.NoError(t, err)

	require.Equal(t, chair, chair2)
}

