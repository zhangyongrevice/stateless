package stateless

import (
	"fmt"
	"strings"
)

type graph struct {
}

func (g *graph) FormatStateMachine(sm *StateMachine) string {
	var sb strings.Builder
	sb.WriteString("digraph {\n\tcompound=true;\n\tnode [shape=Mrecord];\n\trankdir=\"LR\";\n\n")
	for _, sr := range sm.stateConfig {
		if len(sr.Substates) > 0 && sr.Superstate == nil {
			sb.WriteString(g.formatOneCluster(sr))
		} else {
			sb.WriteString(g.formatOneState(sr))
		}
	}
	for _, sr := range sm.stateConfig {
		sb.WriteString(g.formatAllStateTransitions(sm, sr))
	}
	sb.WriteString("\n}")
	return sb.String()
}

func (g *graph) formatOneState(sr *stateRepresentation) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\t%s [label=\"%s", sr.State, sr.State))
	if len(sr.EntryActions) == 0 && len(sr.ExitActions) == 0 {
		sb.WriteString("\"];\n")
		return sb.String()
	}
	sb.WriteString("|")
	es := make([]string, 0, len(sr.EntryActions)+len(sr.ExitActions))
	for _, act := range sr.EntryActions {
		if act.Trigger == nil {
			es = append(es, fmt.Sprintf("enter / %s", act.Description.String()))
		}
	}
	for _, act := range sr.ExitActions {
		es = append(es, fmt.Sprintf("exit / %s", act.Description.String()))
	}
	sb.WriteString(strings.Join(es, "\n"))
	sb.WriteString("\"];\n")
	return sb.String()
}

func (g *graph) formatOneCluster(sr *stateRepresentation) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\nsubgraph cluster_%s {\n\tlabel=\"%s", sr.State, sr.State))
	if len(sr.EntryActions) > 0 || len(sr.ExitActions) > 0 {
		sb.WriteString("\\n----------")
		for _, act := range sr.EntryActions {
			if act.Trigger == nil {
				sb.WriteString("\\nentry / " + act.Description.String())
			}
		}
		for _, act := range sr.ExitActions {
			sb.WriteString("\\nexit / " + act.Description.String())
		}
	}
	sb.WriteString("\";\n")
	for _, substate := range sr.Substates {
		sb.WriteString(g.formatOneState(substate))
	}

	sb.WriteString("}\n")
	return sb.String()
}

func (g *graph) formatAllStateTransitions(sm *StateMachine, sr *stateRepresentation) string {
	var sb strings.Builder
	for _, triggers := range sr.TriggerBehaviours {
		for _, trigger := range triggers {
			switch t := trigger.(type) {
			case *ignoredTriggerBehaviour:
				sb.WriteString(g.formatOneTransition(sr.State, sr.State, t.Trigger, nil, t.Guard))
			case *reentryTriggerBehaviour:
				var actions []string
				for _, ea := range sr.EntryActions {
					if ea.Trigger == t.Trigger {
						actions = append(actions, ea.Description.String())
					}
				}
				sb.WriteString(g.formatOneTransition(sr.State, t.Destination, t.Trigger, actions, t.Guard))
			case *transitioningTriggerBehaviour:
				var actions []string
				dest := sm.stateConfig[t.Destination]
				for _, ea := range dest.EntryActions {
					if ea.Trigger == t.Trigger {
						actions = append(actions, ea.Description.String())
					}
				}
				sb.WriteString(g.formatOneTransition(sr.State, t.Destination, t.Trigger, actions, t.Guard))
			}
		}
	}
	return sb.String()
}

func (g *graph) formatOneTransition(source, destination State, trigger Trigger, actions []string, guards transitionGuard) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprint(trigger))
	if len(actions) > 0 {
		sb.WriteString(" / ")
		sb.WriteString(strings.Join(actions, ", "))
	}
	for _, info := range guards.Guards {
		if sb.Len() > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(fmt.Sprintf("[%s]", info.Description.String()))
	}
	return g.formatOneLine(fmt.Sprint(source), fmt.Sprint(destination), sb.String())
}

func (g *graph) formatOneLine(fromNodeName, toNodeName, label string) string {
	return fmt.Sprintf("\n%s -> %s [style=\"solid\", label=\"%s\"];", fromNodeName, toNodeName, label)
}