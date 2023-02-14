Require autocomm
Require git

Output cinema/movie.gif

# Set up terminal settings
Set FontSize 36
Set Width 1800
Set Height 1100
Set PlaybackSpeed 0.3

Type@50ms "echo 'dev: \n go run . --debug' >> makefile"
Sleep 1s
Enter

Type "git add makefile"
Sleep 1s
Enter

Type "autocomm"
Sleep 1s
Enter

Sleep 8s

Down@500ms 3
Sleep 3s
Up@500ms 2
Sleep 4s
Enter

# Admire the output for a bit
Sleep 3s
Type "git push"
Sleep 3s
