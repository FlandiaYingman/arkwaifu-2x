package verify

import (
	"errors"
	"fmt"
	"hash/adler32"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

// Verify verifies whether the given file is ready-to-use or not.
// If the file and its checksum exist, and they are the same, it returns true.
// If either the file or its checksum does not exist, or the checksums are not the same, it returns false.
// Otherwise, it returns an error.
//
// Verify returns false if any error is returned. So the caller may ignore the error and just check the return value.
func Verify(filePath string) (bool, error) {
	actualSum, err := calcSum(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	expectedSum, err := readSum(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return actualSum == expectedSum, nil
}

// Done calculates the checksum of the given file, and writes it to the file's checksum file.
// The file's checksum file is the filePath with ".sum" appended.
func Done(filePath string) error {
	sum, err := calcSum(filePath)
	if err != nil {
		return err
	}

	sumFilePath := fmt.Sprintf("%s.sum", filePath)
	err = os.WriteFile(sumFilePath, []byte(fmt.Sprintf("%d", sum)), 0644)
	if err != nil {
		return err
	}

	return nil
}

func calcSum(filePath string) (uint32, error) {
	h := adler32.New()
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer func() { _ = f.Close() }()

	_, err = io.Copy(h, f)
	if err != nil {
		return 0, err
	}

	return h.Sum32(), nil
}
func readSum(filePath string) (uint32, error) {
	sumFilePath := fmt.Sprintf("%s.sum", filePath)
	sum, err := ioutil.ReadFile(sumFilePath)
	if err != nil {
		return 0, err
	}
	sum64, err := strconv.ParseUint(string(sum), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(sum64), nil
}
