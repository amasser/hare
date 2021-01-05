package ram

import (
	"errors"
)

var (
	ErrNoTable      = errors.New("hare: table with that name does not exist")
	ErrTableExists  = errors.New("hare: table with that name already exists")
	ErrNoRecord     = errors.New("hare: record with that id does not exist")
	ErrRecordExists = errors.New("hare: record with that id already exists")
)

type Ram struct {
	tables map[string]*table
}

func NewRam(seedData map[string]map[int]string) (*Ram, error) {
	var ram Ram

	if err := ram.init(seedData); err != nil {
		return nil, err
	}

	return &ram, nil
}

func (ram *Ram) Close() error {
	ram.tables = nil

	return nil
}

func (ram *Ram) CreateTable(tableName string) error {
	if ram.TableExists(tableName) {
		return ErrTableExists
	}

	ram.tables[tableName] = newTable()

	return nil
}

func (ram *Ram) DeleteRec(tableName string, id int) error {
	table, err := ram.getTable(tableName)
	if err != nil {
		return err
	}

	if err = table.deleteRec(id); err != nil {
		return err
	}

	return nil
}

func (ram *Ram) GetLastID(tableName string) (int, error) {
	table, err := ram.getTable(tableName)
	if err != nil {
		return 0, err
	}

	return table.getLastID(), nil
}

func (ram *Ram) IDs(tableName string) ([]int, error) {
	table, err := ram.getTable(tableName)
	if err != nil {
		return nil, err
	}

	return table.ids(), nil
}

func (ram *Ram) InsertRec(tableName string, id int, rec []byte) error {
	table, err := ram.getTable(tableName)
	if err != nil {
		return err
	}

	if table.recExists(id) {
		return ErrRecordExists
	}

	table.writeRec(id, rec)

	return nil
}

func (ram *Ram) ReadRec(tableName string, id int) ([]byte, error) {
	table, err := ram.getTable(tableName)
	if err != nil {
		return nil, err
	}

	rec, err := table.readRec(id)
	if err != nil {
		return nil, err
	}

	return rec, err
}

func (ram *Ram) RemoveTable(tableName string) error {
	if !ram.TableExists(tableName) {
		return ErrNoTable
	}

	delete(ram.tables, tableName)

	return nil
}

func (ram *Ram) TableExists(tableName string) bool {
	_, ok := ram.tables[tableName]

	return ok
}

func (ram *Ram) TableNames() []string {
	var names []string

	for k := range ram.tables {
		names = append(names, k)
	}

	return names
}

func (ram *Ram) UpdateRec(tableName string, id int, rec []byte) error {
	table, err := ram.getTable(tableName)
	if err != nil {
		return err
	}

	if !table.recExists(id) {
		return ErrNoRecord
	}

	table.writeRec(id, rec)

	return nil
}

//******************************************************************************
// UNEXPORTED METHODS
//******************************************************************************

func (ram *Ram) getTable(tableName string) (*table, error) {
	table, ok := ram.tables[tableName]
	if !ok {
		return nil, ErrNoTable
	}

	return table, nil
}

func (ram *Ram) getTables() ([]string, error) {
	var tableNames []string

	for name := range ram.tables {
		tableNames = append(tableNames, name)
	}

	return tableNames, nil
}

func (ram *Ram) init(seedData map[string]map[int]string) error {
	ram.tables = make(map[string]*table)

	for tableName, tableData := range seedData {
		ram.tables[tableName] = newTable()

		for id, rec := range tableData {
			if err := ram.InsertRec(tableName, id, []byte(rec)); err != nil {
				return err
			}
		}
	}

	return nil
}