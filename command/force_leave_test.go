package command

import (
	"github.com/hashicorp/serf/cli"
	"github.com/hashicorp/serf/testutil"
	"strings"
	"testing"
)

func TestForceLeaveCommand_implements(t *testing.T) {
	var _ cli.Command = &ForceLeaveCommand{}
}

func TestForceLeaveCommandRun(t *testing.T) {
	a1 := testAgent(t)
	a2 := testAgent(t)
	defer a1.Shutdown()
	defer a2.Shutdown()

	_, err := a1.Join([]string{a2.SerfConfig.MemberlistConfig.BindAddr})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testutil.Yield()

	// Forcibly shutdown a2 so that it appears "failed" in a1
	if err := a2.Shutdown(); err != nil {
		t.Fatalf("err: %s", err)
	}

	testutil.Yield()

	c := &ForceLeaveCommand{}
	args := []string{
		"-rpc-addr=" + a1.RPCAddr,
		a2.SerfConfig.NodeName,
	}
	ui := new(cli.MockUi)

	code := c.Run(args, ui)
	if code != 0 {
		t.Fatalf("bad: %d. %#v", code, ui.ErrorWriter.String())
	}

	if len(a1.Serf().Members()) != 1 {
		t.Fatalf("bad: %#v", a1.Serf().Members())
	}
}

func TestForceLeaveCommandRun_noAddrs(t *testing.T) {
	c := &ForceLeaveCommand{}
	args := []string{"-rpc-addr=foo"}
	ui := new(cli.MockUi)

	code := c.Run(args, ui)
	if code != 1 {
		t.Fatalf("bad: %d", code)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "node name") {
		t.Fatalf("bad: %#v", ui.ErrorWriter.String())
	}
}
