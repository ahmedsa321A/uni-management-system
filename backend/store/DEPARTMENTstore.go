package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"university-management/backend/models"
)

type DeprtmentStore struct {
	db *sql.DB
}

// function to create department recod return department id anr error message
func (d *DeprtmentStore) GetAll() ([]*models.Department, error) {

	query := `
          SELECT department_id, department_name, faculty_id 
          FROM DEPARTMENTS 
         `

	rows, err := d.db.Query(query)

	if err != nil {
		return nil, errors.New("Error getting departments")
	}
	defer rows.Close()
	var departments []*models.Department
	for rows.Next() {
		var department models.Department
		err := rows.Scan(
			&department.DepartmentID,
			&department.DepartmentName,
			&department.FacultyID,
		)
		if err != nil {
			return nil, errors.New("Error scanning department")
		}
		departments = append(departments, &department)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return departments, nil
}
func (d *DeprtmentStore) CreateDepartment(department *models.Department) (int64, error) {
	// assume 0 cannot be a dept_id as it starts from 1
	if department == nil || department.DepartmentName == "" {
		return 0, errors.New("department name is required")
	}

	query := `
        INSERT INTO DEPARTMENTS (department_name, faculty_id)
        VALUES (?, ?)
    `
	result, err := d.db.Exec(query, department.DepartmentName, department.FacultyID)
	if err != nil {
		return 0, fmt.Errorf("failed to create department: %w", err)
	}

	newID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return newID, nil
}

func (d *DeprtmentStore) GetDepartmentById(id int64) (*models.Department, error) {
	if id <= 0 {
		return nil, errors.New("invalid department id")
	}
	query := `
          SELECT department_id, department_name, faculty_id 
          FROM DEPARTMENTS 
          WHERE  department_id = ?`
	var department models.Department
	err := d.db.QueryRow(query, id).Scan(
		&department.DepartmentID,
		&department.DepartmentName,
		&department.FacultyID,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("department not found")
		}
		return nil, fmt.Errorf("failed to retrieve department: %w", err)
	}

	return &department, nil
}

func (d *DeprtmentStore) GetDepartmentByName(name string) (*models.Department, error) {
	if name == "" {
		return nil, errors.New("invalid department name")
	}
	query := `
        SELECT department_id, department_name, faculty_id
        FROM DEPARTMENTS
        WHERE department_name = ?`

	var department models.Department
	err := d.db.QueryRow(query, name).Scan(&department.DepartmentID,
		&department.DepartmentName,
		&department.FacultyID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("department not found")
		}
		return nil, fmt.Errorf("failed to retrieve department: %w", err)
	}
	return &department, nil
}

func (d *DeprtmentStore) DeleteDepartmentById(id int64) error {
	if id <= 0 {
		return errors.New("invalid department id")
	}
	query := `DELETE FROM DEPARTMENTS WHERE department_id = ?`
	err := d.db.QueryRow(query, id).Scan(&sql.ErrNoRows)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("department not found")
		}
		return fmt.Errorf("failed to delete department: %w", err)
	}
	return nil
}

func (d *DeprtmentStore) DeleteDepartmentByName(name string) error {
	if name == "" {
		return errors.New("invalid department name")
	}
	query := `DELETE FROM DEPARTMENTS WHERE department_name = ?`
	err := d.db.QueryRow(query, name).Scan(&sql.ErrNoRows)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("department not found")
		}
		return fmt.Errorf("failed to delete department: %w", err)
	}
	return nil
}

func (d *DeprtmentStore) UpdateDepartmentByID(department *models.Department, newName string, newFacultyID *int) (*models.Department, error) {

	if department == nil || department.DepartmentID <= 0 {
		return nil, errors.New("invalid department or department ID")
	}
	if newName == "" {
		return nil, errors.New("department name is required")
	}

	// Start transaction
	tx, err := d.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `
        UPDATE DEPARTMENTS
        SET department_name = ?, faculty_id = ?
        WHERE department_id = ?
    `
	result, err := tx.Exec(query, newName, newFacultyID, department.DepartmentID)
	if err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return nil, errors.New("cannot update department; invalid faculty_id")
		}
		return nil, fmt.Errorf("failed to update department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, errors.New("department not found")
	}

	query = `
        SELECT department_id, department_name, faculty_id
        FROM DEPARTMENTS
        WHERE department_id = ?
    `
	var updatedDept models.Department
	err = tx.QueryRow(query, department.DepartmentID).Scan(
		&updatedDept.DepartmentID,
		&updatedDept.DepartmentName,
		&updatedDept.FacultyID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated department: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &updatedDept, nil
}
