Create a new version of the action by creating a new tag and pushing it to the repository.
```bash
git add .
git commit -m "Initial commit"
git tag -a v1.0.1 -m "Version 1.0.1"
git push origin main --tags
```
After that go to the repository and create a new release with the same tag name.