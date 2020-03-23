# nordshare
##the app for sharing notes

read password:
* reading note
* reading time to note expiration

write password:
* reading note
* reading time to note expiration
* modifying content of note
* deleting note
---

| read password set | write password set| read access                 | write access           
| :---:             |:---:              |:---:                        |:---:   
| √                 | √                 | w/ read or write password   | with write password
| √                 | X                 | w/ read password            | X
| X                 | √                 | w/o password (open)         | with write password
| X                 | X                 | w/o password (open)         | X
