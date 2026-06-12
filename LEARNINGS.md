## Floor 2 - Create a blob object

1. If you call `os.Create("parentDir/filename")` or sth like it,
   - The Go binary won't make the parent dir for you.
   - You need to make sure the parent dir exists first with
   - `os.Mkdir()`.

2. "Bad file descriptor" probably means the descriptor doesn't
   - have the necessary permissions to perform a certain
   - operation (e.g., write or read) on the file.

3. When converting an int to a string, don't use `string(v)`.
   - It'll evaluate to a string of 1 rune, not a string of digits.
   - Use `strconv.Itoa(v)` instead.
