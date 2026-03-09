"""
File Operations Tests

Tests for file system operations including reading, writing, copying,
moving, deleting files, and directory management.
"""

import os
import json
import pytest
import shutil
import tempfile
from pathlib import Path
from datetime import datetime
from unittest.mock import Mock, patch, mock_open

from file_ops import (
    # File operations
    read_file,
    read_file_lines,
    write_file,
    append_file,
    copy_file,
    move_file,
    delete_file,
    
    # Directory operations
    list_files,
    list_dirs,
    ensure_directory,
    delete_directory,
    copy_directory,
    move_directory,
    
    # File info
    file_exists,
    is_file,
    is_dir,
    get_file_size,
    get_file_extension,
    get_file_name,
    get_file_path,
    get_modified_time,
    get_creation_time,
    
    # File searching
    find_files,
    find_by_extension,
    find_by_name,
    find_by_content,
    
    # File permissions
    is_readable,
    is_writable,
    is_executable,
    set_permissions,
    
    # File comparison
    files_are_equal,
    compare_by_size,
    compare_by_content,
    
    # File hashing
    get_md5,
    get_sha256,
    get_sha1,
    
    # File encoding
    detect_encoding,
    convert_encoding,
    
    # Temporary files
    temp_file,
    temp_directory,
    
    # File locking
    lock_file,
    unlock_file,
    
    # File watching
    watch_file,
    watch_directory,
    
    # Archives
    zip_files,
    unzip_file,
    tar_files,
    untar_file
)

# ============================================================================
# Basic File Operations Tests
# ============================================================================

class TestBasicFileOperations:
    """Tests for basic file operations."""

    def test_write_and_read_file(self, temp_dir):
        """Test writing to and reading from a file."""
        file_path = temp_dir / "test.txt"
        content = "Hello, World!\nThis is a test file."
        
        # Write file
        write_file(file_path, content)
        assert file_path.exists()
        
        # Read file
        read_content = read_file(file_path)
        assert read_content == content
        
        # Read lines
        lines = read_file_lines(file_path)
        assert len(lines) == 2
        assert lines[0] == "Hello, World!"
        assert lines[1] == "This is a test file."

    def test_append_to_file(self, temp_dir):
        """Test appending content to a file."""
        file_path = temp_dir / "append.txt"
        
        # Write initial content
        write_file(file_path, "Line 1\n")
        
        # Append content
        append_file(file_path, "Line 2\n")
        append_file(file_path, "Line 3\n")
        
        # Verify content
        content = read_file(file_path)
        assert content == "Line 1\nLine 2\nLine 3\n"

    def test_write_empty_file(self, temp_dir):
        """Test writing an empty file."""
        file_path = temp_dir / "empty.txt"
        
        write_file(file_path, "")
        assert file_path.exists()
        assert get_file_size(file_path) == 0
        
        content = read_file(file_path)
        assert content == ""

    def test_write_binary_file(self, temp_dir):
        """Test writing and reading binary files."""
        file_path = temp_dir / "binary.bin"
        binary_data = bytes([0x00, 0x01, 0x02, 0x03, 0xFF])
        
        write_file(file_path, binary_data, binary=True)
        read_data = read_file(file_path, binary=True)
        
        assert read_data == binary_data

    def test_write_with_encoding(self, temp_dir):
        """Test writing with different encodings."""
        file_path = temp_dir / "encoding.txt"
        content = "Café Français 中文 русский"
        
        # UTF-8
        write_file(file_path, content, encoding='utf-8')
        read_utf8 = read_file(file_path, encoding='utf-8')
        assert read_utf8 == content
        
        # UTF-16
        write_file(file_path, content, encoding='utf-16')
        read_utf16 = read_file(file_path, encoding='utf-16')
        assert read_utf16 == content

# ============================================================================
# File Copy/Move/Delete Tests
# ============================================================================

class TestFileCopyMoveDelete:
    """Tests for copying, moving, and deleting files."""

    def test_copy_file(self, temp_dir, sample_file):
        """Test copying a file."""
        dest_path = temp_dir / "copy.txt"
        
        copy_file(sample_file, dest_path)
        
        assert dest_path.exists()
        assert files_are_equal(sample_file, dest_path)

    def test_copy_file_to_directory(self, temp_dir, sample_file):
        """Test copying a file to a directory."""
        dest_dir = temp_dir / "subdir"
        ensure_directory(dest_dir)
        
        copy_file(sample_file, dest_dir)
        
        expected_path = dest_dir / sample_file.name
        assert expected_path.exists()
        assert files_are_equal(sample_file, expected_path)

    def test_copy_nonexistent_file(self, temp_dir):
        """Test copying a nonexistent file."""
        src = temp_dir / "nonexistent.txt"
        dst = temp_dir / "dest.txt"
        
        with pytest.raises(FileNotFoundError):
            copy_file(src, dst)

    def test_copy_file_overwrite(self, temp_dir, sample_file):
        """Test overwriting an existing file when copying."""
        dest_path = temp_dir / "dest.txt"
        write_file(dest_path, "original content")
        
        copy_file(sample_file, dest_path, overwrite=True)
        
        assert files_are_equal(sample_file, dest_path)

    def test_copy_file_no_overwrite(self, temp_dir, sample_file):
        """Test copying without overwrite."""
        dest_path = temp_dir / "dest.txt"
        original_content = "original content"
        write_file(dest_path, original_content)
        
        with pytest.raises(FileExistsError):
            copy_file(sample_file, dest_path, overwrite=False)
        
        # Verify original unchanged
        assert read_file(dest_path) == original_content

    def test_move_file(self, temp_dir, sample_file):
        """Test moving a file."""
        dest_path = temp_dir / "moved.txt"
        
        move_file(sample_file, dest_path)
        
        assert dest_path.exists()
        assert not sample_file.exists()

    def test_move_file_to_directory(self, temp_dir, sample_file):
        """Test moving a file to a directory."""
        dest_dir = temp_dir / "subdir"
        ensure_directory(dest_dir)
        
        move_file(sample_file, dest_dir)
        
        expected_path = dest_dir / sample_file.name
        assert expected_path.exists()
        assert not sample_file.exists()

    def test_delete_file(self, temp_dir, sample_file):
        """Test deleting a file."""
        assert sample_file.exists()
        
        delete_file(sample_file)
        
        assert not sample_file.exists()

    def test_delete_nonexistent_file(self, temp_dir):
        """Test deleting a nonexistent file."""
        file_path = temp_dir / "nonexistent.txt"
        
        # Should not raise error by default
        delete_file(file_path, missing_ok=True)
        
        with pytest.raises(FileNotFoundError):
            delete_file(file_path, missing_ok=False)

# ============================================================================
# Directory Operations Tests
# ============================================================================

class TestDirectoryOperations:
    """Tests for directory operations."""

    def test_ensure_directory(self, temp_dir):
        """Test creating directories."""
        dir_path = temp_dir / "a" / "b" / "c"
        
        ensure_directory(dir_path)
        
        assert dir_path.exists()
        assert dir_path.is_dir()

    def test_ensure_existing_directory(self, temp_dir):
        """Test ensuring an existing directory."""
        ensure_directory(temp_dir)  # Should not raise error

    def test_list_files(self, temp_dir, sample_files):
        """Test listing files in a directory."""
        files = list_files(temp_dir)
        
        assert len(files) == 3
        assert all(f.exists() for f in files)

    def test_list_files_with_pattern(self, temp_dir, sample_files):
        """Test listing files with a pattern."""
        txt_files = list_files(temp_dir, pattern="*.txt")
        assert len(txt_files) == 2
        
        py_files = list_files(temp_dir, pattern="*.py")
        assert len(py_files) == 1

    def test_list_files_recursive(self, temp_dir, nested_files):
        """Test listing files recursively."""
        all_files = list_files(temp_dir, recursive=True)
        assert len(all_files) == 4
        
        nested_files = list_files(temp_dir / "subdir", recursive=True)
        assert len(nested_files) == 2

    def test_list_dirs(self, temp_dir, nested_files):
        """Test listing directories."""
        dirs = list_dirs(temp_dir)
        
        assert len(dirs) == 1
        assert dirs[0].name == "subdir"

    def test_delete_directory(self, temp_dir, nested_files):
        """Test deleting a directory."""
        dir_path = temp_dir / "subdir"
        assert dir_path.exists()
        
        delete_directory(dir_path)
        
        assert not dir_path.exists()
        assert (temp_dir / "file1.txt").exists()  # Parent files preserved

    def test_delete_directory_recursive(self, temp_dir, nested_files):
        """Test recursive directory deletion."""
        delete_directory(temp_dir, recursive=True)
        
        assert not temp_dir.exists()

    def test_copy_directory(self, temp_dir, nested_files):
        """Test copying a directory."""
        dest_dir = temp_dir / "copy"
        
        copy_directory(temp_dir / "subdir", dest_dir)
        
        assert dest_dir.exists()
        assert (dest_dir / "nested1.txt").exists()
        assert (dest_dir / "nested2.txt").exists()

    def test_move_directory(self, temp_dir, nested_files):
        """Test moving a directory."""
        dest_dir = temp_dir / "moved"
        
        move_directory(temp_dir / "subdir", dest_dir)
        
        assert dest_dir.exists()
        assert not (temp_dir / "subdir").exists()
        assert (dest_dir / "nested1.txt").exists()

# ============================================================================
# File Information Tests
# ============================================================================

class TestFileInformation:
    """Tests for file information functions."""

    def test_file_exists(self, temp_dir, sample_file):
        """Test checking if file exists."""
        assert file_exists(sample_file)
        assert not file_exists(temp_dir / "nonexistent.txt")

    def test_is_file(self, temp_dir, sample_file):
        """Test checking if path is a file."""
        assert is_file(sample_file)
        assert not is_file(temp_dir)

    def test_is_dir(self, temp_dir, sample_file):
        """Test checking if path is a directory."""
        assert is_dir(temp_dir)
        assert not is_dir(sample_file)

    def test_get_file_size(self, temp_dir, sample_file):
        """Test getting file size."""
        size = get_file_size(sample_file)
        assert size > 0
        
        # Empty file
        empty_file = temp_dir / "empty.txt"
        write_file(empty_file, "")
        assert get_file_size(empty_file) == 0

    def test_get_file_extension(self):
        """Test getting file extension."""
        assert get_file_extension("file.txt") == ".txt"
        assert get_file_extension("file.tar.gz") == ".gz"
        assert get_file_extension("file") == ""
        assert get_file_extension(".gitignore") == ""

    def test_get_file_name(self):
        """Test getting file name from path."""
        assert get_file_name("/path/to/file.txt") == "file.txt"
        assert get_file_name("file.txt") == "file.txt"
        assert get_file_name("/path/to/") == ""

    def test_get_file_path(self):
        """Test getting directory path from file path."""
        assert get_file_path("/path/to/file.txt") == "/path/to"
        assert get_file_path("file.txt") == "."

    def test_get_modified_time(self, sample_file):
        """Test getting file modified time."""
        mtime = get_modified_time(sample_file)
        assert isinstance(mtime, datetime)

    def test_get_creation_time(self, sample_file):
        """Test getting file creation time."""
        ctime = get_creation_time(sample_file)
        assert isinstance(ctime, datetime)

# ============================================================================
# File Searching Tests
# ============================================================================

class TestFileSearching:
    """Tests for file searching functions."""

    def test_find_files(self, temp_dir, sample_files):
        """Test finding files."""
        results = find_files(temp_dir, "file*")
        assert len(results) == 3

    def test_find_by_extension(self, temp_dir, sample_files):
        """Test finding files by extension."""
        txt_files = find_by_extension(temp_dir, ".txt")
        assert len(txt_files) == 2
        
        py_files = find_by_extension(temp_dir, ".py")
        assert len(py_files) == 1

    def test_find_by_name(self, temp_dir, sample_files):
        """Test finding files by name."""
        results = find_by_name(temp_dir, "file1")
        assert len(results) == 1
        assert results[0].name == "file1.txt"

    def test_find_by_content(self, temp_dir, sample_files):
        """Test finding files by content."""
        # Create files with specific content
        file1 = temp_dir / "content1.txt"
        file2 = temp_dir / "content2.txt"
        write_file(file1, "Hello World")
        write_file(file2, "Goodbye World")
        
        results = find_by_content(temp_dir, "Hello")
        assert len(results) == 1
        assert results[0] == file1

    def test_find_by_content_recursive(self, temp_dir, nested_files):
        """Test finding files by content recursively."""
        # Add content to nested files
        nested1 = temp_dir / "subdir" / "nested1.txt"
        nested2 = temp_dir / "subdir" / "nested2.txt"
        write_file(nested1, "Secret content")
        write_file(nested2, "Public content")
        
        results = find_by_content(temp_dir, "Secret", recursive=True)
        assert len(results) == 1
        assert results[0] == nested1

# ============================================================================
# File Permissions Tests
# ============================================================================

class TestFilePermissions:
    """Tests for file permission functions."""

    def test_is_readable(self, sample_file):
        """Test checking if file is readable."""
        assert is_readable(sample_file)

    def test_is_writable(self, sample_file):
        """Test checking if file is writable."""
        assert is_writable(sample_file)

    def test_is_executable(self, sample_file):
        """Test checking if file is executable."""
        # Regular text files shouldn't be executable
        assert not is_executable(sample_file)

    @pytest.mark.skipif(os.name == 'nt', reason="Permission tests not reliable on Windows")
    def test_set_permissions(self, temp_dir, sample_file):
        """Test setting file permissions."""
        # Make read-only
        set_permissions(sample_file, 0o444)
        assert is_readable(sample_file)
        assert not is_writable(sample_file)
        
        # Make writable
        set_permissions(sample_file, 0o644)
        assert is_readable(sample_file)
        assert is_writable(sample_file)

# ============================================================================
# File Comparison Tests
# ============================================================================

class TestFileComparison:
    """Tests for file comparison functions."""

    def test_files_are_equal_same(self, temp_dir):
        """Test comparing identical files."""
        file1 = temp_dir / "file1.txt"
        file2 = temp_dir / "file2.txt"
        content = "Same content"
        
        write_file(file1, content)
        write_file(file2, content)
        
        assert files_are_equal(file1, file2)

    def test_files_are_equal_different(self, temp_dir):
        """Test comparing different files."""
        file1 = temp_dir / "file1.txt"
        file2 = temp_dir / "file2.txt"
        
        write_file(file1, "Content A")
        write_file(file2, "Content B")
        
        assert not files_are_equal(file1, file2)

    def test_files_are_equal_different_sizes(self, temp_dir):
        """Test comparing files of different sizes."""
        file1 = temp_dir / "file1.txt"
        file2 = temp_dir / "file2.txt"
        
        write_file(file1, "Short")
        write_file(file2, "Much longer content")
        
        assert not files_are_equal(file1, file2)

    def test_compare_by_size(self, temp_dir):
        """Test comparing files by size."""
        file1 = temp_dir / "file1.txt"
        file2 = temp_dir / "file2.txt"
        
        write_file(file1, "Same length")
        write_file(file2, "Same length")
        
        assert compare_by_size(file1, file2)

    def test_compare_by_content(self, temp_dir):
        """Test comparing files by content."""
        file1 = temp_dir / "file1.txt"
        file2 = temp_dir / "file2.txt"
        
        write_file(file1, "Content")
        write_file(file2, "Content")
        
        assert compare_by_content(file1, file2)

# ============================================================================
# File Hashing Tests
# ============================================================================

class TestFileHashing:
    """Tests for file hashing functions."""

    def test_get_md5(self, temp_dir, sample_file):
        """Test calculating MD5 hash."""
        hash1 = get_md5(sample_file)
        hash2 = get_md5(sample_file)
        
        assert hash1 == hash2
        assert len(hash1) == 32  # MD5 is 32 hex chars

    def test_get_sha256(self, temp_dir, sample_file):
        """Test calculating SHA256 hash."""
        hash1 = get_sha256(sample_file)
        hash2 = get_sha256(sample_file)
        
        assert hash1 == hash2
        assert len(hash1) == 64  # SHA256 is 64 hex chars

    def test_get_sha1(self, temp_dir, sample_file):
        """Test calculating SHA1 hash."""
        hash1 = get_sha1(sample_file)
        hash2 = get_sha1(sample_file)
        
        assert hash1 == hash2
        assert len(hash1) == 40  # SHA1 is 40 hex chars

    def test_different_files_different_hashes(self, temp_dir):
        """Test that different files have different hashes."""
        file1 = temp_dir / "file1.txt"
        file2 = temp_dir / "file2.txt"
        
        write_file(file1, "Content 1")
        write_file(file2, "Content 2")
        
        assert get_md5(file1) != get_md5(file2)
        assert get_sha256(file1) != get_sha256(file2)

# ============================================================================
# File Encoding Tests
# ============================================================================

class TestFileEncoding:
    """Tests for file encoding detection and conversion."""

    def test_detect_encoding_utf8(self, temp_dir):
        """Test detecting UTF-8 encoding."""
        file_path = temp_dir / "utf8.txt"
        write_file(file_path, "Hello World", encoding='utf-8')
        
        encoding = detect_encoding(file_path)
        assert encoding.lower() in ['utf-8', 'ascii']

    def test_detect_encoding_utf16(self, temp_dir):
        """Test detecting UTF-16 encoding."""
        file_path = temp_dir / "utf16.txt"
        write_file(file_path, "Hello World", encoding='utf-16')
        
        encoding = detect_encoding(file_path)
        assert 'utf-16' in encoding.lower()

    def test_convert_encoding(self, temp_dir):
        """Test converting file encoding."""
        file_path = temp_dir / "convert.txt"
        original_content = "Café Français"
        
        # Write in UTF-8
        write_file(file_path, original_content, encoding='utf-8')
        
        # Convert to UTF-16
        convert_encoding(file_path, 'utf-16', 'utf-8')
        
        # Read back in UTF-16
        content = read_file(file_path, encoding='utf-16')
        assert content == original_content

# ============================================================================
# Temporary File Tests
# ============================================================================

class TestTemporaryFiles:
    """Tests for temporary file creation."""

    def test_temp_file_context_manager(self):
        """Test temporary file context manager."""
        with temp_file() as tmp:
            assert tmp.exists()
            write_file(tmp, "test content")
            assert read_file(tmp) == "test content"
        
        # File should be deleted after context
        assert not tmp.exists()

    def test_temp_file_with_content(self):
        """Test temporary file with initial content."""
        with temp_file(content="initial content") as tmp:
            assert read_file(tmp) == "initial content"

    def test_temp_file_with_suffix(self):
        """Test temporary file with custom suffix."""
        with temp_file(suffix=".txt") as tmp:
            assert tmp.suffix == ".txt"

    def test_temp_directory_context_manager(self):
        """Test temporary directory context manager."""
        with temp_directory() as tmp_dir:
            assert tmp_dir.exists()
            assert tmp_dir.is_dir()
            
            # Create file in temp directory
            test_file = tmp_dir / "test.txt"
            write_file(test_file, "test")
            assert test_file.exists()
        
        # Directory should be deleted after context
        assert not tmp_dir.exists()

# ============================================================================
# File Locking Tests
# ============================================================================

class TestFileLocking:
    """Tests for file locking mechanisms."""

    def test_lock_file(self, sample_file):
        """Test locking a file."""
        lock = lock_file(sample_file)
        assert lock is not None
        
        # Try to acquire lock again (should fail)
        with pytest.raises(Exception):
            lock_file(sample_file, timeout=0.1)
        
        unlock_file(lock)

    def test_unlock_file(self, sample_file):
        """Test unlocking a file."""
        lock = lock_file(sample_file)
        unlock_file(lock)
        
        # Should be able to lock again
        lock2 = lock_file(sample_file)
        assert lock2 is not None
        unlock_file(lock2)

# ============================================================================
# File Watching Tests
# ============================================================================

class TestFileWatching:
    """Tests for file watching functionality."""

    def test_watch_file_creation(self, temp_dir):
        """Test watching file creation."""
        events = []
        
        def callback(event):
            events.append(event)
        
        file_path = temp_dir / "watched.txt"
        
        # Start watching
        watcher = watch_file(file_path, callback)
        
        # Create file
        write_file(file_path, "content")
        
        # Give time for event to trigger
        import time
        time.sleep(0.1)
        
        watcher.stop()
        
        assert len(events) > 0
        assert events[0]['type'] == 'created'

    def test_watch_file_modification(self, sample_file):
        """Test watching file modification."""
        events = []
        
        def callback(event):
            events.append(event)
        
        # Start watching
        watcher = watch_file(sample_file, callback)
        
        # Modify file
        append_file(sample_file, "new content")
        
        # Give time for event to trigger
        import time
        time.sleep(0.1)
        
        watcher.stop()
        
        assert len(events) > 0
        assert events[0]['type'] == 'modified'

# ============================================================================
# Archive Tests
# ============================================================================

class TestArchives:
    """Tests for archive operations."""

    def test_zip_files(self, temp_dir, sample_files):
        """Test creating a ZIP archive."""
        zip_path = temp_dir / "archive.zip"
        
        zip_files(zip_path, sample_files)
        
        assert zip_path.exists()
        assert get_file_size(zip_path) > 0

    def test_unzip_file(self, temp_dir, sample_files):
        """Test extracting a ZIP archive."""
        zip_path = temp_dir / "archive.zip"
        extract_dir = temp_dir / "extracted"
        
        # Create zip
        zip_files(zip_path, sample_files)
        
        # Extract
        unzip_file(zip_path, extract_dir)
        
        assert extract_dir.exists()
        assert (extract_dir / "file1.txt").exists()
        assert (extract_dir / "file2.txt").exists()
        assert (extract_dir / "script.py").exists()

    def test_tar_files(self, temp_dir, sample_files):
        """Test creating a TAR archive."""
        tar_path = temp_dir / "archive.tar.gz"
        
        tar_files(tar_path, sample_files, compression='gz')
        
        assert tar_path.exists()
        assert get_file_size(tar_path) > 0

    def test_untar_file(self, temp_dir, sample_files):
        """Test extracting a TAR archive."""
        tar_path = temp_dir / "archive.tar.gz"
        extract_dir = temp_dir / "extracted"
        
        # Create tar
        tar_files(tar_path, sample_files, compression='gz')
        
        # Extract
        untar_file(tar_path, extract_dir)
        
        assert extract_dir.exists()
        assert (extract_dir / "file1.txt").exists()
        assert (extract_dir / "file2.txt").exists()
        assert (extract_dir / "script.py").exists()

# ============================================================================
# Edge Cases and Error Handling Tests
# ============================================================================

class TestEdgeCases:
    """Tests for edge cases and error handling."""

    def test_read_nonexistent_file(self, temp_dir):
        """Test reading a nonexistent file."""
        file_path = temp_dir / "nonexistent.txt"
        
        with pytest.raises(FileNotFoundError):
            read_file(file_path)

    def test_read_directory_as_file(self, temp_dir):
        """Test reading a directory as if it were a file."""
        with pytest.raises(IsADirectoryError):
            read_file(temp_dir)

    def test_write_to_nonexistent_directory(self, temp_dir):
        """Test writing to a file in a nonexistent directory."""
        file_path = temp_dir / "subdir" / "file.txt"
        
        # Should create directory automatically
        write_file(file_path, "content")
        assert file_path.exists()

    def test_very_long_filename(self, temp_dir):
        """Test handling very long filenames."""
        long_name = "a" * 255 + ".txt"
        file_path = temp_dir / long_name
        
        write_file(file_path, "content")
        assert file_path.exists()

    def test_file_with_special_characters(self, temp_dir):
        """Test handling files with special characters in names."""
        special_name = "!@#$%^&*()_+{}[]|\\;:'\",.<>?`~.txt"
        file_path = temp_dir / special_name
        
        write_file(file_path, "content")
        assert file_path.exists()

    def test_very_large_file(self, temp_dir):
        """Test handling very large files (100MB)."""
        file_path = temp_dir / "large.bin"
        
        # Create 100MB file
        chunk = b'x' * 1024 * 1024  # 1MB
        with open(file_path, 'wb') as f:
            for _ in range(100):
                f.write(chunk)
        
        size = get_file_size(file_path)
        assert size == 100 * 1024 * 1024

    def test_symlink_handling(self, temp_dir, sample_file):
        """Test handling symbolic links."""
        link_path = temp_dir / "link.txt"
        
        try:
            os.symlink(sample_file, link_path)
            
            # Should follow symlink
            assert is_file(link_path)
            assert files_are_equal(link_path, sample_file)
            
            # Delete symlink
            delete_file(link_path)
            assert not link_path.exists()
            assert sample_file.exists()  # Original should remain
            
        except OSError:
            pytest.skip("Symlinks not supported on this platform")

# ============================================================================
# Performance Tests
# ============================================================================

class TestPerformance:
    """Performance tests for file operations."""

    @pytest.mark.benchmark
    def test_read_file_performance(self, benchmark, large_file):
        """Benchmark file reading performance."""
        result = benchmark(read_file, large_file)
        assert len(result) > 0

    @pytest.mark.benchmark
    def test_write_file_performance(self, benchmark, temp_dir):
        """Benchmark file writing performance."""
        file_path = temp_dir / "bench.txt"
        content = "x" * 1024 * 1024  # 1MB
        
        benchmark(write_file, file_path, content)

    @pytest.mark.benchmark
    def test_copy_file_performance(self, benchmark, large_file, temp_dir):
        """Benchmark file copying performance."""
        dest = temp_dir / "copy.txt"
        
        benchmark(copy_file, large_file, dest)

    @pytest.mark.benchmark
    def test_hash_file_performance(self, benchmark, large_file):
        """Benchmark file hashing performance."""
        benchmark(get_sha256, large_file)

# ============================================================================
# Concurrency Tests
# ============================================================================

class TestConcurrency:
    """Tests for concurrent file operations."""

    def test_concurrent_reads(self, temp_dir, sample_file):
        """Test multiple threads reading the same file."""
        import threading
        
        results = []
        
        def read_file_thread():
            content = read_file(sample_file)
            results.append(content)
        
        threads = [threading.Thread(target=read_file_thread) for _ in range(10)]
        
        for t in threads:
            t.start()
        for t in threads:
            t.join()
        
        assert len(results) == 10
        assert all(r == results[0] for r in results)

    def test_concurrent_writes(self, temp_dir):
        """Test multiple threads writing to different files."""
        import threading
        
        file_paths = []
        
        def write_file_thread(i):
            path = temp_dir / f"thread_{i}.txt"
            file_paths.append(path)
            write_file(path, f"Content from thread {i}")
        
        threads = [threading.Thread(target=write_file_thread, args=(i,)) for i in range(10)]
        
        for t in threads:
            t.start()
        for t in threads:
            t.join()
        
        assert len(file_paths) == 10
        for i, path in enumerate(file_paths):
            assert path.exists()
            assert read_file(path) == f"Content from thread {i}"

# ============================================================================
# Fixtures
# ============================================================================

@pytest.fixture
def temp_dir():
    """Create a temporary directory for tests."""
    dirpath = tempfile.mkdtemp()
    yield Path(dirpath)
    shutil.rmtree(dirpath)

@pytest.fixture
def sample_file(temp_dir):
    """Create a sample file for testing."""
    file_path = temp_dir / "sample.txt"
    content = """This is a sample file.
It has multiple lines.
For testing purposes."""
    write_file(file_path, content)
    return file_path

@pytest.fixture
def sample_files(temp_dir):
    """Create multiple sample files."""
    files = []
    for name in ["file1.txt", "file2.txt", "script.py"]:
        path = temp_dir / name
        write_file(path, f"Content of {name}")
        files.append(path)
    return files

@pytest.fixture
def nested_files(temp_dir, sample_files):
    """Create nested directory structure with files."""
    subdir = temp_dir / "subdir"
    ensure_directory(subdir)
    
    nested1 = subdir / "nested1.txt"
    nested2 = subdir / "nested2.txt"
    write_file(nested1, "Nested file 1")
    write_file(nested2, "Nested file 2")
    
    return {
        'root': sample_files,
        'subdir': subdir,
        'nested': [nested1, nested2]
    }

@pytest.fixture
def large_file(temp_dir):
    """Create a large file for performance testing."""
    file_path = temp_dir / "large.bin"
    chunk = b'x' * 1024 * 1024  # 1MB
    with open(file_path, 'wb') as f:
        for _ in range(10):  # 10MB file
            f.write(chunk)
    return file_path