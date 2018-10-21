# TypedGit

This project is inspired by the original RawGit project : https://github.com/rgrove/rawgit

Since this service is going to shut down, I decided to write my own thing

## How it works ?

It works exactly like RawGit, with /github at the beginning:
```
GET /github/<owner>/<repo>/<branch>/<path_to_the_file>

GET /github/AlexFlipnote/Discord_Theme/master/theme.css
```

For GitLab, just replace /github by /gitlab:
```
GET /gitlab/Bowser65/React-Modesta/master/demos/img/wallpaper3.jpg
```
