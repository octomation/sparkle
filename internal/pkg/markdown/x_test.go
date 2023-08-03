package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.octolab.org/pointer"
)

func TestProperties_Decode(t *testing.T) {
	props := map[string]any{
		"uid":       "10000000-2000-3000-4000-500000000000",
		"aliases":   []string{"alias1", "alias2"},
		"tags":      []string{"tag1", "tag2"},
		"date":      "2020-01-01",
		"sleep":     "8h",
		"health":    "good",
		"mood":      "happy",
		"tracking":  "https://example.com/",
		"duration":  "1h",
		"timestamp": 1234567890,
	}

	var p DynamicProperties
	var (
		diary = new(DiaryEntry)
		mood  = new(MoodTracker)
		track = new(TimeTracker)
	)
	for _, component := range []Component{diary, mood, track} {
		if component.IsPresent(props) {
			p.Components = append(p.Components, component)
		}
	}
	assert.NoError(t, p.Decode(props))
	assert.Equal(t, props["uid"], p.UID)
	assert.Equal(t, props["aliases"], p.Aliases)
	assert.Equal(t, props["tags"], p.Tags)
	assert.Equal(t, props["date"], diary.Date)
	assert.Equal(t, props["sleep"], pointer.ValueOfString(mood.Sleep))
	assert.Equal(t, props["health"], pointer.ValueOfString(mood.Health))
	assert.Equal(t, props["mood"], pointer.ValueOfString(mood.Mood))
	assert.Equal(t, props["tracking"], pointer.ValueOfString(track.Tracking))
	assert.Equal(t, props["duration"], pointer.ValueOfString(track.Duration))
	assert.Equal(t, props["timestamp"], p.Rest["timestamp"])
}
