package seed

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

type Step struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}
type Module struct {
	Code  string `json:"code"`
	Title string `json:"title"`
	Steps []Step `json:"steps"`
}
type Payload struct {
	Modules []Module `json:"modules"`
}

// Run reads internal/seed/algorithm.json and seeds modules/steps exactly as provided.
func Run(db *sqlx.DB) error {
	here, _ := os.Getwd()
	path := filepath.Join(here, "internal", "seed", "algorithm.json")
	b, err := os.ReadFile(path)
	if err != nil { return err }
	var p Payload
	if err := json.Unmarshal(b, &p); err != nil { return err }

	tx := db.MustBegin()
	// insert modules
	for i, m := range p.Modules {
		tx.Exec(`INSERT INTO checklist_modules (code,title,sort_order)
		         VALUES ($1,$2,$3)
		         ON CONFLICT (code) DO UPDATE SET title=$2, sort_order=$3`,
			m.Code, m.Title, i+1)
	}
	// map module ids
	type Row struct{ ID, Code string }
	var rows []Row
	tx.Select(&rows, `SELECT id, code FROM checklist_modules`)
	idBy := map[string]string{}
	for _, r := range rows { idBy[r.Code]=r.ID }

	// insert steps
	for _, m := range p.Modules {
		for j, s := range m.Steps {
			tx.Exec(`INSERT INTO checklist_steps (module_id, code, title, sort_order)
			         VALUES ($1,$2,$3,$4)
			         ON CONFLICT (code) DO UPDATE SET title=$3, sort_order=$4`,
				idBy[m.Code], s.Code, s.Title, j+1)
		}
	}
	return tx.Commit()
}
