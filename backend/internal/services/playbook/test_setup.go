package playbook

import (
	"path/filepath"
	"runtime"

	"github.com/jmoiron/sqlx"
)

// SetupTestPlaybook loads the test playbook JSON and initializes a PlaybookManager for tests.
// It inserts the playbook into playbook_versions and sets it as active for the given tenant.
func SetupTestPlaybook(db *sqlx.DB, tenantID string) (*Manager, error) {
	// Find the test_playbook.json file
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	// testutils is where test_playbook.json lives
	playbookPath := filepath.Join(basepath, "../../testutils/test_playbook.json")
	
	// Use EnsureActiveForTenant to load and activate
	return EnsureActiveForTenant(db, playbookPath, tenantID)
}
