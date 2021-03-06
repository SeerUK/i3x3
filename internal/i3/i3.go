package i3

import (
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
)

// FindOutputs fetches an array of outputs from i3 via i3-msg. This will return all outputs, not
// just the active ones.
func FindOutputs() ([]Output, error) {
	var outputs []Output

	out, err := exec.Command("i3-msg", "-t", "get_outputs").Output()
	if err != nil {
		return []Output{}, err
	}

	err = json.Unmarshal(out, &outputs)
	if err != nil {
		return []Output{}, err
	}

	return outputs, nil
}

// FindWorkspaces fetches an array of workspaces from i3 via i3-msg. This will return all
// workspaces, not just the visible ones, or focused ones.
func FindWorkspaces() ([]Workspace, error) {
	var workspaces []Workspace

	out, err := exec.Command("i3-msg", "-t", "get_workspaces").Output()
	if err != nil {
		return []Workspace{}, err
	}

	err = json.Unmarshal(out, &workspaces)
	if err != nil {
		return []Workspace{}, err
	}

	return workspaces, nil
}

// MoveToWorkspace tells i3 to move the current container to the given workspace. It does not also
// switch to the workspace. Any error running the i3-msg command will be returned.
func MoveToWorkspace(workspace float64) error {
	ws := fmt.Sprintf("%v", workspace)

	return exec.Command("i3-msg", "move", "container", "to", "workspace", ws).Run()
}

// SwitchToWorkspace tells i3 to switch to the given workspace. Any error running the i3-msg command
// will be returned.
func SwitchToWorkspace(workspace float64) error {
	ws := fmt.Sprintf("%v", workspace)

	return exec.Command("i3-msg", "workspace", ws).Run()
}

// MoveWorkspaceToOutput takes a given workspace, and moves it to the given output (used for
// re-arranging workspaces on differing numbers of outputs).
func MoveWorkspaceToOutput(workspaceNum float64, outputName string) error {
	// Firstly, switch to the workspace that will be moved.
	SwitchToWorkspace(float64(workspaceNum))

	// Then move the workspace (you can't move a workspace unless you're on it).
	return exec.Command("i3-msg", "move", "workspace", "to", "output", outputName).Run()
}

// ActiveOutputsNum counts the number of active outputs in the given slice of Outputs. This could
// be, but is unlikely to be 0.
func ActiveOutputsNum(outputs []Output) float64 {
	return float64(len(ActiveOutputs(outputs)))
}

// ActiveOutputs returns the active outputs in the given slice of Outputs. This could be, but is
// unlikely to be empty.
func ActiveOutputs(outputs []Output) []Output {
	var activeOutputs []Output

	for _, output := range outputs {
		if output.Active {
			activeOutputs = append(activeOutputs, output)
		}
	}

	return activeOutputs
}

// CurrentOutputNum calculates the current "display number" that i3x3 uses internally based on
// the workspace number, and the number of outputs. We avoid trying to figure out the physical
// layout of displays because that will be both complicated, and error prone. This method works best
// if your only method of navigating workspaces is by using i3x3.
func CurrentOutputNum(workspaceNum float64, outputsNum float64) float64 {
	mod := math.Mod(workspaceNum, outputsNum)

	if mod == 0 {
		return outputsNum
	}

	return mod
}

// CurrentWorkspaceNum gets the currently focused workspace number from the given workspaces.
func CurrentWorkspaceNum(workspaces []Workspace) float64 {
	for _, workspace := range workspaces {
		if workspace.Focused {
			return float64(workspace.Num)
		}
	}

	return 1.0
}

// MaxWorkspaceNum finds the workspace with the highest number in the given slice of workspaces.
func MaxWorkspaceNum(workspaces []Workspace) float64 {
	max := 0.0

	for _, workspace := range workspaces {
		max = math.Max(max, float64(workspace.Num))
	}

	return max
}
