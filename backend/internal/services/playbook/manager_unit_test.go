package playbook

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager_NodeWorldID_Unit(t *testing.T) {
	m := &Manager{
		NodeWorlds: map[string]string{
			"node1": "worldA",
			"node2": "worldB",
		},
	}

	assert.Equal(t, "worldA", m.NodeWorldID("node1"))
	assert.Equal(t, "worldB", m.NodeWorldID("node2"))
	assert.Empty(t, m.NodeWorldID("unknown"))
}

func TestManager_GetNodesByWorld_Unit(t *testing.T) {
	m := &Manager{
		NodeWorlds: map[string]string{
			"n1": "w1",
			"n2": "w1",
			"n3": "w2",
		},
	}

	nodesW1 := m.GetNodesByWorld("w1")
	assert.Len(t, nodesW1, 2)
	assert.Contains(t, nodesW1, "n1")
	assert.Contains(t, nodesW1, "n2")

	nodesW2 := m.GetNodesByWorld("w2")
	assert.Len(t, nodesW2, 1)
	assert.Equal(t, "n3", nodesW2[0])

	nodesEmpty := m.GetNodesByWorld("wUnknown")
	assert.Empty(t, nodesEmpty)
}

func TestManager_NodeDefinition_Unit(t *testing.T) {
	node1 := Node{ID: "n1", Type: "step"}
	m := &Manager{
		Nodes: map[string]Node{
			"n1": node1,
		},
	}

	n, ok := m.NodeDefinition("n1")
	assert.True(t, ok)
	assert.Equal(t, "n1", n.ID)

	_, ok = m.NodeDefinition("unknown")
	assert.False(t, ok)
}

func TestRequirementsExpansion(t *testing.T) {
	rawJSON := `{
		"playbook_id": "pb1",
		"version": "1.0",
		"worlds": [
			{
				"id": "w1",
				"nodes": [
					{
						"id": "n1_course",
						"type": "course",
						"title": {"en": "Course Node"},
						"requirements": {
							"courseId": "c_123"
						}
					},
					{
						"id": "n2_pay",
						"type": "payment",
						"title": {"en": "Pay Node"},
						"requirements": {
							"amount": 5000,
							"currency": "KZT"
						}
					}
				]
			}
		]
	}`

	var pb Playbook
	err := json.Unmarshal([]byte(rawJSON), &pb)
	assert.NoError(t, err)

	nodes, _ := indexNodes(pb)
	
	// Check Course Node
	courseNode, ok := nodes["n1_course"]
	assert.True(t, ok)
	assert.Equal(t, "c_123", courseNode.Requirements.CourseID)

	// Check Payment Node
	payNode, ok := nodes["n2_pay"]
	assert.True(t, ok)
	assert.Equal(t, 5000, payNode.Requirements.Amount)
	assert.Equal(t, "KZT", payNode.Requirements.Currency)
}
