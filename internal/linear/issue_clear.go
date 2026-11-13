package linear

// IssueClearField is an IssueUpdateInput field Linear accepts an explicit
// empty value for. Fields outside this set (title, stateId, ...) reject null
// with "Argument Validation Error".
type IssueClearField string

const (
	ClearProject  IssueClearField = "projectId"
	ClearAssignee IssueClearField = "assigneeId"
	ClearCycle    IssueClearField = "cycleId"
	ClearParent   IssueClearField = "parentId"
	ClearDueDate  IssueClearField = "dueDate"
	ClearEstimate IssueClearField = "estimate"
	ClearLabels   IssueClearField = "labelIds"
)

func (f IssueClearField) emptyValue() any {
	if f == ClearLabels {
		return []string{}
	}
	return nil
}
