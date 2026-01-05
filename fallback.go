package patch2pr

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
)

func gitApply(data []byte, f *gitdiff.File) (out string, err error) {
	var b bytes.Buffer
	var cwd, dir string

	if cwd, err = os.Getwd(); err != nil {
		return
	}

	if dir, err = os.MkdirTemp("", "patch2pr-apply-*"); err != nil {
		return
	}

	os.Chdir(dir)
	defer os.Chdir(cwd)

	patchFile := filepath.Join(dir, "patch.patch")
	patchContents := []byte(f.String())
	if err = os.WriteFile(patchFile, patchContents, os.ModePerm); err != nil {
		return
	}

	cleanup := func() {
		if err == nil {
			os.RemoveAll(dir)
			return
		}
		os.Stderr.Write(patchContents)
		os.Stderr.WriteString("\n")
		if val, ok := err.(*exec.ExitError); ok {
			os.Stderr.Write(val.Stderr)
		}
		os.Stderr.WriteString("\n")
		os.Stderr.WriteString(dir + "\n")
	}
	defer cleanup()

	file := filepath.Join(dir, f.OldName)
	if err = os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return
	}
	if err = os.WriteFile(file, data, f.OldMode); err != nil {
		return
	}

	// only single file needs patching (parial patch)
	// TODO it can occur multiple times in the patch, see if we handle that case correctly
	// TODO --3way patch can't be done outside a git repo
	cmd := exec.Command("git", "apply", patchFile, "--include="+f.OldName)
	if _, err = cmd.Output(); err != nil {
		return
	}

	newFile, err := os.Open(filepath.Join(dir, f.NewName))
	if err != nil {
		return
	}

	w := bufio.NewWriter(&b)
	if _, err = io.Copy(w, newFile); err != nil {
		return
	}

	out = b.String()
	return
}

func base64GitApply(data []byte, f *gitdiff.File) (out string, err error) {
	var in string
	if in, err = gitApply(data, f); err != nil {
		return
	}
	out = base64.StdEncoding.EncodeToString([]byte(in))
	return
}
