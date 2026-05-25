package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

const testBin = "./clock-gate-test"
const testGameFile = "../../testdata/2011_01_NO_GB.json"
const testDataDir = "../../testdata/"

func TestMain(m *testing.M) {
	out, err := exec.Command("go", "build", "-o", testBin, ".").CombinedOutput()
	if err != nil {
		os.Stderr.Write(out)
		os.Exit(1)
	}
	code := m.Run()
	os.Remove(testBin)
	os.Exit(code)
}

func TestSmokeTickZero(t *testing.T) {
	cmd := exec.Command(testBin, "--tick", "0", testGameFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got: %v\noutput: %s", err, out)
	}
	s := string(out)
	if !strings.Contains(s, "Q1") && !strings.Contains(s, "Elapsed: 0s") {
		t.Errorf("expected output to contain Q1 or Elapsed: 0s, got: %s", s)
	}
}

func TestSmokeTickMidgame(t *testing.T) {
	cmd := exec.Command(testBin, "--tick", "1800", testGameFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got: %v\noutput: %s", err, out)
	}
	if len(strings.TrimSpace(string(out))) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestSmokeTickOverflow(t *testing.T) {
	cmd := exec.Command(testBin, "--tick", "999999", testGameFile)
	out, _ := cmd.CombinedOutput()
	if cmd.ProcessState.ExitCode() == 0 {
		t.Fatalf("expected non-zero exit, got exit 0\noutput: %s", out)
	}
	if !strings.Contains(string(out), "exceeds") {
		t.Errorf("expected output to contain 'exceeds', got: %s", out)
	}
}

func TestSmokeFormatJSON(t *testing.T) {
	cmd := exec.Command(testBin, "--tick", "900", "--format", "json", testGameFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got: %v\noutput: %s", err, out)
	}
	s := strings.TrimSpace(string(out))
	if !strings.HasPrefix(s, "{") {
		t.Errorf("expected JSON output starting with '{', got: %s", s)
	}
}

func TestSmokeList(t *testing.T) {
	cmd := exec.Command(testBin, "--list", testDataDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got: %v\noutput: %s", err, out)
	}
	if !strings.Contains(string(out), ".json") {
		t.Errorf("expected output to contain .json filenames, got: %s", out)
	}
}
