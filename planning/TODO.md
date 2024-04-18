# TODO

-   [] Build and test on windows
-   [] Ensure the output file can be created. No point in spinning minutes through data to fail at the last step.
-   [] Support clean shutdown. Ctrl+c, context cancel and ensure freq table is saved
-   [] Appears the csv reader support stripping out comments, so I should use that instead (Comment rune on reader)
-   [] Document all the possible ways of using the CLI args E.g -o, --out, -out, --o, -o=something.
-   [] Document the default output path name resolving. [--help and README]
-   [] Document some examples. [--help and README]
-   [] Document that the Word ngrams doesn't rip out punctionation
    -   [] Add an extra update task that can be run to filter out certain things. For example rip out where the words are
        just "\* \*" etc.
-   [] Add github actions to build release binaries for supported platforms. Update README for installation steps.

## Done

-   [x] Rename go-analysis to go-analyse. It sounds better :-D
    -   [x] Rename every mention in the repo
    -   [x] Rename the git repo and update local refs
-   [x] Promote internal/alphabet package to text/alphabet.
-   [x] Tokens, monogram should have specialized fast path. Do a benchmark first and after. Document it.
-   [x] Support reading from a zip file
-   [x] -v should be -verbose, so make -version (instead of -v)
-   [x] Add a command to list out the available languages (so if no lang file, show builtins)
