package mappers

import "github.com/shhac/lin/internal/linear"

// MapDocumentParents renders each parent Linear reports for a document after
// an update, keyed by kind ("project", "team", "initiative", "issue"). A
// document has exactly one parent, so a well-formed result has one entry.
func MapDocumentParents(doc linear.DocumentUpdateDocumentUpdateDocumentPayloadDocument) map[string]map[string]any {
	parents := map[string]map[string]any{}
	if doc.Project != nil {
		parents["project"] = map[string]any{"id": doc.Project.Id, "name": doc.Project.Name, "slugId": doc.Project.SlugId}
	}
	if doc.Team != nil {
		parents["team"] = map[string]any{"id": doc.Team.Id, "key": doc.Team.Key, "name": doc.Team.Name}
	}
	if doc.Initiative != nil {
		parents["initiative"] = map[string]any{"id": doc.Initiative.Id, "name": doc.Initiative.Name}
	}
	if doc.Issue != nil {
		parents["issue"] = map[string]any{"id": doc.Issue.Id, "identifier": doc.Issue.Identifier}
	}
	return parents
}
