package markdown

import "github.com/mitchellh/mapstructure"

type Component interface {
	Decode(map[string]any) (map[string]any, error)
	IsPresent(map[string]any) bool
}

type StaticProperties struct {
	Core        `mapstructure:",squash" yaml:",inline"`
	DiaryEntry  `mapstructure:",squash" yaml:",inline"`
	MoodTracker `mapstructure:",squash" yaml:",inline"`
	TimeTracker `mapstructure:",squash" yaml:",inline"`
	Timestamp   `mapstructure:",squash" yaml:",inline"`
	Rest        map[string]any `mapstructure:",remain" yaml:",inline"`
}

type DynamicProperties struct {
	Core       `mapstructure:",squash" yaml:",inline"`
	Components []Component    `mapstructure:"-" yaml:"-"`
	Rest       map[string]any `mapstructure:",remain" yaml:",inline"`
}

func (p *DynamicProperties) Decode(props map[string]any) error {
	for _, component := range p.Components {
		if component.IsPresent(props) {
			var err error
			props, err = component.Decode(props)
			if err != nil {
				return err
			}
		}
	}
	if err := mapstructure.Decode(props, p); err != nil {
		return err
	}
	return nil
}

type Core struct {
	UID     string   `yaml:"uid"`
	Aliases []string `yaml:"aliases,omitempty"`
	Tags    []string `yaml:"tags,omitempty"`
}

type DiaryEntry struct {
	Date string `yaml:"date"`
}

func (c *DiaryEntry) Decode(props map[string]any) (map[string]any, error) {
	return Decode(c, props)
}

func (*DiaryEntry) IsPresent(props map[string]any) bool {
	for _, key := range []string{"date"} {
		if _, present := props[key]; !present {
			return false
		}
	}
	return true
}

type MoodTracker struct {
	Sleep  *string `yaml:"sleep"`
	Health *string `yaml:"health"`
	Mood   *string `yaml:"mood"`
}

func (c *MoodTracker) Decode(props map[string]any) (map[string]any, error) {
	return Decode(c, props)
}

func (*MoodTracker) IsPresent(props map[string]any) bool {
	for _, key := range []string{"sleep", "health", "mood"} {
		if _, present := props[key]; !present {
			return false
		}
	}
	return true
}

type TimeTracker struct {
	Tracking *string `yaml:"tracking"`
	Duration *string `yaml:"duration"`
}

func (c *TimeTracker) Decode(props map[string]any) (map[string]any, error) {
	return Decode(c, props)
}

func (*TimeTracker) IsPresent(props map[string]any) bool {
	for _, key := range []string{"tracking", "duration"} {
		if _, present := props[key]; !present {
			return false
		}
	}
	return true
}

type Timestamp struct {
	CreatedAt *string `yaml:"created_at"`
	UpdatedAt *string `yaml:"updated_at"`
}

func (c *Timestamp) Decode(props map[string]any) (map[string]any, error) {
	return Decode(c, props)
}

func (*Timestamp) IsPresent(props map[string]any) bool {
	for _, key := range []string{"created_at", "updated_at"} {
		if _, present := props[key]; !present {
			return false
		}
	}
	return true
}

func Decode[T comparable](component *T, props map[string]any) (map[string]any, error) {
	var data = struct {
		Component *T             `mapstructure:",squash" yaml:",inline"`
		Rest      map[string]any `mapstructure:",remain" yaml:",inline"`
	}{component, nil}

	if err := mapstructure.Decode(props, &data); err != nil {
		return nil, err
	}
	return data.Rest, nil
}
