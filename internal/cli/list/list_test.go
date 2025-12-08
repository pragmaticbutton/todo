package list_test

import (
	"bytes"
	"path/filepath"
	"testing"

	listcmd "github.com/pragmaticbutton/todo/internal/cli/list"
	"github.com/pragmaticbutton/todo/internal/domain/list"
)

func TestListCmd_Golden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		lists  []list.List
		golden string
	}{
		{
			name:   "empty",
			lists:  []list.List{},
			golden: filepath.Join("testdata", "list", "empty.golden"),
		},
		{
			name: "mixed",
			lists: []list.List{
				{ID: 2, Name: "Work", Description: "Tasks"},
				{ID: 1, Name: "Groceries", Description: "Weekly"},
			},
			golden: filepath.Join("testdata", "list", "mixed.golden"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newListServiceWithLists(t, tt.lists...)
			cmd := listcmd.NewListListsCmd(svc)

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetArgs([]string{})

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute list: %v", err)
			}

			assertGolden(t, buf.String(), tt.golden)
		})
	}
}
