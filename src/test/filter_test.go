package main

import (
	"encoding/json"
	"testing"
)

func TestFilterFilesByName(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		filter   string
		expected string
	}{
		{
			name:     "Matching file",
			jsonData: `{"mount_path":"/mnt/storage","files":[{"name":"report.pdf","type":"file","is_symlink":false,"size_bytes":1024,"permissions":"rw-r--r--","hard_links":1,"inode":12345,"owner_uid":1000,"owner_gid":1000,"block_size":4096,"blocks_allocated":1,"last_modified":"2025-03-12T00:00:00Z","last_access":"2025-03-12T00:00:00Z"},{"name":"notes.txt","type":"file","is_symlink":false,"size_bytes":2048,"permissions":"rw-r--r--","hard_links":1,"inode":12346,"owner_uid":1000,"owner_gid":1000,"block_size":4096,"blocks_allocated":1,"last_modified":"2025-03-12T00:00:00Z","last_access":"2025-03-12T00:00:00Z"},{"name":"project.docx","type":"file","is_symlink":false,"size_bytes":4096,"permissions":"rw-r--r--","hard_links":1,"inode":12347,"owner_uid":1000,"owner_gid":1000,"block_size":4096,"blocks_allocated":1,"last_modified":"2025-03-12T00:00:00Z","last_access":"2025-03-12T00:00:00Z"}]}`,
			filter:   "note",
			expected: `{"mount_path":"/mnt/storage","files":[{"name":"notes.txt","type":"file","is_symlink":false,"size_bytes":2048,"permissions":"rw-r--r--","hard_links":1,"inode":12346,"owner_uid":1000,"owner_gid":1000,"block_size":4096,"blocks_allocated":1,"last_modified":"2025-03-12T00:00:00Z","last_access":"2025-03-12T00:00:00Z"}]}`,
		},
		{
			name:     "No matching file",
			jsonData: `{"mount_path":"/mnt/storage","files":[{"name":"report.pdf","type":"file","is_symlink":false,"size_bytes":1024,"permissions":"rw-r--r--","hard_links":1,"inode":12345,"owner_uid":1000,"owner_gid":1000,"block_size":4096,"blocks_allocated":1,"last_modified":"2025-03-12T00:00:00Z","last_access":"2025-03-12T00:00:00Z"},{"name":"project.docx","type":"file","is_symlink":false,"size_bytes":4096,"permissions":"rw-r--r--","hard_links":1,"inode":12347,"owner_uid":1000,"owner_gid":1000,"block_size":4096,"blocks_allocated":1,"last_modified":"2025-03-12T00:00:00Z","last_access":"2025-03-12T00:00:00Z"}]}`,
			filter:   "note",
			expected: `{"mount_path":"/mnt/storage","files":[]}`,
		},
		{
			name:     "Case insensitive match",
			jsonData: `{"mount_path":"/mnt/storage","files":[{"name":"Report.pdf","type":"file","is_symlink":false,"size_bytes":1024,"permissions":"rw-r--r--","hard_links":1,"inode":12345,"owner_uid":1000,"owner_gid":1000,"block_size":4096,"blocks_allocated":1,"last_modified":"2025-03-12T00:00:00Z","last_access":"2025-03-12T00:00:00Z"},{"name":"NOTES.txt","type":"file","is_symlink":false,"size_bytes":2048,"permissions":"rw-r--r--","hard_links":1,"inode":12346,"owner_uid":1000,"owner_gid":1000,"block_size":4096,"blocks_allocated":1,"last_modified":"2025-03-12T00:00:00Z","last_access":"2025-03-12T00:00:00Z"},{"name":"Project.docx","type":"file","is_symlink":false,"size_bytes":4096,"permissions":"rw-r--r--","hard_links":1,"inode":12347,"owner_uid":1000,"owner_gid":1000,"block_size":4096,"blocks_allocated":1,"last_modified":"2025-03-12T00:00:00Z","last_access":"2025-03-12T00:00:00Z"}]}`,
			filter:   "note",
			expected: `{"mount_path":"/mnt/storage","files":[{"name":"NOTES.txt","type":"file","is_symlink":false,"size_bytes":2048,"permissions":"rw-r--r--","hard_links":1,"inode":12346,"owner_uid":1000,"owner_gid":1000,"block_size":4096,"blocks_allocated":1,"last_modified":"2025-03-12T00:00:00Z","last_access":"2025-03-12T00:00:00Z"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filteredJSON, err := FilterFilesByName([]byte(tt.jsonData), tt.filter)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var expected, actual FilesResponse
			json.Unmarshal([]byte(tt.expected), &expected)
			json.Unmarshal(filteredJSON, &actual)

			if len(expected.Files) != len(actual.Files) {
				t.Errorf("expected %d files, got %d", len(expected.Files), len(actual.Files))
			}

			for i := range expected.Files {
				if expected.Files[i].Name != actual.Files[i].Name {
					t.Errorf("expected file name %s, got %s", expected.Files[i].Name, actual.Files[i].Name)
				}
			}
		})
	}
}
