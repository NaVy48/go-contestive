package judge

import (
	"contestive/judge/unzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
)

const testDirectory = "/tmp/judge"

func checkDir(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func unpackProblem(archive []byte, dest string) (err error) {
	defer func() {
		if err != nil {
			os.RemoveAll(dest)
		}
	}()

	err = unzip.Unzip(archive, dest)
	if err != nil {
		return err
	}

	err = CppCompile(path.Join(dest, "files", "check.cpp"), path.Join(dest, "files", "check"))
	if err != nil {
		return err
	}

	err = RunShell(dest, "doall.sh")
	if err != nil {
		return err
	}

	return nil
}

func copyFile(src, dest string, mod fs.FileMode) error {
	fmt.Printf("Copying from '%s'\n", src)
	fmt.Printf("Copying to '%s'\n", dest)

	destFile, err := os.Create(dest)
	destFile.Chmod(mod)
	if err != nil {
		return err
	}
	defer destFile.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func RunShell(dir, script string) error {
	fmt.Printf("running %s/%s\n", dir, script)
	cmd := exec.Command("/bin/sh", script)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func CppCompile(src, dest string) error {
	cmd := exec.Command("g++", "-static", "-DONLINE_JUDGE", "-lm", "-s", "-x", "c++", "-O2", "-std=c++11", "-o", dest, src)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
