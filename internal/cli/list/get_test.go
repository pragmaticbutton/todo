package list_test

import (
	"bytes"
	"path/filepath"
	"strconv"
	"testing"

	listcmd "github.com/pragmaticbutton/todo/internal/cli/list"
	"github.com/pragmaticbutton/todo/internal/domain/list"
)

func TestGetCmd_Golden(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		list   list.List
		golden string
	}{
		{
			name: "basic",
			list: list.List{
				ID:          1,
				Name:        "Groceries",
				Description: "Weekly shopping",
			},
			golden: filepath.Join("testdata", "get", "basic.golden"),
		},
		{
			name: "no description",
			list: list.List{
				ID:   2,
				Name: "Chores",
			},
			golden: filepath.Join("testdata", "get", "no_description.golden"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := newListServiceWithLists(t, tc.list)
			cmd := listcmd.NewGetCmd(svc)

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetArgs([]string{strconv.Itoa(int(tc.list.ID))})

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute get command: %v", err)
			}

			assertGolden(t, buf.String(), tc.golden)
		})
	}
}

func TestGetCmd_Integration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		lists      []list.List
		args       []string
		expectedID uint32
		checkFn    func(t *testing.T, lst *list.List)
	}{
		{
			name:       "gets existing list",
			lists:      []list.List{{ID: 1, Name: "Groceries", Description: "Weekly"}},
			args:       []string{"1"},
			expectedID: 1,
			checkFn: func(t *testing.T, lst *list.List) {
				if lst.Name != "Groceries" || lst.Description != "Weekly" {
					t.Fatalf("unexpected list %+v", lst)
				}
			},
		},
		{
			name: "gets list when multiple exist",
			lists: []list.List{
				{ID: 1, Name: "Groceries"},
				{ID: 2, Name: "Work", Description: "Tasks"},
			},
			args:       []string{"2"},
			expectedID: 2,
			checkFn: func(t *testing.T, lst *list.List) {
				if lst.Name != "Work" {
					t.Fatalf("expected Work list, got %s", lst.Name)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newListServiceWithLists(t, tt.lists...)
			cmd := listcmd.NewGetCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute get command: %v", err)
			}

			lst, err := svc.GetList(tt.expectedID)
			if err != nil {
				t.Fatalf("get list: %v", err)
			}

			tt.checkFn(t, lst)
		})
	}
}

func TestGetCmd_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		lists []list.List
		args  []string
	}{
		{
			name: "invalid id",
			args: []string{"abc"},
		},
		{
			name:  "missing list",
			lists: []list.List{{ID: 1, Name: "Groceries"}},
			args:  []string{"2"},
		},
		{
			name: "missing argument",
			args: []string{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newListServiceWithLists(t, tt.lists...)
			cmd := listcmd.NewGetCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err == nil {
				t.Fatalf("expected error for case %q", tt.name)
			}
		})
	}
}
