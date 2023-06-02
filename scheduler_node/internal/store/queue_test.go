package store

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/stretchr/testify/assert"
)


func TestPeek(t *testing.T) {

	// Creating multiple workflows.
	newWorkflow1 := createWorkflow()
	defer deleteWorkflow(newWorkflow1)

	newWorkflow2 := createWorkflow()
	defer deleteWorkflow(newWorkflow2)

	newWorkflow3 := createWorkflow()
	defer deleteWorkflow(newWorkflow3)
	// Completing a workflow.
	err := testStore.CompleteWorkflow(ctx, newWorkflow1.ID)
	assert.NoError(t, err)

	// We expect newWorkflow2 to be the earliest non-completed workflow.
	workflow, err := testStore.Peek(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, workflow)
	assert.Equal(t, newWorkflow2.ID, workflow.ID)
}

func TestGetQueue(t *testing.T) {
	t.Parallel()

	// Creating multiple workflows.
	newWorkflow1 := createWorkflow()
	defer deleteWorkflow(newWorkflow1)

	newWorkflow2 := createWorkflow()
	defer deleteWorkflow(newWorkflow2)

	newWorkflow3 := createWorkflow()
	defer deleteWorkflow(newWorkflow3)

	// Completing a workflow.
	err := testStore.CompleteWorkflow(ctx, newWorkflow1.ID)
	assert.NoError(t, err)

	// Retrieving the queue.
	queue, err := testStore.GetQueue(ctx)

	assert.NoError(t, err)
	assert.NotEmpty(t, queue)

	// Verifying that the completed workflow is not in the queue.
	for _, workflow := range queue {
		assert.NotEqual(t, newWorkflow1.ID, workflow.ID)
	}
}

