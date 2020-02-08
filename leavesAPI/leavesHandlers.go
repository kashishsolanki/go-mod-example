package leaves

import (
	"encoding/json"
	"fmt"
	"net/http"

	users "github.com/kashishsolanki/dt-hrms-golang/userAPI"

	"github.com/gorilla/mux"
	"github.com/kashishsolanki/dt-hrms-golang/db"
)

// ApplyLeaves will apply leaves for employee to management
//
// Request Type: POST
//
// URL - /v1/api/leaves/apply
func ApplyLeaves(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Apply leaves..")

	// Set content type returned to JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var employeeLeave EmployeeLeave
	_ = json.NewDecoder(r.Body).Decode(&employeeLeave)

	db, err := db.Open()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	rows, err := db.Query(`SELECT COUNT(*) as count FROM employee_leave`)
	if err != nil {
		http.Error(w, "Error getting applied leaves count", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			http.Error(w, "Error getting applied leaves count", http.StatusInternalServerError)
			return
		}
	}
	empLeaveID := count + 1
	fmt.Println("Count : ", count)
	empLeaveInsert, err := db.Prepare(`
		INSERT INTO employee_leave
			(
				id,
				reason,
				from_date,
				to_date,
				leave_type,
				from_date_day,
				to_date_day,
				total_days,
				leave_status,
				applied_by
			)
		VALUES (?,?,?,?,?,?,?,?,?,?) `)

	if err != nil {
		panic(err.Error())
	}
	fmt.Println("employee_leave_reason :: ", employeeLeave.Reason, empLeaveID)

	userID, err := getUserID(r)
	fmt.Println("Token User ID :: ", userID)
	if err != nil {
		fmt.Println("Eror while getting email of user from token : ", err)
		http.Error(w, "Error while getting email of user from token", http.StatusInternalServerError)
		return
	}

	_, insertErr := empLeaveInsert.Exec(
		empLeaveID,
		employeeLeave.Reason,
		employeeLeave.FromDate,
		employeeLeave.ToDate,
		employeeLeave.LeaveType,
		employeeLeave.FromDateDay,
		employeeLeave.ToDateDay,
		employeeLeave.TotalDays,
		employeeLeave.LeaveStatus,
		userID)

	if insertErr != nil {
		fmt.Println(insertErr)
		http.Error(w, "Error applying leaves", http.StatusInternalServerError)
		return
	}

	// Get user info for manage user leave and maintain its data
	userInfo, userErr := users.GetUserInfo(userID, "id")
	if userErr != nil {
		fmt.Println("Error while getting user data :: ", userErr)
		http.Error(w, "Error while getting  user data ", http.StatusInternalServerError)
		return
	}

	// Update leaves as per users applied numbers of leaves
	if employeeLeave.LeaveType == 10 {
		userInfo.CasualLeave -= employeeLeave.TotalDays
	} else if employeeLeave.LeaveType == 20 {
		userInfo.SickLeave -= employeeLeave.TotalDays
	}
	leaveErr := users.UpdateUserLeave(userInfo.ID, userInfo.SickLeave, userInfo.CasualLeave)
	if leaveErr != nil {
		http.Error(w, "Error in updating employee leaves", http.StatusInternalServerError)
		return
	}

	// send mail for leave management
	sendLeaveMail(userInfo, employeeLeave)

	// Write OK back to client
	w.WriteHeader(http.StatusOK)

}

// GetAllLeaves will get list of all leaves
//
// Request Type: GET
//
// URL - /v1/api/leaves
func GetAllLeaves(w http.ResponseWriter, r *http.Request) {

	var employeeLeaves []EmployeeLeaveInfo

	db, err := db.Open()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	queryValues := r.URL.Query()
	userIDParam := queryValues.Get("user_id")
	fmt.Println(userIDParam)

	isUserAdmin, isUserAdminErr := isUserAdmin(r)
	if isUserAdminErr != nil {
		// handle this error
		fmt.Println(isUserAdminErr)
		http.Error(w, "Error in check user is admin or not", http.StatusInternalServerError)
		return
	}

	q := `SELECT 
			id,
			reason,
			from_date,
			to_date,
			leave_type,
			from_date_day,
			to_date_day,
			total_days,
			leave_status,
			applied_by
		FROM employee_leave`

	if userIDParam != "" {
		if isUserAdmin {
			q += " where applied_by = " + userIDParam + " ORDER BY from_date ASC"
		} else {
			q += " where applied_by = " + userIDParam + " ORDER BY from_date ASC"
		}
	} else {
		q += " where leave_status = 0 and DATE(from_date) > DATE(NOW()) ORDER BY from_date ASC"
	}

	fmt.Println("QUery :: ", q)

	rows, err := db.Query(q)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error getting employees leaves", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	employeeLeave := EmployeeLeaveInfo{}
	for rows.Next() {
		err = rows.Scan(
			&employeeLeave.ID,
			&employeeLeave.Reason,
			&employeeLeave.FromDate,
			&employeeLeave.ToDate,
			&employeeLeave.LeaveType,
			&employeeLeave.FromDateDay,
			&employeeLeave.ToDateDay,
			&employeeLeave.TotalDays,
			&employeeLeave.LeaveStatus,
			&employeeLeave.UserID)

		if err != nil {
			// handle this error
			fmt.Println(err)
			http.Error(w, "Error while get all data", http.StatusInternalServerError)
			return
		}

		fmt.Println("Token User ID :: ", employeeLeave.UserID)

		// Get user info for manage user leave and maintain its data
		userInfo, userErr := users.GetUserInfo(fmt.Sprint(employeeLeave.UserID), "id")
		if userErr != nil {
			fmt.Println("Error while getting user data :: ", userErr)
			http.Error(w, "Error while getting  user data ", http.StatusInternalServerError)
			return
		}

		employeeLeave.UserFirstName = userInfo.FirstName
		employeeLeave.UserLastName = userInfo.LastName
		employeeLeaves = append(employeeLeaves, employeeLeave)
	}

	jsonResponse, jsonError := json.Marshal(employeeLeaves)

	if jsonError != nil {
		fmt.Println(jsonError)
		http.Error(w, "error in retriving data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// UpdateLeave will update employee leave details
//
// Request Type: PUT
//
// URL - /v1/api/leaves/id
func UpdateLeave(w http.ResponseWriter, r *http.Request) {
	leaveID := mux.Vars(r)["id"]
	fmt.Println(leaveID)

	if leaveID == "" {
		http.Error(w, "Leave id is missing", http.StatusBadRequest)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	employeeLeaveInfo, leaveInfoErr := GetEmployeeLeaveInfo(leaveID, "id")
	if leaveInfoErr != nil {
		fmt.Println(err)
		http.Error(w, "Error while getting employee leave data : ", http.StatusInternalServerError)
		return
	}

	var employeeLeave EmployeeLeave
	_ = json.NewDecoder(r.Body).Decode(&employeeLeave)
	empLeaveUpdate, err := db.Prepare(`
		UPDATE employee_leave
			SET
				reason = ?,
				from_date = ?,
				to_date = ?,
				from_date_day = ?,
				to_date_day = ?,
				total_days = ?,
				leave_status = ?
			WHERE
				id = ? `)

	if err != nil {
		panic(err.Error())
	}
	_, empLeaveErr := empLeaveUpdate.Exec(
		employeeLeave.Reason,
		employeeLeave.FromDate,
		employeeLeave.ToDate,
		employeeLeave.FromDateDay,
		employeeLeave.ToDateDay,
		employeeLeave.TotalDays,
		employeeLeave.LeaveStatus,
		leaveID,
	)

	if empLeaveErr != nil {
		http.Error(w, "Error while update leaves", http.StatusInternalServerError)
		return
	}
	userEmail, err := getUserEmail(r)
	if err != nil {
		fmt.Println("Eror while getting email of user from token : ", err)
		http.Error(w, "Error while getting  email of user from token ", http.StatusInternalServerError)
		return
	}
	// Get user info for manage user leave and maintain its data
	userInfo, userErr := users.GetUserInfo(userEmail, "email")
	if userErr != nil {
		fmt.Println("Error while getting user data :: ", userErr)
		http.Error(w, "Error while getting user data ", http.StatusInternalServerError)
		return
	}

	fmt.Println(userInfo.SickLeave, employeeLeaveInfo.TotalDays, employeeLeave.TotalDays)
	// Update leaves as per users applied numbers of leaves
	if employeeLeave.LeaveType == 10 {
		userInfo.CasualLeave = (userInfo.CasualLeave + employeeLeaveInfo.TotalDays) - employeeLeave.TotalDays
	} else if employeeLeave.LeaveType == 20 {
		userInfo.SickLeave = (userInfo.SickLeave + employeeLeaveInfo.TotalDays) - employeeLeave.TotalDays
	}
	fmt.Println("Leaves :: sick : ", userInfo.SickLeave, ", casual : ", userInfo.CasualLeave)
	leaveErr := users.UpdateUserLeave(userInfo.ID, userInfo.SickLeave, userInfo.CasualLeave)
	if leaveErr != nil {
		http.Error(w, "Error in updating employee leaves", http.StatusInternalServerError)
		return
	}

	// send mail for leave management
	sendLeaveMail(userInfo, employeeLeave)

	// fmt.Println(empLeaveResponse, empLeaveErr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(employeeInfo)
}

// ManageLeave will approve or cancel leave of employee
//
// Request Type: PUT
//
// URL - /v1/api/leaves/manage/id
func ManageLeave(w http.ResponseWriter, r *http.Request) {
	leaveID := mux.Vars(r)["id"]
	fmt.Println(leaveID)

	if leaveID == "" {
		http.Error(w, "Leave id is missing", http.StatusBadRequest)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	var employeeLeave EmployeeLeave
	_ = json.NewDecoder(r.Body).Decode(&employeeLeave)
	empLeaveUpdate, err := db.Prepare(`
		UPDATE employee_leave
			SET
				leave_status = ?
			WHERE
				id = ? `)

	if err != nil {
		panic(err.Error())
	}
	empLeaveResponse, empLeaveErr := empLeaveUpdate.Exec(
		employeeLeave.LeaveStatus,
		leaveID,
	)

	if empLeaveErr != nil {
		http.Error(w, "Error while update leave", http.StatusInternalServerError)
		return
	}
	if employeeLeave.LeaveStatus != 0 {
		fmt.Println(empLeaveResponse)
		employeeLeaveInfo, leaveErr := GetEmployeeLeaveInfo(leaveID, "id")
		if leaveErr != nil {
			fmt.Println("Error while getting employee leave data : ", err)
			http.Error(w, "Error while getting employee leave data", http.StatusInternalServerError)
			return
		}

		userInfo, userErr := users.GetUserInfo(fmt.Sprint(employeeLeaveInfo.AppliedBy), "id")
		if userErr != nil {
			fmt.Println("Error while getting user data :: ", userErr)
			http.Error(w, "Error while getting user data ", http.StatusInternalServerError)
			return
		}
		if employeeLeave.LeaveStatus == 2 || employeeLeave.LeaveStatus == 3 {
			// Update leaves as per users applied numbers of leaves
			if employeeLeaveInfo.LeaveType == 10 {
				userInfo.CasualLeave = (userInfo.CasualLeave + employeeLeaveInfo.TotalDays)
			} else if employeeLeaveInfo.LeaveType == 20 {
				userInfo.SickLeave = (userInfo.SickLeave + employeeLeaveInfo.TotalDays)
			}
			updateLeaveErr := users.UpdateUserLeave(userInfo.ID, userInfo.SickLeave, userInfo.CasualLeave)
			if updateLeaveErr != nil {
				http.Error(w, "Error in updating employee leaves", http.StatusInternalServerError)
				return
			}
		}
		sendManageLeaveMail(userInfo, employeeLeave.LeaveStatus, employeeLeaveInfo)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(employeeInfo)
}
