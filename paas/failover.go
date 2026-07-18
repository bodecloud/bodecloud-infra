package paas

import (
	"encoding/json"
	"fmt"
)

func (failover *FailoverSpec) MarshalJSON() ([]byte, error) {
	if failover == nil {
		return []byte("null"), nil
	}
	return json.Marshal(serializeFailoverSpec(failover))
}

func (failover *FailoverSpec) UnmarshalJSON(data []byte) error {
	if failover == nil {
		return fmt.Errorf("nil FailoverSpec")
	}
	parsed, err := parseFailoverSpecMapBytes(data)
	if err != nil {
		return err
	}
	*failover = *parsed
	return nil
}

func failoverSpecFromMap(mapped map[string]interface{}) (*FailoverSpec, error) {
	if len(mapped) == 0 {
		return nil, fmt.Errorf("failover must be a map")
	}
	failover := &FailoverSpec{}
	if _, ok := mapped["enabled"]; ok {
		failover.Enabled = toBool(mapped["enabled"])
		failover.EnabledSet = true
	}
	failover.Type = toString(mapped["type"])
	failover.Port = toString(mapped["port"])
	failover.HealthcheckPath = toString(mapped["healthcheck_path"])
	failover.HealthcheckInterval = toString(mapped["healthcheck_interval"])
	failover.MaxRetries = toInt(mapped["max_retries"])
	failover.RedeployOnPeer = toBool(mapped["redeploy_on_peer"])
	failover.Singleton = toBool(mapped["singleton"])
	failover.SingletonElection = toString(mapped["singleton_election"])
	failover.Strategy = toString(mapped["strategy"])
	failover.PreferLocal = toBool(mapped["prefer_local"])
	if nodesRaw, ok := asMap(mapped["nodes"]); ok && len(nodesRaw) > 0 {
		failover.Nodes = map[string]*FailoverNode{}
		for name, rawNode := range nodesRaw {
			nodeMap, ok := asMap(rawNode)
			if !ok {
				continue
			}
			node, err := failoverNodeFromMap(nodeMap)
			if err != nil {
				return nil, fmt.Errorf("failed to parse failover node %s: %w", name, err)
			}
			failover.Nodes[name] = node
		}
	}
	if failover.Extensions == nil {
		failover.Extensions = map[string]interface{}{}
	}
	for key, raw := range mapped {
		switch key {
		case "enabled", "type", "port", "healthcheck_path", "healthcheck_interval", "max_retries", "redeploy_on_peer", "singleton", "singleton_election", "strategy", "prefer_local", "nodes", "extensions":
			continue
		default:
			failover.Extensions[key] = deepCopyValue(raw)
		}
	}
	if ext, ok := asMap(mapped["extensions"]); ok {
		for key, raw := range ext {
			failover.Extensions[key] = deepCopyValue(raw)
		}
	}
	if len(failover.Extensions) == 0 {
		failover.Extensions = nil
	}
	return failover, nil
}

func parseFailoverSpecMapBytes(data []byte) (*FailoverSpec, error) {
	var mapped map[string]interface{}
	if err := json.Unmarshal(data, &mapped); err != nil {
		return nil, err
	}
	return failoverSpecFromMap(mapped)
}

func failoverNodeFromMap(mapped map[string]interface{}) (*FailoverNode, error) {
	if len(mapped) == 0 {
		return &FailoverNode{}, nil
	}
	node := &FailoverNode{
		Status:     toString(mapped["status"]),
		LastSeen:   toString(mapped["last_seen"]),
		Priority:   toInt(mapped["priority"]),
		URL:        toString(mapped["url"]),
		Weight:     toInt(mapped["weight"]),
		Extensions: map[string]interface{}{},
	}
	for key, raw := range mapped {
		switch key {
		case "status", "last_seen", "priority", "url", "weight", "extensions":
			continue
		default:
			node.Extensions[key] = deepCopyValue(raw)
		}
	}
	if ext, ok := asMap(mapped["extensions"]); ok {
		for key, raw := range ext {
			node.Extensions[key] = deepCopyValue(raw)
		}
	}
	if len(node.Extensions) == 0 {
		node.Extensions = nil
	}
	return node, nil
}
