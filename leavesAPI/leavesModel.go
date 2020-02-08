package leaves

/**
LeaveType {
	10: CasualLeave,
	20: SickLeave
}

FromDateType {
	0: FullDay,
	1: FirstHalf,
	2: SecondHalf
}

ToDateType {
	0: FullDay,
	1: FirstHalf,
	2: SecondHalf
}

LeaveStatus {
	0: Pending,
	1: Approved,
	2: Rejected,
	3: Cancelled
}
*/

// EmployeeLeave struct for define fields
type EmployeeLeave struct {
	ID          int    `json:"id"`
	Reason      string `json:"reason"`
	FromDate    string `json:"from_date"`
	ToDate      string `json:"to_date"`
	LeaveType   int    `json:"leave_type"`
	FromDateDay int    `json:"from_date_day"`
	ToDateDay   int    `json:"to_date_day"`
	TotalDays   int    `json:"total_days"`
	LeaveStatus int    `json:"leave_status"`
	AppliedBy   int    `json:"applied_by"`
}

// EmployeeLeaveInfo struct for return all leave and user details
type EmployeeLeaveInfo struct {
	ID            int    `json:"id"`
	Reason        string `json:"reason"`
	FromDate      string `json:"from_date"`
	ToDate        string `json:"to_date"`
	LeaveType     int    `json:"leave_type"`
	FromDateDay   int    `json:"from_date_day"`
	ToDateDay     int    `json:"to_date_day"`
	TotalDays     int    `json:"total_days"`
	LeaveStatus   int    `json:"leave_status"`
	UserFirstName string `json:"user_first_name"`
	UserLastName  string `json:"user_last_name"`
	UserID        int    `json:"user_id"`
}
