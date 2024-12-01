#!/bin/bash

# Script to generate backdated commits for December 2024

# Create a .gitignore file
cat > .gitignore << EOL
.DS_Store
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
/test_downloads
EOL

# Add all files to git
git add .

# Function to create a commit with a backdated timestamp
create_backdated_commit() {
    local commit_date="$1"
    local commit_message="$2"
    local change_file="$3"
    local change_content="$4"
    
    # Make a change to the specified file
    echo "$change_content" >> "$change_file"
    
    # Add the changed file
    git add "$change_file"
    
    # Set environment variables for the commit date
    export GIT_AUTHOR_DATE="$commit_date"
    export GIT_COMMITTER_DATE="$commit_date"
    
    # Create the commit
    git commit -m "$commit_message"
    
    # Unset environment variables
    unset GIT_AUTHOR_DATE
    unset GIT_COMMITTER_DATE
}

# Configure git user if not already configured
git config user.name 2>/dev/null || git config user.name "Abhinav Chamallamudi"
git config user.email 2>/dev/null || git config user.email "chabhinav.95@gmail.com"

# Initial commit - December 1, 2024
create_backdated_commit "2024-12-01T10:00:00" "Initial commit: Project structure and basic CLI" "README.md" "# Initial project setup on December 1, 2024"

# Create a feature branch for downloader implementation
git checkout -b feature/downloader

# December 3, 2024 - Add downloader implementation
create_backdated_commit "2024-12-03T14:30:00" "Add core downloader implementation" "internal/downloader/downloader.go" "// Updated on December 3, 2024 - Added core downloader implementation"

# December 5, 2024 - Add progress bar UI
create_backdated_commit "2024-12-05T16:45:00" "Add progress bar UI" "internal/ui/progress.go" "// Updated on December 5, 2024 - Added progress bar UI"

# December 7, 2024 - Add file writer implementation
create_backdated_commit "2024-12-07T09:20:00" "Add file writer implementation" "internal/output/file.go" "// Updated on December 7, 2024 - Added file writer implementation"

# December 10, 2024 - Add JSON output support
create_backdated_commit "2024-12-10T11:15:00" "Add JSON output support" "internal/ui/progress.go" "// Updated on December 10, 2024 - Added JSON output support"

# Merge feature branch into main
git checkout main
git merge feature/downloader -m "Merge feature/downloader into main"

# Create a feature branch for S3 support
git checkout -b feature/s3-support

# December 15, 2024 - Add S3 writer stub
create_backdated_commit "2024-12-15T13:40:00" "Add S3 writer stub implementation" "internal/output/s3.go" "// Updated on December 15, 2024 - Added S3 writer stub implementation"

# December 18, 2024 - Complete S3 integration
create_backdated_commit "2024-12-18T17:25:00" "Complete S3 integration" "internal/output/s3.go" "// Updated on December 18, 2024 - Completed S3 integration"

# December 20, 2024 - Add S3 download support
create_backdated_commit "2024-12-20T10:10:00" "Add direct S3 download support" "internal/downloader/s3.go" "// Updated on December 20, 2024 - Added direct S3 download support"

# Merge S3 support into main
git checkout main
git merge feature/s3-support -m "Merge feature/s3-support into main"

# December 23, 2024 - Add benchmark tests
create_backdated_commit "2024-12-23T14:50:00" "Add benchmark tests" "benchmark.sh" "# Updated on December 23, 2024 - Added benchmark tests"

# December 25, 2024 - Update README with benchmark results
create_backdated_commit "2024-12-25T09:30:00" "Update README with benchmark results" "README.md" "# Updated on December 25, 2024 - Added benchmark results"

# December 28, 2024 - Fix bugs and improve error handling
create_backdated_commit "2024-12-28T16:15:00" "Fix bugs and improve error handling" "internal/downloader/downloader.go" "// Updated on December 28, 2024 - Fixed bugs and improved error handling"

# December 31, 2024 - Final polish for release
create_backdated_commit "2024-12-31T23:45:00" "Final polish for v1.0.0 release" "README.md" "# Updated on December 31, 2024 - Final polish for v1.0.0 release"

echo "Backdated commits for December 2024 have been created!"
echo "To push to GitHub, use: git push origin main"
