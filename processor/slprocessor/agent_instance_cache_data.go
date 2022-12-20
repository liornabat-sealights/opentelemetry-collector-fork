package slprocessor

import "fmt"

type AgentInstanceCacheData struct {
	agentInstanceId string
}

func NewAgentInstanceCacheData() *AgentInstanceCacheData {
	return &AgentInstanceCacheData{
		agentInstanceId: "",
	}
}

func (d *AgentInstanceCacheData) SetAgentInstanceId(agentInstanceId string) *AgentInstanceCacheData {
	d.agentInstanceId = agentInstanceId
	return d
}

func (d *AgentInstanceCacheData) Validate() error {
	if d.agentInstanceId == "" {
		return fmt.Errorf("Agent instance id is empty")
	}

	return nil
}

var _ model = (*AgentInstanceCacheData)(nil)
