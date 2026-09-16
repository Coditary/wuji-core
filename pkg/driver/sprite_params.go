package driver

import (
	"fmt"
	"strings"
)

// SpriteAction identifies a preset animation type for sprite generation.
type SpriteAction string

const (
	SpriteActionIdle   SpriteAction = "idle"
	SpriteActionWalk   SpriteAction = "walk"
	SpriteActionRun    SpriteAction = "run"
	SpriteActionAttack SpriteAction = "attack"
	SpriteActionJump   SpriteAction = "jump"
	SpriteActionCast   SpriteAction = "cast"
	SpriteActionDeath  SpriteAction = "death"
	SpriteActionCustom SpriteAction = "custom"
)

// SpriteView identifies the camera/view preset for sprite generation.
type SpriteView string

const (
	SpriteViewSide       SpriteView = "side"
	SpriteViewFront      SpriteView = "front"
	SpriteViewBack       SpriteView = "back"
	SpriteViewTop        SpriteView = "top"
	SpriteViewIsometric  SpriteView = "isometric"
	SpriteViewThreeQuarter SpriteView = "three-quarter"
)

func AllSpriteActions() []SpriteAction {
	return []SpriteAction{
		SpriteActionIdle,
		SpriteActionWalk,
		SpriteActionRun,
		SpriteActionAttack,
		SpriteActionJump,
		SpriteActionCast,
		SpriteActionDeath,
		SpriteActionCustom,
	}
}

func AllSpriteViews() []SpriteView {
	return []SpriteView{
		SpriteViewSide,
		SpriteViewFront,
		SpriteViewBack,
		SpriteViewTop,
		SpriteViewIsometric,
		SpriteViewThreeQuarter,
	}
}

func (a SpriteAction) IsValid() bool {
	for _, known := range AllSpriteActions() {
		if a == known {
			return true
		}
	}
	return false
}

func (v SpriteView) IsValid() bool {
	for _, known := range AllSpriteViews() {
		if v == known {
			return true
		}
	}
	return false
}

// ParseSpriteAction normalizes an animation type name. Empty string is allowed.
func ParseSpriteAction(raw string) (SpriteAction, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	if normalized == "" {
		return "", nil
	}
	switch normalized {
	case "idle", "stand":
		return SpriteActionIdle, nil
	case "walk", "walking":
		return SpriteActionWalk, nil
	case "run", "running":
		return SpriteActionRun, nil
	case "attack", "hit":
		return SpriteActionAttack, nil
	case "jump", "jumping":
		return SpriteActionJump, nil
	case "cast", "spell", "magic":
		return SpriteActionCast, nil
	case "death", "die", "dying":
		return SpriteActionDeath, nil
	case "custom":
		return SpriteActionCustom, nil
	}
	action := SpriteAction(normalized)
	if !action.IsValid() {
		return "", fmt.Errorf("unknown sprite action %q (valid: %s)", raw, joinSpriteActions())
	}
	return action, nil
}

// ParseSpriteView normalizes a camera/view preset. Empty string is allowed.
func ParseSpriteView(raw string) (SpriteView, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	if normalized == "" {
		return "", nil
	}
	switch normalized {
	case "side", "side-view", "profile":
		return SpriteViewSide, nil
	case "front", "front-view":
		return SpriteViewFront, nil
	case "back", "rear", "back-view":
		return SpriteViewBack, nil
	case "top", "top-down", "overhead":
		return SpriteViewTop, nil
	case "iso", "isometric":
		return SpriteViewIsometric, nil
	case "3/4", "three-quarter", "threequarter", "3-quarter":
		return SpriteViewThreeQuarter, nil
	}
	view := SpriteView(normalized)
	if !view.IsValid() {
		return "", fmt.Errorf("unknown sprite view %q (valid: %s)", raw, joinSpriteViews())
	}
	return view, nil
}

// ParseSpriteDirections validates directional frame counts (4 or 8). Zero is allowed.
func ParseSpriteDirections(raw int) (int, error) {
	if raw == 0 {
		return 0, nil
	}
	if raw == 4 || raw == 8 {
		return raw, nil
	}
	return 0, fmt.Errorf("sprite directions must be 4 or 8 (got %d)", raw)
}

func joinSpriteActions() string {
	parts := make([]string, len(AllSpriteActions()))
	for i, a := range AllSpriteActions() {
		parts[i] = string(a)
	}
	return strings.Join(parts, ", ")
}

func joinSpriteViews() string {
	parts := make([]string, len(AllSpriteViews()))
	for i, v := range AllSpriteViews() {
		parts[i] = string(v)
	}
	return strings.Join(parts, ", ")
}
