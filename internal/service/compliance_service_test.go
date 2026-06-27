package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComplianceServiceDisabled(t *testing.T) {
	svc := NewComplianceService(nil)

	assert.True(t, svc.IsAcknowledged())

	status, err := svc.Status()
	require.NoError(t, err)
	assert.True(t, status.Acknowledged)
	assert.Equal(t, "disabled", status.Version)

	require.NoError(t, svc.Acknowledge(AcknowledgeRequest{}))
}
