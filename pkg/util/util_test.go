/*
Copyright © 2020 Corey Daley <cdaley@redhat.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package util

import (
	"fmt"
	"os"
	"testing"
)

func TestIsGoFile(t *testing.T) {
	tests := []struct {
		filename string
		isGoFile bool
	}{
		{
			"test.go",
			true,
		},
		{
			"test.txt",
			false,
		},
	}

	for _, test := range tests {
		f, err := os.Create(fmt.Sprintf("%s/%s", os.TempDir(), test.filename))
		if err != nil {
			t.Fatalf("TempFile %s: %s", test.filename, err)
		}

		fi, err := f.Stat()
		if err != nil {
			t.Fatalf("TempFile %s: %s", test.filename, err)
		}

		if IsGoFile(fi) != test.isGoFile {
			t.Fatalf("TempFile %s: wanted %t, got %t, details: %#v", test.filename, test.isGoFile, IsGoFile(fi), fi)
		}
	}
}

func TestExcludeDirsRegExp(t *testing.T) {
	tests := []struct {
		name        string
		excludeDirs []string
		have        string
		want        bool
	}{
		{
			name:        "basic test the proper regexp (dir: '/home/go/src/project/test/pkg')",
			excludeDirs: []string{"test"},
			have:        "/home/go/src/project/test/pkg",
			want:        true,
		},
		{
			name:        "basic test the proper regexp (dir: '/home/go/src/project/vendor/pkg')",
			excludeDirs: []string{"test"},
			have:        "/home/go/src/project/vendor/pkg",
			want:        false,
		},
		{
			name:        "basic test the proper regexp for dir with subdir (dir: '/home/go/src/project/test/generated')",
			excludeDirs: []string{"test/generated"},
			have:        "/home/go/src/project/test/generated",
			want:        true,
		},
		{
			name:        "basic test the proper regexp for dir with subdir (dir: '/home/go/src/project/test/pkg')",
			excludeDirs: []string{"test/generated"},
			have:        "/home/go/src/project/test/pkg",
			want:        false,
		},
		{
			name:        "basic test the proper regexp for dir with subdir (dir: '/home/go/src/project/vendor/pkg')",
			excludeDirs: []string{"test/generated"},
			have:        "/home/go/src/project/vendor/pkg",
			want:        false,
		},
		{
			name:        "basic test the proper regexp for dir is nil",
			excludeDirs: nil,
			have:        "/home/go/src/project/test/pkg",
			want:        false,
		},
		{
			name:        "basic test the proper regexp for dir is empty",
			excludeDirs: []string{},
			have:        "/home/go/src/project/test/pkg",
			want:        false,
		},
	}

	for _, test := range tests {
		r := ExcludedDirsRegExp(test.excludeDirs)

		if len(r) != len(test.excludeDirs) {
			t.Fatalf("ExcludeDirs %s, wanted: %d, got %d", test.name, len(test.excludeDirs), len(r))
		}

		if len(r) == 0 {
			continue
		}

		if match := r[0].MatchString(test.have); match != test.want {
			t.Fatalf("ExcludeDirs %s, wanted: %t, got %t", test.name, test.want, match)
		}
	}
}
