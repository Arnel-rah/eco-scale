package scheduler

import (
	"strings"

	"github.com/Arnel-rah/eco-scale/config"
	"github.com/docker/docker/api/types"
)

type ActionType string

const (
	ActionStop  ActionType = "STOP"
	ActionScale ActionType = "SCALE"
	ActionNone  ActionType = "NONE"
)

type TargetContainer struct {
	ID       string
	Name     string
	Policy   config.ContainerPolicy
	Required ActionType
}

func AnalyzeSystem(policies []config.ContainerPolicy, activeContainers []types.Container, alertMode bool) []TargetContainer {
	var targets []TargetContainer

	if !alertMode {
		return targets
	}

	for _, c := range activeContainers {
		cName := ""
		if len(c.Names) > 0 {
			cName = strings.TrimPrefix(c.Names[0], "/")
		}

		for _, p := range policies {
			if p.ContainerName == cName {
				action := ActionScale
				if p.CPULimit == 0 {
					action = ActionStop
				}

				targets = append(targets, TargetContainer{
					ID:       c.ID[:10],
					Name:     cName,
					Policy:   p,
					Required: action,
				})
				break
			}
		}
	}

	return targets
}
