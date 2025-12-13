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
