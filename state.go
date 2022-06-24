package main

type ButtonState int

const (
	ButtonPurple ButtonState = iota
	ButtonBlue
	ButtonGreen
	ButtonYellow
	ButtonOrange
	ButtonRed
	ButtonDead
)

var ButtonStates = []ButtonState{
	ButtonPurple, ButtonBlue, ButtonGreen, ButtonYellow, ButtonOrange, ButtonRed,
}

// map role id to button state
var roleToButtonState = map[string]ButtonState{}

func (b ButtonState) Role() string {
	switch b {
	case ButtonPurple:
		return "purple"
	case ButtonBlue:
		return "blue"
	case ButtonGreen:
		return "green"
	case ButtonYellow:
		return "yellow"
	case ButtonOrange:
		return "orange"
	case ButtonRed:
		return "red"
	case ButtonDead:
		return "invalid"
	}

	return ""
}

func (b ButtonState) Channel() string {
	switch b {
	case ButtonPurple:
		return "button-🟣"
	case ButtonBlue:
		return "button-🔵"
	case ButtonGreen:
		return "button-🟢"
	case ButtonYellow:
		return "button-🟡"
	case ButtonOrange:
		return "button-🟠"
	case ButtonRed:
		return "button-🔴"
	case ButtonDead:
		return "button-⚫"
	}

	return ""
}

// returns color of the state in hex
func (b ButtonState) Color() int {
	switch b {
	case ButtonPurple:
		return 0x9b59b6
	case ButtonBlue:
		return 0x3498db
	case ButtonGreen:
		return 0x57f287
	case ButtonYellow:
		return 0xfee75c
	case ButtonOrange:
		return 0xe67e22
	case ButtonRed:
		return 0xed4245
	default:
		return 0x000000
	}
}

func (b ButtonState) Next() ButtonState {
	switch b {
	case ButtonPurple:
		return ButtonBlue
	case ButtonBlue:
		return ButtonGreen
	case ButtonGreen:
		return ButtonYellow
	case ButtonYellow:
		return ButtonOrange
	case ButtonOrange:
		return ButtonRed
	case ButtonRed:
		return ButtonDead
	case ButtonDead:
		return ButtonPurple
	default:
		return ButtonDead
	}
}

func ButtonStateFromState(s string) ButtonState {
	switch s {
	case "button-🟣":
		return ButtonPurple
	case "button-🔵":
		return ButtonBlue
	case "button-🟢":
		return ButtonGreen
	case "button-🟡":
		return ButtonYellow
	case "button-🟠":
		return ButtonOrange
	case "button-🔴":
		return ButtonRed
	default:
		return ButtonDead
	}
}

func ButtonStateFromRole(s string) ButtonState {
	switch s {
	case "purple":
		return ButtonPurple
	case "blue":
		return ButtonBlue
	case "green":
		return ButtonGreen
	case "yellow":
		return ButtonYellow
	case "orange":
		return ButtonOrange
	case "red":
		return ButtonRed
	default:
		return ButtonDead
	}
}
