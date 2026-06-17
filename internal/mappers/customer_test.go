package mappers

import (
	"testing"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/ptr"
)

func TestMapCustomerSummary(t *testing.T) {
	got := MapCustomerSummary(linear.CustomerSummaryFields{
		Id:                   "c1",
		Name:                 "Acme Corp",
		SlugId:               "acme-corp",
		Url:                  "https://linear.app/acme/customer/acme-corp",
		Revenue:              ptr.To(50000),
		ApproximateNeedCount: 4,
		Owner:                &linear.CustomerSummaryFieldsOwnerUser{Id: "u1", Name: "Dana"},
		Status:               linear.CustomerSummaryFieldsStatusCustomerStatus{Id: "s1", Name: "Active"},
		Tier:                 &linear.CustomerSummaryFieldsTierCustomerTier{Id: "t1", DisplayName: "Enterprise"},
	})

	if got["name"] != "Acme Corp" {
		t.Errorf("name = %v", got["name"])
	}
	if got["revenue"] != 50000 {
		t.Errorf("revenue = %v", got["revenue"])
	}
	tier, ok := got["tier"].(map[string]any)
	if !ok || tier["displayName"] != "Enterprise" {
		t.Errorf("tier = %v", got["tier"])
	}
}

func TestMapCustomerSummary_OmitsNilOptionals(t *testing.T) {
	got := MapCustomerSummary(linear.CustomerSummaryFields{
		Id:     "c1",
		Name:   "Acme Corp",
		Status: linear.CustomerSummaryFieldsStatusCustomerStatus{Id: "s1", Name: "Active"},
	})

	if _, ok := got["revenue"]; ok {
		t.Error("revenue should be absent when nil")
	}
	if _, ok := got["owner"]; ok {
		t.Error("owner should be absent when nil")
	}
	if _, ok := got["tier"]; ok {
		t.Error("tier should be absent when nil")
	}
}

func TestMapCustomerNeedSummary_Important_IssueLinked(t *testing.T) {
	got := MapCustomerNeedSummary(linear.CustomerNeedSummaryFields{
		Id:        "n1",
		Priority:  1,
		Body:      ptr.To("Wants SSO"),
		CreatedAt: "2026-01-01T00:00:00.000Z",
		Customer:  &linear.CustomerNeedSummaryFieldsCustomer{Id: "c1", Name: "Acme Corp"},
		Issue:     &linear.CustomerNeedSummaryFieldsIssue{Identifier: "ENG-123", Title: "Add SSO"},
	})

	if got["important"] != true {
		t.Errorf("important = %v, want true", got["important"])
	}
	issue, ok := got["issue"].(map[string]any)
	if !ok || issue["identifier"] != "ENG-123" {
		t.Errorf("issue = %v", got["issue"])
	}
	if _, ok := got["project"]; ok {
		t.Error("project should be absent for an issue-linked need")
	}
}

func TestMapCustomerNeedSummary_NotImportant(t *testing.T) {
	got := MapCustomerNeedSummary(linear.CustomerNeedSummaryFields{
		Id:       "n2",
		Priority: 0,
		Project:  &linear.CustomerNeedSummaryFieldsProject{Id: "p1", Name: "Q3 Roadmap"},
	})

	if got["important"] != false {
		t.Errorf("important = %v, want false", got["important"])
	}
	if _, ok := got["issue"]; ok {
		t.Error("issue should be absent for a project-linked need")
	}
}
