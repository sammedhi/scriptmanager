A simple CLI utility made with Cobra to manage single scripts living on an ftp server.

This is made for file that are not manage in a dev setup in which the user don't have control over the file environment, if you do have control please, just use git.

```
Usage:
  scm [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  fetch       fetch the script at the target path
  help        Help about any command
  open        open one of the scripts fetched previously
  upload      Upload a script to the server

Flags:
  -h, --help   help for scm
```

The workflow with scm is the following, first you can fetch a file on an ftp server using the **fetch** command this will fetch the file at the given url / path.

What scm will do is store the file in a custom directory and keep the records of the file origin, meaning that from now you can easily update or reupload the file just by calling **upload** or **fetch** using the --name parameter.