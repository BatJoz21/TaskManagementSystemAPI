package models

import "taskmanagementsystem.localhost/tmsapi/database"

type Dashboard struct {
	TotalUser      int `json:"total_user"`
	TotalTask      int `json:"total_task"`
	TaskInProgress int `json:"task_in_progress"`
	TaskCancelled  int `json:"task_cancelled"`
	TaskComplete   int `json:"task_complete"`
}

func GetDataForDashboard() (*Dashboard, error) {
	query := `SELECT
    	(SELECT COUNT(*) FROM theusers) as total_user,
    	(SELECT COUNT(*) FROM tasks) as total_task,
    	(SELECT COUNT(*) FROM tasks WHERE status_id = 1) as task_in_progress,
    	(SELECT COUNT(*) FROM tasks WHERE status_id = 2) as task_cancelled,
    	(SELECT COUNT(*) FROM tasks WHERE status_id = 3) as task_complete`
	row := database.DB.QueryRow(query)

	var data Dashboard
	err := row.Scan(&data.TotalUser, &data.TotalTask, &data.TaskInProgress, &data.TaskCancelled, &data.TaskComplete)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
